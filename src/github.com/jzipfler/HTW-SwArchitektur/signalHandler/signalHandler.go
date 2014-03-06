// signalHandler
package signalHandler

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

func SignalHandlerMitChannel(q chan bool) {
	var quit bool

	c := make(chan os.Signal, 1)
	// Signal für CTRL-C abfangen...
	signal.Notify(c, syscall.SIGINT)

	// foreach signal received
	for signal := range c {
		fmt.Println("\nSignal empfangen...")
		fmt.Println(signal.String())
		switch signal {
		case syscall.SIGINT:
			quit = true
		}
		if quit {
			//TODO: Alles schließen!!!
			//os.Exit(0) // os.Exit() muss nicht sein, da das Hauptprogramm auf den Channel wartet.
		}
		// report the value of quit via the channel
		q <- quit
	}
}

func SignalHandler() {
	c := make(chan os.Signal, 1)
	// Signal für CTRL-C abfangen...
	signal.Notify(c, syscall.SIGINT)

	// foreach signal received
	for signal := range c {
		fmt.Println("\nSignal empfangen...")
		fmt.Println(signal.String())
		switch signal {
		case syscall.SIGINT:
			os.Exit(0)
		}
	}
}
