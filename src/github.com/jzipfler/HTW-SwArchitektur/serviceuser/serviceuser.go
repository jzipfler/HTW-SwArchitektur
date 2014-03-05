package main

import (
	"github.com/jzipfler/HTW-SwArchitektur/service"
	"fmt"
	"time"
)

func main() {
	fmt.Println("running...")
	
	_, err := service.CallService("concatenate", "isprime", "random", "isrndprime")
	if err != nil {
		fmt.Println("error: CallService()")
	}
	
	for {
		time.Sleep(time.Millisecond * 2000)
		
		/*rnd, err := service.CallService("random")
		if err == nil {
			fmt.Println("random():", rnd)
		} else {
			fmt.Println("error: CallService()")
		}*/
		
		isprime, err := service.CallService("isrndprime")
		if err == nil {
			fmt.Println(isprime)
		} else {
			fmt.Println("error: CallService()")
		}
		
	}
}
