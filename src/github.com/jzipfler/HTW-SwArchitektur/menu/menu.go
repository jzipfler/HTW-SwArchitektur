// menu
package main

import (
	"fmt"
	"github.com/jzipfler/HTW-SwArchitektur/signalHandler"
	"os"
)

const (
	VERSION           string = "v"
	quit_menu         string = "q"
	zeigeServiceListe string = "1"
	zeigeServiceInfos string = "2"
	starteService     string = "3"
	send_echo_menu    string = "9"
	MENU_HEADER       string = "------------Menü---------------"
	FOOTER            string = "-------------------------------"
	menu_content      string = "Wählen Sie aus den folgenden Einträgen:\n\n" +
		zeigeServiceListe + "\tServiceliste anzeigen" + "\n" +
		zeigeServiceInfos + "\tServicebeschreibung anzeigen" + "\n" +
		starteService + "\tService aufrufen / starten" + "\n" +
		send_echo_menu + "\tNachricht an Server senden" + "\n" +
		"\n" +
		quit_menu + "\tProgramm beenden\t"
)

func main() {
	signalChannel := make(chan bool)
	go signalHandler.SignalHandler(signalChannel)
	// we start the main loop code in a goroutine
	go menu()
	// Blocke die "main" Funktion so lange bis ein Signal empfangen wird.
	select {
	// Sollte ein signal Empfangen werden, speicher es in quit...
	case <-signalChannel:
	}
}

func menu() {
	for line := ""; line != quit_menu; {
		fmt.Printf("%s\n%s\n%s\n", MENU_HEADER, menu_content, FOOTER)
		fmt.Scan(&line) // Lese einen String von der Standardeingabe!
		switch {
		//Verlasse das Programm mit Exit-Code 1
		case line == quit_menu:
			fmt.Println("\nProgramm wird beendet...")
			os.Exit(0)
		case line == zeigeServiceListe:
			fmt.Println("TODO:::")
		case line == zeigeServiceInfos:
			fmt.Println("TODO:::")
		case line == starteService:
			fmt.Println("TODO:::")
		case line == send_echo_menu:
			sendEcho()
		}
	}
}

func sendEcho() {
	fmt.Println("\n\nIn der sendEcho() Funktion\n")
}
