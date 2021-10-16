// Package fmthostel is responsible to read hostel server responses and format them for user.
// Good response starts with RESULT_<COMMAND ANSWERED>.
// Error response starts with ERROR
package fmthostel

import (
	"fmt"
	"strings"
)

// FetchDisplayResponse fetches server responses then formats and displays the answer.
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

// displayError prints hostel server welcome message.
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

// displayError prints error message to user.
func displayError(splitsResponse []string) {

	fmt.Print("/!\\ Operation failed ")

	if len(splitsResponse) == 2 {
		fmt.Print("with message: " + splitsResponse[1])
	}

	fmt.Print("\n")
}

// displayLoginSuccess prints that the user successfully logged in.
func displayLoginSuccess() {

	fmt.Println("Login success")
}

// displayLogoutSuccess prints that the user successfully logged out.
func displayLogoutSuccess() {

	fmt.Println("Logout success")
}

// displayError prints booking summary confirmation.
func displayBook(splitResponse string) {
	args := strings.Split(splitResponse, " ")

	fmt.Println("You successfully booked room ", args[0], " for ", args[2], " night(s), starting night", args[1])
}

// displayRoomlist prints the room list with their state.
func displayRoomlist(splitResponse string) {

	states := strings.Split(splitResponse, ",")

	fmt.Println("Room no : state ")
	for i, v := range  states {
		fmt.Println(i + 1," : ", v )
	}
}

// displayRoomlist prints the free room found.
// displayRoomlist prints the free room found.
func displayFreeroom(splitResponse string) {

	args := strings.Split(splitResponse, " ")

	if args[0] == "0" {
		fmt.Println("No rooms free from night ", args[1], " for ", args[2], " night(s).")
	} else {
		fmt.Println("Room ", args[0], " is free from night ", args[1], " during ", args[2], " night(s).")
	}
}


