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
	"prr.configuration/config"
)

// StartClient connects to the server. This function starts a goroutine that reads from server and the function itself
// reads user inputs. It loops until application is shutdown (CTRL + C) or connection closed by server.
func StartClient(server *config.Server) {

	// Connection
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", server.Ip, server.Port))
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
