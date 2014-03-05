package main

import (
	"fmt"
	"github.com/jzipfler/HTW-SwArchitektur/service"
)

var serviceConcatenate = service.ServiceInfo{
	"concatenate",
	"variable",
	"Concatenate tow services service1(service2()) as service.",
	[]service.ArgumentInfo{
		{"service1", "string", "first service name"},
		{"service2", "string", "second service name"},
		{"service", "string", "new service name"},
	},
}

// Create composite service service1(service2()) as servicenew
func createCompositeService(service1, service2, servicenew string) {
	serviceInfo := service.ServiceInfo{
		servicenew,
		"variable",
		"service1(service2()).",
		[]service.ArgumentInfo{
			{"void", "void", "void"},
		},
	}
	
	service.RunService(&serviceInfo, func (servicecall *service.ServiceCall) string {
		result, _ := service.CallService(service2)
		result, _ = service.CallService(service1, result)
		return result
	})
}

// Main function of the "concatenate" service
func concatenateHandler(servicecall *service.ServiceCall) string {
	go createCompositeService(servicecall.Arguments[0], servicecall.Arguments[1], servicecall.Arguments[2])
	
	return "1"
}

func main() {
	// register "concatenate" as service
	fmt.Println("running...")
	err := service.RunService(&serviceConcatenate, concatenateHandler)
	if err != nil {
		fmt.Println("Error occured: ")
		fmt.Println(err)
	}
}
