package main

import (
	"fmt"
	"github.com/jzipfler/HTW-SwArchitektur/service"
	"math/rand"
	"strconv"
	"time"
)

var serviceRandom = service.ServiceInfo{
	"random",
	"int",
	"Generates a random int",
	[]service.ArgumentInfo{
		{"void", "void", "no arguments"},
	},
}

// Main function of the "random" service
func randomHandler(servicecall *service.ServiceCall) string {
	number := rand.Int()
	
	fmt.Println("random():", number)

	return strconv.Itoa(number)
}

func main() {
	rand.Seed(time.Now().Unix())

	// register "random" as service
	fmt.Println("running...")
	var err error = service.RunService(&serviceRandom, randomHandler)
	if err != nil {
		fmt.Println("Error occured: ")
		fmt.Println(err)
	}
}
