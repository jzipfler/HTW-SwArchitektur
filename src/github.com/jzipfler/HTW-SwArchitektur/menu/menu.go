// menu
package main

import (
	"fmt"
	"github.com/jzipfler/HTW-SwArchitektur/service"
	"github.com/jzipfler/HTW-SwArchitektur/signalHandler"
	"os"
	"strconv"
)

const (
	VERSION                string = "v"
	quit_menu              string = "q"
	ZEIGE_SERVICE_LISTE    string = "1"
	ZEIGE_SERVICE_INFOS    string = "2"
	AUFRUFEN_SERVICE       string = "3"
	send_echo_menu         string = "9"
	MENU_HEADER            string = "------------Menü---------------"
	AUSGABE_HEADER         string = "-----------Ausgabe-------------"
	FEHLER_HEADER          string = "-----------Fehler--------------"
	FOOTER                 string = "-------------------------------"
	AUSGABE_BEENDEN        string = "Programm wird beendet."
	ZEILENUMBRUCH          string = "\n"
	EINGABE                string = "Eingabe: "
	EINGABE_SERVICE_NAME   string = "Eingabe des Service-Namen: "
	SERVICE_INFOS_VORGEHEN string = "Sie wollen sich Informationen zu einem Service anzeigen lassen." + ZEILENUMBRUCH +
		"Geben Sie dazu bitte den Namen des Services an."
	AUFRUFEN_SERVICE_VORGEHEN string = "Sie wollen einen Service ausführen." + ZEILENUMBRUCH +
		"Geben Sie dazu bitte den Namen des Services an."
	menu_content string = "Wählen Sie aus den folgenden Einträgen:\n\n" +
		ZEIGE_SERVICE_LISTE + "\tServiceliste anzeigen" + ZEILENUMBRUCH +
		ZEIGE_SERVICE_INFOS + "\tServicebeschreibung anzeigen" + ZEILENUMBRUCH +
		AUFRUFEN_SERVICE + "\tService aufrufen / starten" + ZEILENUMBRUCH +
		send_echo_menu + "\tNachricht an Server senden" + ZEILENUMBRUCH +
		ZEILENUMBRUCH +
		quit_menu + "\tProgramm beenden"
)

func main() {
	// Starte den SignalHandler als goroutine, die bei
	// einem empfangenem Signal das Programm beendet.
	go signalHandler.SignalHandler()
	menu()
}

func menu() {
	for line := ""; line != quit_menu; {
		fmt.Printf("%s\n%s\n%s\n", MENU_HEADER, menu_content, FOOTER)
		fmt.Print(EINGABE)
		fmt.Scan(&line) // Lese einen String von der Standardeingabe!
		fmt.Print(ZEILENUMBRUCH)
		switch {
		//Verlasse das Programm mit Exit-Code 1
		case line == quit_menu:
			informationenAusgeben(AUSGABE_BEENDEN, false)
			os.Exit(0) // oder "return"
		case line == ZEIGE_SERVICE_LISTE:
			zeigeServiceListe()
		case line == ZEIGE_SERVICE_INFOS:
			zeigeServiceInformation()
		case line == AUFRUFEN_SERVICE:
			aufrufenService()
		case line == send_echo_menu:
			sendEcho()
		default:
			informationenAusgeben("Kein solcher Menüpunkt vorhanden.", true)
		}
	}
}

func zeigeServiceListe() {
	serviceListe, err := service.GetServiceList()
	if err != nil {
		informationenAusgeben(err.Error(), true)
	}
	fmt.Println(*serviceListe)
	fmt.Println(ZEILENUMBRUCH)
	for k, v := range *serviceListe {
		fmt.Println("String: " + k)
		fmt.Println("Values:" + ZEILENUMBRUCH)
		fmt.Println("Adresse: " + v.Address + ZEILENUMBRUCH +
			"Service-Name: " + v.Info.Name + ZEILENUMBRUCH +
			"Service-Beschreibung: " + v.Info.Description + ZEILENUMBRUCH +
			ZEILENUMBRUCH +
			"Rückgabewert: " + v.Info.ResultType)
		for argumente := range v.Info.Arguments {
			fmt.Println(strconv.Itoa(argumente) + ". Argument: " + ZEILENUMBRUCH +
				"\tName: " + v.Info.Arguments[argumente].Name + ZEILENUMBRUCH +
				"\tTyp: " + v.Info.Arguments[argumente].Type + ZEILENUMBRUCH +
				"\tBeschreibung: " + v.Info.Arguments[argumente].Description)
		}
	}
}

func zeigeServiceInformation() {
	var serviceName string
	fmt.Println(ZEILENUMBRUCH + SERVICE_INFOS_VORGEHEN)
	fmt.Print(EINGABE_SERVICE_NAME)
	fmt.Scan(&serviceName)
	serviceInformation, err := service.GetServiceInfo(serviceName)
	if err != nil {
		informationenAusgeben(err.Error(), true)
	}
	fmt.Println(*serviceInformation)
}

func aufrufenService() {
	var serviceName string
	fmt.Println(ZEILENUMBRUCH + AUFRUFEN_SERVICE_VORGEHEN)
	fmt.Print(EINGABE_SERVICE_NAME)
	fmt.Scan(&serviceName)
	serviceAusgabe, err := service.CallService(serviceName)
	if err != nil {
		informationenAusgeben(err.Error(), true)
	}
	fmt.Println(serviceAusgabe)
}

func informationenAusgeben(infos string, fehlerAusgabe bool) {
	if fehlerAusgabe {
		fmt.Println(FEHLER_HEADER)
	} else {
		fmt.Println(AUSGABE_HEADER)
	}
	fmt.Println(infos)
	fmt.Println(FOOTER + ZEILENUMBRUCH)
}

func sendEcho() {
	fmt.Println("\n\nIn der sendEcho() Funktion\n")
}
