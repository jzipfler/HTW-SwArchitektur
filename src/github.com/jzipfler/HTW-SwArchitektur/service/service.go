package service

import (
	"encoding/json"
	"fmt"
	"net"
	"time"
)

// Information about service argument/parameter
type ArgumentInfo struct {
	Name        string
	Type        string
	Description string
}

// Information about service
type ServiceInfo struct {
	Name        string
	ResultType  string
	Description string
	Arguments   []ArgumentInfo
}

// Call parameter for a service
type ServiceCall struct {
	Name      string
	Arguments []string
}

// Return value of a service
type ServiceResult struct {
	Result string
}

// Service handler function
type ServiceHandler func(*ServiceCall) string

// Service lookup request (name to address)
type LookupRequest struct {
	ServiceName string
}

// Response to lookup request (holds service address)
type LookupResponse struct {
	Address net.TCPAddr
}

var (
	// Multicast address registry ip resolution
	MULTICAT_ADDR, _ = net.ResolveUDPAddr("udp4", "224.0.0.1:32001")
	// Used to send multicast messages to self
	MULTICAT_SELF_ADDR, _ = net.ResolveUDPAddr("udp4", "127.0.0.1:32001")
	// Any UDP address
	UDP_ANY_ADDR, _ = net.ResolveUDPAddr("udp4", "0.0.0.0:0")
	// Any TCP address
	TCP_ANY_ADDR, _ = net.ResolveTCPAddr("tcp4", "0.0.0.0:0")
	// Used udp protocol (udp, udp4, udp6)
	UDP_PROTOCOL = "udp"
	// Used tcp protocol (tcp, tcp4, tcp6)
	TCP_PROTOCOL = "tcp"
	// Max packt/buffer size for send/receive
	PACKET_SIZE = 0x10000
	// Map: name --> service info
	services = make(map[string]ServiceInfo)
	// Map: name --> service address
	servicestoaddr = make(map[string]string)
)

// Returns the address of the currently active registry (if any).
func LookupRegistryAddress() (*net.TCPAddr, error) {
	request := LookupRequest{}
	response := LookupResponse{}
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

// Returns the address of the given service name.
func LookupServiceAddress(name string) (*net.TCPAddr, error) {
	request := LookupRequest{name}
	response := LookupResponse{}
	buffer := make([]byte, PACKET_SIZE)

	address, err := LookupRegistryAddress()
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
	err = json.Unmarshal(buffer[:length], &response)
	if err != nil {
		return nil, err
	}

	return &(response.Address), nil
}

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

// Registers a service in the registry.
func RegisterService(serviceinfo *ServiceInfo, handler ServiceHandler) error {
	address, err := LookupRegistryAddress()
	if err != nil {
		return err
	}

	connection, err := net.DialTCP(TCP_PROTOCOL, nil, address)
	if err != nil {
		return err
	}
	address, _ = net.ResolveTCPAddr(TCP_PROTOCOL, connection.LocalAddr().String())
	
	err = connection.SetLinger(0)
	if err != nil {
		return err
	}
	
	bytes, err := json.Marshal(serviceinfo)
	if err != nil {
		return err
	}

	_, err = connection.Write(bytes)
	if err != nil {
		return err
	}

	connection.Close()

	listener, err := net.ListenTCP(TCP_PROTOCOL, address)
	if err != nil {
		return err
	}

	defer listener.Close()

	for {
		connection, err := listener.AcceptTCP()
		if err == nil {
			go handleServiceConnection(connection, handler)
		}
	}

	return nil
}

// Answers multicast reguests for the registry address.
func registryLookupService(address *net.TCPAddr) error {
	response := LookupResponse{*address}
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

// Registers new services and service name lookup.
func handleRegistryConnection(connection *net.TCPConn) error {
	serviceinfo := ServiceInfo{}
	lookuprequest := LookupRequest{}
	lookupresponse := LookupResponse{}
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
	err = json.Unmarshal(buffer[:length], &serviceinfo)
	if err != nil {
		return err
	}

	if lookuprequest.ServiceName != "" {
		fmt.Println("lookup service: ", lookuprequest)
		address, _ := net.ResolveTCPAddr(TCP_PROTOCOL, servicestoaddr[lookuprequest.ServiceName])
		lookupresponse.Address = *address
		bytes, _ := json.Marshal(lookupresponse)
		connection.Write(bytes)
	} else if serviceinfo.Name != "" {
		fmt.Println("service registered:", connection.RemoteAddr().String(), serviceinfo)
		services[serviceinfo.Name] = serviceinfo
		servicestoaddr[serviceinfo.Name] = connection.RemoteAddr().String()
	}

	return nil
}

// Starts a registry server.
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

func CallService(name string, args ...string) (string, error) {
	servicecall := ServiceCall{name, args}
	serviceresult := ServiceResult{}
	buffer := make([]byte, PACKET_SIZE)

	address, err := LookupServiceAddress(name)
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
