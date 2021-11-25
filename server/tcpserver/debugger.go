// Package tcpserver Package debug implements debug log functions.
package tcpserver

import (
	"fmt"
	"prr.configuration/config"
	"server/hostel"
	"time"
)

// NbLoggedClient Number of logged client. It is race safe since incremented and decremented in loginRequest.execute
// and logoutRequest.execute. All this, because the execute methods are handled in hostelManager function
// that respects Communicating Sequential Processes.
var NbLoggedClient = 0

// LogRisk prefix message with RISK which means that this log is executed from a critical zone.
func LogRisk(message string) {
	debugLog("RISK) " + message)
}

// LogSafe prefix message with SAFE which means that this log is not executed from a critical zone.
func LogSafe(message string) {
	debugLog("SAFE) " + message)
}

func LogMutex(message string){
	debugLog("MUTEX) " + message)
}

// LogRequestResult logs if communication is a success or failed depending on success argument.
func LogRequestResult(response hostel.Response) {
	if response.Success {
		LogRisk(response.Message + " SUCCESS ")
	} else {
		LogRisk(response.Message + " ERROR ")
	}
}

// LogRequestHandling  logs that communication handle is starting.
func LogRequestHandling(communication hostel.Request) {
	LogRisk(communication.ToString() + " HANDLING ")
}

// debugLog logs the message with DEBUG >> prefix.
func debugLog(message string) {
	if config.IsDebug() {
		fmt.Println("DEBUG >>", time.Now().Format(time.Stamp), message)
	}
}
