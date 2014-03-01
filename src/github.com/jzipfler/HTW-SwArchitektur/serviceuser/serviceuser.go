package main

import (
	"github.com/jzipfler/HTW-SwArchitektur/service"
	"fmt"
	"time"
)

func main() {
	fmt.Println("running...")
	
	// Lookup "random"-service address
	for {
		time.Sleep(time.Millisecond * 2000)
		/*address, err := service.LookupServiceAddress("random")
		if err == nil {
			fmt.Println(address.String())
		} else {
			fmt.Println(err.Error())
		}*/
		rnd, err := service.CallService("random")
		if err == nil {
			fmt.Println("random():", rnd)
		} else {
			fmt.Println("error: CallService()")
		}
	}
}
