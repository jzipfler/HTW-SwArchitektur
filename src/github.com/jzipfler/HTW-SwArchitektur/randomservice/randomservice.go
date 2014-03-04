package main

import (
	"fmt"
	"github.com/jzipfler/HTW-SwArchitektur/service"
	"math/rand"
	"strconv"
	"time"
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
	number := rand.Intn(1000)
	
	fmt.Println("random():", number)

	return strconv.Itoa(number)
}

func main() {
	rand.Seed(time.Now().Unix())

	// register "random"-service
	fmt.Println("running...")
	var err error = service.RunService(&serviceRandom, randomHandler)
	if err != nil {
		fmt.Println("Error occured: ")
		fmt.Println(err)
	}
}
