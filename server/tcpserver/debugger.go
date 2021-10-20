// Package tcpserver implements logic to run server in debug mode.
package tcpserver

import "fmt"

// DebugMode enables debug mode if set to true.
// Debug mode consist of logging server actions and sleeping 20s when 2 clients logs in, in order to
// let enough time to create a "race" condition.
var DebugMode = false

// Number of logged client. It is race safe since incremented and decremented in loginRequest.execute
// and logoutRequest.execute. All this, because the execute methods are handled in hostelManager function
// that respects Communicating Sequential Processes.
var nbLoggedClient = 0


// debugLogRisk prefix message with RISK which means that this log is executed from a critical zone.
func debugLogRisk(message string) {
	debugLog("RISK) " + message)
}

// debugLogSafe prefix message with SAFE which means that this log is not executed from a critical zone.
func debugLogSafe(message string) {
	debugLog("SAFE) " + message)
}

// debugLog logs the message with DEBUG >> prefix.
func debugLog(message string) {
	fmt.Println("DEBUG >> ", message)
}

// debugLogRequestResult logs if request is a success or failed depending on success argument.
func debugLogRequestResult (request hostelRequestable, success bool) {
	if success {
		debugLogRisk(request.toString() + " SUCCESS ")
	} else {
		debugLogRisk(request.toString() + " ERROR ")
	}
}

// debugLogRequestHandling  logs that request handle is starting.
func debugLogRequestHandling (request hostelRequestable) {
	debugLogRisk(request.toString() + " HANDLING ")
}