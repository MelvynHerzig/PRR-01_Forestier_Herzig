//// Package debug implements debug log functions.
package tcpserver
//
//import (
//	"fmt"
//	"server/tcpserver/clients"
//)
//
//// NbLoggedClient Number of logged client. It is race safe since incremented and decremented in loginRequest.execute
//// and logoutRequest.execute. All this, because the execute methods are handled in hostelManager function
//// that respects Communicating Sequential Processes.
//var NbLoggedClient = 0
//
//// LogRisk prefix message with RISK which means that this log is executed from a critical zone.
//func LogRisk(message string) {
//	debugLog("RISK) " + message)
//}
//
//// LogSafe prefix message with SAFE which means that this log is not executed from a critical zone.
//func LogSafe(message string) {
//	debugLog("SAFE) " + message)
//}
//
//// LogRequestResult logs if request is a success or failed depending on success argument.
//func LogRequestResult(request clients.HostelRequestable, success bool) {
//	if success {
//		LogRisk(request.ToString() + " SUCCESS ")
//	} else {
//		LogRisk(request.ToString() + " ERROR ")
//	}
//}
//
//// LogRequestHandling  logs that request handle is starting.
//func LogRequestHandling(request clients.HostelRequestable) {
//	LogRisk(request.ToString() + " HANDLING ")
//}
//
//// debugLog logs the message with DEBUG >> prefix.
//func debugLog(message string) {
//	fmt.Println("DEBUG >> ", message)
//}