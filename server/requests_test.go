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

var conn1, conn2 net.Conn
var reader1, reader2 *bufio.Reader

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

	if len(config.GetServers()) < 2 {
		log.Fatal("Tests need at least 2 servers to test all functionality")
	}

	server1 := config.GetServerById(0)
	server2 := config.GetServerById(1)
	var err error

	// Connection
	conn1, err = net.Dial("tcp", fmt.Sprintf("%s:%d", server1.Ip, server1.Port))
	if err != nil {
		log.Fatal(err)
	}
	conn2, err = net.Dial("tcp", fmt.Sprintf("%s:%d", server2.Ip, server2.Port))
	if err != nil {
		log.Fatal(err)
	}

	reader1 = bufio.NewReader(conn1)
	reader2 = bufio.NewReader(conn2)

	// Read server introduction (useless for testing)
	reader1.ReadString('\n')
	reader2.ReadString('\n')

	return func(tb testing.TB) {
		log.Println("Teardown suite")
		conn1.Close()
		conn2.Close()
	}
}

// TestServerRequest is the main test. Setup test suite and then launch every tests
func TestServerRequest(t *testing.T) {
	teardownSuite := setupSuite(t)

	defer teardownSuite(t)

	table := []struct {
		conn             net.Conn
		reader           *bufio.Reader
		name             string
		input            string
		expectedContains string
	}{
		// Test when not logged
		{conn1, reader1, "Try LOGOUT command without login", "LOGOUT", "ERROR"},
		{conn1, reader1, "Try BOOK command without login", "BOOK 1 1 1", "ERROR"},
		{conn1, reader1, "Try FREEROOM command without login", "FREEROOM 1 1", "ERROR"},
		{conn1, reader1, "Try ROOMLIST command without login", "ROOMLIST 1", "ERROR"},

		// Test login
		{conn1, reader1, "Login without username", "LOGIN", "ERROR"},
		{conn1, reader1, "Good login command", "LOGIN Quentin", "RESULT_LOGIN"},
		{conn1, reader1, "Try login when logged (same username)", "LOGIN Quentin", "ERROR"},
		{conn1, reader1, "Try login when logged (other username)", "LOGIN Melvyn", "ERROR"},

		/////////////////////// BOOK COMMAND WITH BAD ARGUMENTS /////////////////////////
		// Test number of arguments
		{conn1, reader1, "BOOK without argument", "BOOK", "ERROR"},
		{conn1, reader1, "BOOK with 1 argument", "BOOK 1", "ERROR"},
		{conn1, reader1, "BOOK with 2 arguments", "BOOK 1 1", "ERROR"},
		{conn1, reader1, "BOOK with more than 3 arguments", "BOOK 1 1 1 1", "ERROR"},

		// Test negative arguments
		{conn1, reader1, "BOOK with negative first argument", "BOOK -1 1 1", "ERROR"},
		{conn1, reader1, "BOOK with negative second argument", "BOOK 1 -1 1", "ERROR"},
		{conn1, reader1, "BOOK with negative third argument", "BOOK 1 1 -1", "ERROR"},

		// Test with bad arguments
		{conn1, reader1, "BOOK with bad room", "BOOK 11 1 1", "ERROR"},
		{conn1, reader1, "BOOK with bad day", "BOOK 1 11 1", "ERROR"},
		{conn1, reader1, "BOOK with bad duration", "BOOK 1 1 11", "ERROR"},

		/////////////////////// FREEROOM COMMAND WITH BAD ARGUMENTS /////////////////////////
		{conn1, reader1, "FREEROOM without argument", "FREEROOM", "ERROR"},
		{conn1, reader1, "FREEROOM with 1 argument", "FREEROOM 1", "ERROR"},
		{conn1, reader1, "FREEROOM with more than 2 arguments", "FREEROOM 1 1 1", "ERROR"},

		// Test negative arguments
		{conn1, reader1, "FREEROOM with negative first argument", "FREEROOM -1 1", "ERROR"},
		{conn1, reader1, "FREEROOM with negative second argument", "FREEROOM 1 -1", "ERROR"},

		// Test with bad arguments
		{conn1, reader1, "FREEROOM with bad arrival", "FREEROOM 11 1", "ERROR"},
		{conn1, reader1, "FREEROOM with bad duration", "FREEROOM 1 11", "ERROR"},

		/////////////////////// ROOMLIST COMMAND WITH BAD ARGUMENTS /////////////////////////
		{conn1, reader1, "ROOMLIST without argument", "ROOMLIST", "ERROR"},
		{conn1, reader1, "ROOMLIST with more than 1 arguments", "ROOMLIST 1 1", "ERROR"},

		// Test negative arguments
		{conn1, reader1, "ROOMLIST with negative argument", "ROOMLIST -1", "ERROR"},

		// Test with bad arguments
		{conn1, reader1, "ROOMLIST with bad arrival", "ROOMLIST 11", "ERROR"},

		// All following tests should be run with clean server and follow this order
		{conn1, reader1, "ROOMLIST ONLY FREE", "ROOMLIST 1", "Free,Free"},
		{conn1, reader1, "FREEROOM before book", "FREEROOM 1 1", "RESULT_FREEROOM 1"},
		{conn1, reader1, "BOOK a room", "BOOK 1 1 1", "RESULT_BOOK"},
		{conn1, reader1, "FREEROOM after book", "FREEROOM 1 1", "RESULT_FREEROOM 2"},
		{conn1, reader1, "ROOMLIST after book", "ROOMLIST 1", "Self reserved,Free"},

		{conn1, reader1, "BOOK room already booked", "BOOK 1 1 1", "ERROR"},

		{conn1, reader1, "LOGOUT SUCESSFUL", "LOGOUT", "RESULT_LOGOUT"},
		{conn1, reader1, "LOGIN as other user", "LOGIN Melvyn", "RESULT_LOGIN"},

		{conn1, reader1, "ROOMLIST with room1 booked by someone", "ROOMLIST 1", "Occupied,Free"},

		{conn1, reader1, "BOOK other room", "BOOK 2 1 1", "RESULT_BOOK"},
		{conn1, reader1, "ROOMLIST with own reservation", "ROOMLIST 1", "Occupied,Self reserved"},

		{conn1, reader1, "FREEROOM on full booked day", "FREEROOM 1 1", "RESULT_FREEROOM 0"},

		// Check if another server is updated
		{conn2, reader2, "LOGIN with username used in another server", "LOGIN Melvyn", "ERROR"},
		{conn2, reader2, "LOGIN on second server", "LOGIN server2", "RESULT_LOGIN"},
		{conn2, reader2, "Check if servers are sync", "FREEROOM 1 1", "RESULT_FREEROOM 0"},
	}

	for _, tc := range table {
		t.Run(tc.name, func(t *testing.T) {
			_, err := fmt.Fprintf(tc.conn, tc.input+"\n")
			if err != nil {
				t.Errorf("Can't send command to server")
			}

			message, _ := tc.reader.ReadString('\n')

			if !strings.Contains(message, tc.expectedContains) {
				t.Errorf("expected %s, got %s", tc.expectedContains, message)
			}
		})
	}

}
