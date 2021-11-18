// Package tcpclient implements a simple TCP client to communicate with a TCP server that
// manages hostel rooms.
package tcpclient

import (
	"bufio"
	"client/fmthostel"
	"fmt"
	"io"
	"log"
	"net"
	"os"
)

// StartClient connects to the server. This function goroutine reads from user and another goroutine reads from
// server. This loops until application is shutdown or connection closed.
func StartClient(servAddress string) {

	// Connection
	conn, err := net.Dial("tcp", servAddress + ":3000")
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
			fmt.Print("> ")
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
