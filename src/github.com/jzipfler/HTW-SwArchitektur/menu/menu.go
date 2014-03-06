// menu
package main

import (
	"bytes"
	"fmt"
	"github.com/jzipfler/HTW-SwArchitektur/service"
	"github.com/jzipfler/HTW-SwArchitektur/signalhandler"
	"os"
	"strconv"
)

const (
	VERSION                string = "v"
	quit_menu              string = "q"
	ZEIGE_SERVICE_LISTE    string = "1"
	ZEIGE_SERVICE_INFOS    string = "2"
	AUFRUFEN_SERVICE       string = "3"
	KEIN_MENUEEINTRAG      string = "Kein solcher Menüpunkt vorhanden."
	MENU_HEADER            string = "------------Menü---------------"
	AUSGABE_HEADER         string = "-----------Ausgabe-------------"
	FEHLER_HEADER          string = "-----------Fehler--------------"
	FOOTER                 string = "-------------------------------"
	SERVICE_HEADER         string = "___________Service_____________"
	SERVICE_FOOTER         string = "_______________________________"
	AUSGABE_BEENDEN        string = "Programm wird beendet."
	ZEILENUMBRUCH          string = "\n"
	EINGABE                string = "Eingabe: "
	EINGABE_SERVICE_NAME   string = "Eingabe des Service-Namen: "
	SERVICE_INFOS_VORGEHEN string = "Sie wollen sich Informationen zu einem Service anzeigen lassen." + ZEILENUMBRUCH +
		"Geben Sie dazu bitte den Namen des Services an."
	AUFRUFEN_SERVICE_VORGEHEN string = "Sie wollen einen Service ausführen." + ZEILENUMBRUCH +
		"Geben Sie dazu bitte den Namen des Services an."
	AUFRUFEN_SERVICE_PARAMETER_INFO string = "Nachdem Sie den Service gewählt haben," + ZEILENUMBRUCH +
		"müssen Sie nun die erforderlichen Parameter" + ZEILENUMBRUCH +
		"eingeben. Danach wird der Service ausgeführt."
	menu_content string = "Wählen Sie aus den folgenden Einträgen:\n\n" +
		ZEIGE_SERVICE_LISTE + "\tServiceliste anzeigen" + ZEILENUMBRUCH +
		ZEIGE_SERVICE_INFOS + "\tServicebeschreibung anzeigen" + ZEILENUMBRUCH +
		AUFRUFEN_SERVICE + "\tService aufrufen / starten" + ZEILENUMBRUCH +
		ZEILENUMBRUCH +
		quit_menu + "\tProgramm beenden"
)

func main() {
	// Starte den SignalHandler als goroutine, die bei
	// einem empfangenem Signal das Programm beendet.
	go signalhandler.SignalHandler()
	menu()
}

// Zeigt ein Menü an und verarbeitet die darauf folgende Eingabe.
// Dabei verzweigt die Verarbeitung in andere Funktionen.
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
		default:
			informationenAusgeben(KEIN_MENUEEINTRAG, true)
		}
	}
}

// Wenn im Menü die ServiceListe ausgewählt wurde, wird
// diese Funktion aufgerufen, die die ServiceListe von
// der Registry abfragt und ausgibt.
func zeigeServiceListe() {
	var buffer bytes.Buffer
	serviceListe, err := service.GetServiceList()
	if err != nil {
		informationenAusgeben(err.Error(), true)
		return
	}
	for key, serviceInfoAdresse := range *serviceListe {
		buffer.Reset()
		buffer.WriteString(ZEILENUMBRUCH +
			SERVICE_HEADER + ZEILENUMBRUCH + ZEILENUMBRUCH +
			"\t   " + key + ZEILENUMBRUCH +
			SERVICE_FOOTER + ZEILENUMBRUCH)
		buffer.WriteString("SERVICE CONTRACT:" + ZEILENUMBRUCH + ZEILENUMBRUCH)
		buffer.WriteString(verarbeiteServiceInfoAddress(serviceInfoAdresse))
		informationenAusgeben(buffer.String(), false)
	}
}

// Wenn im Menü der Punkt für die Information eines Services
// aufgerufen wurde, wird in diesem Menü nach dem Service
// befragt, über den Informationen eingeholt werden sollen.
// Danach werden die verfügbaren Informationen über diesen
// ausgegeben.
func zeigeServiceInformation() {
	var serviceName string
	fmt.Println(ZEILENUMBRUCH + SERVICE_INFOS_VORGEHEN)
	fmt.Print(EINGABE_SERVICE_NAME)
	fmt.Scan(&serviceName)
	serviceInformation, err := service.GetServiceInfo(serviceName)
	if err != nil {
		informationenAusgeben(err.Error(), true)
		return
	}
	informationenAusgeben(verarbeiteServiceInfoAddress(*serviceInformation), false)
}

// Wenn der Punkt zum ausführen eines Services aus dem Menü
// ausgewählt wurde, wird diese Funktion aufgerufen.
// Diese fragt nach, welcher Service gestartet werden soll
// und verarbeitet alle Informationen zum starten des Services.
func aufrufenService() {
	var serviceName string
	var serviceAusgabe string
	var err error
	fmt.Println(ZEILENUMBRUCH + AUFRUFEN_SERVICE_VORGEHEN)
	fmt.Print(EINGABE_SERVICE_NAME)
	fmt.Scan(&serviceName)
	serviceInformation, err := service.GetServiceInfo(serviceName)
	if err != nil {
		informationenAusgeben(err.Error(), true)
		return
	}
	if serviceInformation.Info.Arguments[0].Type != "void" {
		//informationenAusgeben("Service hat mehrere Parameter.\nDies wird noch nicht unterstützt.", true)
		anzahlParameter := len(serviceInformation.Info.Arguments)
		parameter := make([]string, anzahlParameter)
		fmt.Println(ZEILENUMBRUCH)
		fmt.Println(AUFRUFEN_SERVICE_PARAMETER_INFO)
		fmt.Println(ZEILENUMBRUCH)
		for i := 0; i < len(serviceInformation.Info.Arguments); i++ {
			fmt.Printf("Der %d te Parameter ist vom Typ:\t\t%s\n", (i + 1), serviceInformation.Info.Arguments[i].Type)
			fmt.Println("Dazu gehört folgende Beschreibung:\t" + serviceInformation.Info.Arguments[i].Description)
			fmt.Printf("Parameter %d eingeben: ", (i + 1))
			fmt.Scan(&parameter[i])
			fmt.Println(ZEILENUMBRUCH)
		}
		switch len(serviceInformation.Info.Arguments) {
		case 1:
			serviceAusgabe, err = service.CallService(serviceName, parameter[0])
		case 2:
			serviceAusgabe, err = service.CallService(serviceName, parameter[0], parameter[1])
		case 3:
			serviceAusgabe, err = service.CallService(serviceName, parameter[0], parameter[1], parameter[2])
		default:
			informationenAusgeben("Unbekannte Anzahl Parameter", true)
			return
		}
	} else {
		serviceAusgabe, err = service.CallService(serviceName)
	}
	if err != nil {
		informationenAusgeben(err.Error(), true)
	}
	informationenAusgeben(serviceAusgabe, false)
}

// Diese Funktion dient als Hilfsfunktion.
// Diese wird für jede Ausgabe die gemacht wird genutzt
// um über und unter dem Text einen "Header" bzw. "Footer" anzuzeigen.
func informationenAusgeben(infos string, fehlerAusgabe bool) {
	if fehlerAusgabe {
		fmt.Println(FEHLER_HEADER)
	} else {
		fmt.Println(AUSGABE_HEADER)
	}
	fmt.Println(infos)
	fmt.Println(FOOTER + ZEILENUMBRUCH)
}

// Diese Methode verarbeitet die übergebene ServiceInfoAddress
// und gibt einen String zurück indem diese Datenstruktur
// formatiert und mit Bezeichnern abgelegt ist.
func verarbeiteServiceInfoAddress(serviceInfoAddress service.ServiceInfoAddress) string {
	var buffer bytes.Buffer
	buffer.WriteString("Adresse: " + serviceInfoAddress.Address + ZEILENUMBRUCH +
		"Service-Name: " + serviceInfoAddress.Info.Name + ZEILENUMBRUCH +
		"Service-Beschreibung: " + serviceInfoAddress.Info.Description + ZEILENUMBRUCH +
		ZEILENUMBRUCH +
		"Rückgabewert: " + serviceInfoAddress.Info.ResultType + ZEILENUMBRUCH)
	for index := range serviceInfoAddress.Info.Arguments {
		buffer.WriteString(strconv.Itoa(index+1) + ". Argument: " + ZEILENUMBRUCH +
			"\tBezeichnung:\t" + serviceInfoAddress.Info.Arguments[index].Name + ZEILENUMBRUCH +
			"\tTyp:\t\t" + serviceInfoAddress.Info.Arguments[index].Type + ZEILENUMBRUCH +
			"\tBeschreibung:\t" + serviceInfoAddress.Info.Arguments[index].Description + ZEILENUMBRUCH)
	}
	return buffer.String()
}
