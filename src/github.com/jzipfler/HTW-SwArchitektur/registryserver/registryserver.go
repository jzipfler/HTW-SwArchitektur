package main

import (
	"github.com/jzipfler/HTW-SwArchitektur/service"
	"fmt"
)

func main() {
	// start registry server
	fmt.Println("running...")
	service.RunRegistryServer()
}
