// Package main implements a simple client TCP to communication with a TCP server that
// manages hostel rooms.
package main

import (
	"Client/fmthostel"
	"bufio"
	"io"
	"log"
	"net"
	"os"
)

// main Connects to the server and start a writing goroutine.
func main() {

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

	// Reading user input
	if _, err := io.Copy(conn, os.Stdin); err != nil {
		log.Fatal(err)
	}

	// Closing
	conn.Close()
	<-done
}