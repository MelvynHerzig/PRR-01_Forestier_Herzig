// Package fmthostel is responsible for reading hostel server response message
// and format them for user.
// Good response start with RESULT_<COMMAND ANSWERED>.
// Errors response start with ERROR
package fmthostel

import (
	"fmt"
	"strings"
)

// FetchDisplayResponse Fetch server responses then format and display the answer.
func FetchDisplayResponse(response string) {
	trimResponse := strings.TrimSpace(response)
	splitsResponse := strings.SplitN(trimResponse, " ", 2)

	switch splitsResponse[0] {
	case "WELCOME":
		displayWelcome(splitsResponse[1])
	case "ERROR":
		displayError(splitsResponse)
	case "RESULT_LOGIN":
		displayLoginSuccess()
	case "RESULT_LOGOUT":
		displayLogoutSuccess()
	case "RESULT_BOOK":
		displayBook(splitsResponse[1])
	case "RESULT_ROOMLIST":
		displayRoomlist(splitsResponse[1])
	case "RESULT_FREEROOM":
		displayFreeroom(splitsResponse[1])
	default:
		fmt.Println("Unknown server response")
	}
}

// displayError Prints hostel server welcome message.
func displayWelcome(splitResponse string) {
	args := strings.Split(splitResponse, "- ")

	// welcoming
	fmt.Println(args[0])

	// Commands
	fmt.Println("Available commands:")
	for i := 1; i < len(args); i++ {
		fmt.Println("- ", args[i])
	}
}

// displayError Prints error message to user.
func displayError(splitsResponse []string) {

	fmt.Print("/!\\ Operation failed ")

	if len(splitsResponse) == 2 {
		fmt.Print("with message: " + splitsResponse[1])
	}

	fmt.Print("\n")
}

// displayLoginSuccess Prints that the user successfully logged in.
func displayLoginSuccess() {

	fmt.Println("Login success")
}

// displayLogoutSuccess Prints that the user successfully logged out.
func displayLogoutSuccess() {

	fmt.Println("Logout success")
}

// displayError Prints booking summary confirmation.
func displayBook(splitResponse string) {
	args := strings.Split(splitResponse, " ")

	fmt.Println("You successfully booked room ", args[0], " for ", args[2], " night(s), starting day", args[1])
}

// displayRoomlist Prints the room list with their state.
func displayRoomlist(splitResponse string) {

	states := strings.Split(splitResponse, ",")

	fmt.Println("Room no : state ")
	for i, v := range  states {
		fmt.Println(i + 1," : ", v )
	}
}

// displayRoomlist Prints the free room found.
func displayFreeroom(splitResponse string) {

	args := strings.Split(splitResponse, " ")

	fmt.Println("Room  ", args[0], " is free from day ", args[1], " during ", args[2], " nights.")
}


