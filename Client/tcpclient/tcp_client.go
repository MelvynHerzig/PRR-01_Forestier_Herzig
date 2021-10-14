// Package tcpclient implements a simple client TCP to communicate with a TCP server that
// manages hostel rooms.
package tcpclient

import (
	"Client/fmthostel"
	"bufio"
	"io"
	"log"
	"net"
	"os"
)

// StartClient Connects to the server. This function go routine reads from user and another go routine reads from
// server. This loops until application is shutdown or connection closed.
func StartClient() {

	// Connection
	conn, err := net.Dial("tcp", "localhost:8000")
	if err != nil {
		log.Fatal(err)
	}

	// Ending signal channel.
	done := make(chan struct{})

	// Start goroutine to read from server.
	go func() {

		input := bufio.NewScanner(conn)
		for input.Scan() {
			fmthostel.FetchDisplayResponse(input.Text())
		}

		log.Println("done")

		done <- struct{}{} // signal the main goroutine
	}()

	// Reads user inputs.
	if _, err := io.Copy(conn, os.Stdin); err != nil {
		log.Fatal(err)
	}

	// Closing
	conn.Close()
	<-done
}
