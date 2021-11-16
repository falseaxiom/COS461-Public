/*****************************************************************************
 * server-go.go
 * Name:  Monique Legaspi
 * NetId: mlegaspi
 *****************************************************************************/

package main

import (
	"fmt"
	// "io"
	"log"
	"net"
	"os"
	// "bufio"
)

const RECV_BUFFER_SIZE = 2048

/* TODO: server()
 * Open socket and wait for client to connect
 * Print received message to stdout
 */
func server(server_port string) {
	// necessary variables
	buf := make([]byte, RECV_BUFFER_SIZE)

	// STEP 1: get info about server ip, throw error if fail
	// STEP 2: make socket
	// STEP 3: bind
	// STEP 4: listen for incoming connections
	listen, err := net.Listen("tcp", ":"+server_port)
	if err != nil {
		fmt.Println("server: failed to listen:", err)
		return
	}

	for true {
		// STEP 5: wait for connections, accept when we get one
		connect, err := listen.Accept()
		if err != nil {
			fmt.Println("server: failed to accept connection:", err)
			return
		}

		// STEP 6: receive message from client & print to stdout
		for true {
			msgsize, err := connect.Read(buf[0:])
			if msgsize == 0 { // if msgsize is 0, client message is over
				break
			}
			if err != nil {
				fmt.Println("server: failed to read message:", err)
				return
			}

			// print to stdout
			fmt.Print(string(buf[0:msgsize]))
		}

		// STEP 7: close the new socket
		connect.Close()
	}

}

// Main parses command-line arguments and calls server function
func main() {
	if len(os.Args) != 2 {
		log.Fatal("Usage: ./server-go [server port]")
	}
	server_port := os.Args[1]
	server(server_port)
}
