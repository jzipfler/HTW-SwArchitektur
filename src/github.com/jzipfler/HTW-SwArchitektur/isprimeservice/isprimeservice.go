package main

import (
	"fmt"
	"github.com/jzipfler/HTW-SwArchitektur/service"
	"math/big"
	"strconv"
)

var serviceIsPrime = service.ServiceInfo{
	"isprime",
	"bool",
	"Performs 16 Miller-Rabin tests to check whether x is prime.",
	[]service.ArgumentInfo{
		{"x", "int", "number to test"},
	},
}

// Main function of the "isprime" service
func isprimeHandler(servicecall *service.ServiceCall) string {
	number, _ := strconv.Atoi(servicecall.Arguments[0])
	
	result := big.NewInt(int64(number)).ProbablyPrime(16)
	
	str := fmt.Sprintf("isprime(%d) = %t", number, result)
	
	fmt.Println(str)
	
	return str
}

func main() {
	// register "isprime" as service
	fmt.Println("running...")
	var err error = service.RunService(&serviceIsPrime, isprimeHandler)
	if err != nil {
		fmt.Println("Error occured: ")
		fmt.Println(err)
	}
}
