// Tests some requests to send to the server.
package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"prr.configuration/config"
	"strings"
	"testing"
)

var conn net.Conn
var reader *bufio.Reader

// setupSuite setups the tests environment. Starts the server and creates the tcp connection.
func setupSuite(tb testing.TB) func(tb testing.TB) {
	log.Println("Setup suite")
	log.Println("These tests are made to test a clean instance of server. Please be sure every server is restarted before executing these tests.")

	config.InitSimple("../config.json")
	if config.GetRoomsCount() != 2 {
		log.Fatal("Tests are made to work when the number of rooms = 2 because 1 is not enough, and with more room, it is harder to fill the hostel.\n" +
			"Please change your config.json and restart all servers")
	}

	if config.GetNightsCount() > 10 {
		log.Fatal("Limits are hardcoded, it means the number of nights can't be more than 10")
	}

	server := config.GetServerRandomly()
	var err error

	// Connection
	conn, err = net.Dial("tcp", fmt.Sprintf("%s:%d", server.Ip, server.Port))
	if err != nil {
		log.Fatal(err)
	}

	reader = bufio.NewReader(conn)

	// Read server introduction (useless for testing)
	reader.ReadString('\n')

	return func(tb testing.TB) {
		log.Println("Teardown suite")
		conn.Close()
	}
}

// TestServerRequest is the main test. Setup test suite and then launch every tests
func TestServerRequest(t *testing.T) {
	teardownSuite := setupSuite(t)

	defer teardownSuite(t)

	table := []struct {
		name             string
		input            string
		expectedContains string
	}{
		// Test when not logged
		{"Try LOGOUT command without login", "LOGOUT", "ERROR"},
		{"Try BOOK command without login", "BOOK 1 1 1", "ERROR"},
		{"Try FREEROOM command without login", "FREEROOM 1 1", "ERROR"},
		{"Try ROOMLIST command without login", "ROOMLIST 1", "ERROR"},

		// Test login
		{"Login without username", "LOGIN", "ERROR"},
		{"Good login command", "LOGIN Quentin", "RESULT_LOGIN"},
		{"Try login when logged (same username)", "LOGIN Quentin", "ERROR"},
		{"Try login when logged (other username)", "LOGIN Melvyn", "ERROR"},

		/////////////////////// BOOK COMMAND WITH BAD ARGUMENTS /////////////////////////
		// Test number of arguments
		{"BOOK without argument", "BOOK", "ERROR"},
		{"BOOK with 1 argument", "BOOK 1", "ERROR"},
		{"BOOK with 2 arguments", "BOOK 1 1", "ERROR"},

		// Test negative arguments
		{"BOOK with negative first argument", "BOOK -1 1 1", "ERROR"},
		{"BOOK with negative second argument", "BOOK 1 -1 1", "ERROR"},
		{"BOOK with negative third argument", "BOOK 1 1 -1", "ERROR"},

		// Test with bad arguments
		{"BOOK with bad room", "BOOK 11 1 1", "ERROR"},
		{"BOOK with bad day", "BOOK 1 11 1", "ERROR"},
		{"BOOK with bad duration", "BOOK 1 1 11", "ERROR"},

		/////////////////////// FREEROOM COMMAND WITH BAD ARGUMENTS /////////////////////////
		{"FREEROOM without argument", "FREEROOM", "ERROR"},
		{"FREEROOM with 1 argument", "FREEROOM 1", "ERROR"},

		// Test negative arguments
		{"FREEROOM with negative first argument", "FREEROOM -1 1", "ERROR"},
		{"FREEROOM with negative second argument", "FREEROOM 1 -1", "ERROR"},

		// Test with bad arguments
		{"FREEROOM with bad arrival", "FREEROOM 11 1", "ERROR"},
		{"FREEROOM with bad duration", "FREEROOM 1 11", "ERROR"},

		/////////////////////// ROOMLIST COMMAND WITH BAD ARGUMENTS /////////////////////////
		{"ROOMLIST without argument", "ROOMLIST", "ERROR"},

		// Test negative arguments
		{"ROOMLIST with negative argument", "ROOMLIST -1", "ERROR"},

		// Test with bad arguments
		{"ROOMLIST with bad arrival", "ROOMLIST 11", "ERROR"},

		// All following tests should be run with clean server and follow this order
		{"ROOMLIST ONLY FREE", "ROOMLIST 1", "Free,Free"},
		{"FREEROOM before book", "FREEROOM 1 1", "RESULT_FREEROOM 1"},
		{"BOOK a room", "BOOK 1 1 1", "RESULT_BOOK"},
		{"FREEROOM after book", "FREEROOM 1 1", "RESULT_FREEROOM 2"},
		{"ROOMLIST after book", "ROOMLIST 1", "Self reserved,Free"},

		{"BOOK room already booked", "BOOK 1 1 1", "ERROR"},

		{"LOGOUT SUCESSFUL", "LOGOUT", "RESULT_LOGOUT"},
		{"LOGIN as other user", "LOGIN Melvyn", "RESULT_LOGIN"},

		{"ROOMLIST with room1 booked by someone", "ROOMLIST 1", "Occupied,Free"},

		{"BOOK other room", "BOOK 2 1 1", "RESULT_BOOK"},
		{"ROOMLIST with own reservation", "ROOMLIST 1", "Occupied,Self reserved"},

		{"FREEROOM on full booked day", "FREEROOM 1 1", "No free room"},
	}

	for _, tc := range table {
		t.Run(tc.name, func(t *testing.T) {
			_, err := fmt.Fprintf(conn, tc.input+"\n")
			if err != nil {
				t.Errorf("Can't send command to server")
			}

			message, _ := reader.ReadString('\n')

			if !strings.Contains(message, tc.expectedContains) {
				t.Errorf("expected %s, got %s", tc.expectedContains, message)
			}
		})
	}

}
