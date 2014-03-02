package main

import (
	"fmt"
	"github.com/jzipfler/HTW-SwArchitektur/service"
)

// Example service info definition
var serviceRandom = service.ServiceInfo{
	"random",
	"int",
	"Generates a random int",
	[]service.ArgumentInfo{
		{"void", "void", "no arguments"},
	},
}

// Main function of the "random service"
func randomHandler(servicecall *service.ServiceCall) string {
	// TODO: generate and return random number
	fmt.Println("random():", 42)

	return "42"
}

func main() {
	// register "random"-service
	fmt.Println("running...")
	var err error = service.RegisterService(&serviceRandom, randomHandler)
	if err != nil {
		fmt.Println("Error occured: ")
		fmt.Println(err)
	}
}
