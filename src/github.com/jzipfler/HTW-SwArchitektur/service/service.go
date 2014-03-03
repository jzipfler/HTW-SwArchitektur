// Package service provides functions and types to create
// and manage services in a LAN network. It also offers
// functions to query information about running services.
// The service discovery is done via a central service called
// "Registry", which itself is discovered via multicast. This
// means is it not necessary to configure static address information.
package service

import (
	"encoding/json"
	"fmt"
	"net"
	"time"
	"strconv"
)

// Information about service argument/parameter.
type ArgumentInfo struct {
	Name        string
	Type        string
	Description string
}

// Information about a service.
type ServiceInfo struct {
	Name        string
	ResultType  string
	Description string
	Arguments   []ArgumentInfo
}

// Information about service that belongs to a specific address.
type ServiceInfoAddress struct {
	Address string
	Info    ServiceInfo
}

// Call parameter for a service.
// This structure is sent to a service when it's invoked.
type ServiceCall struct {
	Name      string
	Arguments []string
}

// Return value of a service.
// This structure is sent upon return of a service.
type ServiceResult struct {
	Result string
}

// Definition of the Service handler function, which will be
// invoked when the service is being called.
type ServiceHandler func(*ServiceCall) string

// Service information lookup request. This is used to query
// information about a service. Valid values are "address"
// (which returns the network address of the given service name,
// "info" (which returns information about the given service name
// and "list" (which returns a map (name to info) containing all
// available services.
type LookupInfoRequest struct {
	Operation   string
	ServiceName string
}

// Response to a service address lookup request (holds service address).
type LookupAddressResponse struct {
	Address net.TCPAddr
}

var (
	// Multicast address for resolution of the registry address.
	MULTICAT_ADDR, _ = net.ResolveUDPAddr("udp4", "224.0.0.1:32001")
	// Used to send multicast messages to own address.
	MULTICAT_SELF_ADDR, _ = net.ResolveUDPAddr("udp4", "127.0.0.1:32001")
	// Any UDP address.
	UDP_ANY_ADDR, _ = net.ResolveUDPAddr("udp4", "0.0.0.0:0")
	// Any TCP address.
	TCP_ANY_ADDR, _ = net.ResolveTCPAddr("tcp4", "0.0.0.0:0")
	// UDP protocol to use (any of: "udp", "udp4", "udp6").
	UDP_PROTOCOL = "udp"
	// TCP protocol to use (any of: "tcp", "tcp4", "tcp6")
	TCP_PROTOCOL = "tcp"
	// Maximum packet/buffer size for send/receive calls.
	PACKET_SIZE = 0x10000
	// Operation for LookupInfoRequest: get service address.
	OPERATION_ADDRESS = "address"
	// Operation for LookupInfoRequest: get service info.
	OPERATION_INFO = "info"
	// Operation for LookupInfoRequest: get service list.
	OPERATION_LIST = "list"
	// Map which contains information about all available services.
	// The mapping is from service name to service information.
	services = make(map[string]ServiceInfoAddress)
)

// Returns the address of any registry which is currently active.
func GetRegistryAddress() (*net.TCPAddr, error) {
	request := LookupInfoRequest{OPERATION_ADDRESS, "registry"}
	response := LookupAddressResponse{}
	buffer := make([]byte, PACKET_SIZE)

	connection, err := net.ListenUDP(UDP_PROTOCOL, UDP_ANY_ADDR)
	if err != nil {
		return nil, err
	}
	defer connection.Close()

	connection.SetReadDeadline(time.Now().Add(time.Second * 4))
	bytes, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}
	_, err = connection.WriteToUDP(bytes, MULTICAT_ADDR)
	if err != nil {
		return nil, err
	}
	_, err = connection.WriteToUDP(bytes, MULTICAT_SELF_ADDR)
	if err != nil {
		return nil, err
	}
	length, address, err := connection.ReadFromUDP(buffer)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(buffer[:length], &response)
	if err != nil {
		return nil, err
	}

	return &net.TCPAddr{address.IP, response.Address.Port, address.Zone}, nil
}

// Get service information for the given operation as JOSN.
// Valid operations are:
// * "address"
// * "info"
// * "list"
func GetServiceData(operation, name string) ([]byte, error) {
	request := LookupInfoRequest{operation, name}
	buffer := make([]byte, PACKET_SIZE)

	address, err := GetRegistryAddress()
	if err != nil {
		return nil, err
	}

	connection, err := net.DialTCP(TCP_PROTOCOL, nil, address)
	if err != nil {
		return nil, err
	}
	defer connection.Close()

	connection.SetReadDeadline(time.Now().Add(time.Second * 4))
	bytes, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}
	_, err = connection.Write(bytes)
	if err != nil {
		return nil, err
	}
	length, err := connection.Read(buffer)
	if err != nil {
		return nil, err
	}

	return buffer[:length], nil
}

// Returns the address for the given service name.
func GetServiceAddress(name string) (*net.TCPAddr, error) {
	response := LookupAddressResponse{}
	buffer, err := GetServiceData(OPERATION_ADDRESS, name);
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(buffer, &response)
	if err != nil {
		return nil, err
	}

	return &response.Address, nil
}

// Returns ServiceInfoAddress for the given service name.
func GetServiceInfo(name string) (*ServiceInfoAddress, error) {
	response := ServiceInfoAddress{}
	buffer, err := GetServiceData(OPERATION_ADDRESS, name);
	if err != nil {
		return nil, err
	}
	
	err = json.Unmarshal(buffer, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

// Returns a map (map[string]ServiceInfoAddress) containing all services.
func GetServiceList() (*map[string]ServiceInfoAddress, error) {
	response := make(map[string]ServiceInfoAddress)
	buffer, err := GetServiceData(OPERATION_ADDRESS, "");
	if err != nil {
		return nil, err
	}
	
	err = json.Unmarshal(buffer, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

// Handles connections to a service and calls the handler specified in RunService().
func handleServiceConnection(connection *net.TCPConn, handler ServiceHandler) error {
	servicecall := ServiceCall{}
	buffer := make([]byte, PACKET_SIZE)

	defer connection.Close()

	length, err := connection.Read(buffer)
	if err != nil {
		return err
	}
	err = json.Unmarshal(buffer[:length], &servicecall)
	if err != nil {
		return err
	}

	ret := handler(&servicecall)

	bytes, err := json.Marshal(ServiceResult{ret})
	if err != nil {
		return err
	}
	_, err = connection.Write(bytes)
	if err != nil {
		return err
	}

	return nil
}

// Registers and starts a service. Any requests to the service are given to
// the user defined handler. Note that this function blocks forever.
func RunService(serviceinfo *ServiceInfo, handler ServiceHandler) error {
	address, err := GetRegistryAddress()
	if err != nil {
		return err
	}

	connection, err := net.DialTCP(TCP_PROTOCOL, nil, address)
	if err != nil {
		return err
	}
	defer connection.Close()
	
	listener, err := net.ListenTCP(TCP_PROTOCOL, TCP_ANY_ADDR)
	if err != nil {
		return err
	}
	defer listener.Close()
	
	address, err = net.ResolveTCPAddr(TCP_PROTOCOL, listener.Addr().String())
	if err != nil {
		return err
	}
	
	bytes, err := json.Marshal(ServiceInfoAddress{strconv.Itoa(address.Port), *serviceinfo})
	if err != nil {
		return err
	}

	_, err = connection.Write(bytes)
	if err != nil {
		return err
	}

	for {
		connection, err := listener.AcceptTCP()
		if err == nil {
			go handleServiceConnection(connection, handler)
		}
	}

	return nil
}

// Server which listens for incoming multicast requests. Upon receive of a
// request it sends the registry address to the asking client.
func registryLookupService(address *net.TCPAddr) error {
	response := LookupAddressResponse{*address}
	buffer := make([]byte, PACKET_SIZE)

	connection, err := net.ListenMulticastUDP(UDP_PROTOCOL, nil, MULTICAT_ADDR)
	if err != nil {
		return err
	}
	defer connection.Close()

	bytes, err := json.Marshal(response)
	if err != nil {
		return err
	}

	for {
		_, sender, err := connection.ReadFromUDP(buffer)
		if err != nil {
			return err
		}
		_, err = connection.WriteToUDP(bytes, sender)
		if err != nil {
			return err
		}
	}

	return nil
}

// Handles new connections to the registry server. For example
// address lookup requests or query service info requests.
func handleRegistryConnection(connection *net.TCPConn) error {
	serviceinfoaddress := ServiceInfoAddress{}
	lookuprequest := LookupInfoRequest{}
	buffer := make([]byte, PACKET_SIZE)

	defer connection.Close()

	length, err := connection.Read(buffer)
	if err != nil {
		return err
	}
	err = json.Unmarshal(buffer[:length], &lookuprequest)
	if err != nil {
		return err
	}
	err = json.Unmarshal(buffer[:length], &serviceinfoaddress)
	if err != nil {
		return err
	}

	if lookuprequest.Operation == OPERATION_ADDRESS {
		fmt.Println("service address:", lookuprequest)
		address, _ := net.ResolveTCPAddr(TCP_PROTOCOL, services[lookuprequest.ServiceName].Address)
		bytes, _ := json.Marshal(LookupAddressResponse{*address})
		connection.Write(bytes)
	} else if lookuprequest.Operation == OPERATION_INFO {
		fmt.Println("service info:", lookuprequest)
		bytes, _ := json.Marshal(services[lookuprequest.ServiceName])
		connection.Write(bytes)
	} else if lookuprequest.Operation == OPERATION_LIST {
		fmt.Println("service list:", lookuprequest)
		bytes, _ := json.Marshal(services)
		connection.Write(bytes)
	} else if serviceinfoaddress.Address != "" {
		address, _ := net.ResolveTCPAddr(TCP_PROTOCOL, connection.RemoteAddr().String())
		address.Port, _ = strconv.Atoi(serviceinfoaddress.Address)
		serviceinfoaddress.Address = address.String()
		fmt.Println("service registered:", address, serviceinfoaddress)
		services[serviceinfoaddress.Info.Name] = serviceinfoaddress
	}

	return nil
}

// Starts a registry server on "0.0.0.0" alias any address. Note that this
// function blocks forever.
func RunRegistryServer() error {
	listener, err := net.ListenTCP(TCP_PROTOCOL, TCP_ANY_ADDR)
	if err != nil {
		return err
	}
	address, err := net.ResolveTCPAddr(TCP_PROTOCOL, listener.Addr().String())
	if err != nil {
		return err
	}

	defer listener.Close()

	go registryLookupService(address)

	for {
		connection, err := listener.AcceptTCP()
		if err == nil {
			go handleRegistryConnection(connection)
		}
	}
}

// Invokes the service specified by name with the given arguments.
func CallService(name string, args ...string) (string, error) {
	servicecall := ServiceCall{name, args}
	serviceresult := ServiceResult{}
	buffer := make([]byte, PACKET_SIZE)

	address, err := GetServiceAddress(name)
	if err != nil {
		return "", err
	}
	connection, err := net.DialTCP(TCP_PROTOCOL, nil, address)
	if err != nil {
		return "", err
	}
	defer connection.Close()

	bytes, err := json.Marshal(servicecall)
	if err != nil {
		return "", err
	}
	_, err = connection.Write(bytes)
	if err != nil {
		return "", err
	}

	length, err := connection.Read(buffer)
	if err != nil {
		return "", err
	}
	err = json.Unmarshal(buffer[:length], &serviceresult)
	if err != nil {
		return "", err
	}

	return serviceresult.Result, nil
}
