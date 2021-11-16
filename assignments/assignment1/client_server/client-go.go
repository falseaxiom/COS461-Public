/*****************************************************************************
 * client-go.go
 * Name:  Monique Legaspi
 * NetId: mlegaspi
 *****************************************************************************/

package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"os"
)

const SEND_BUFFER_SIZE = 2048

/* TODO: client()
 * Open socket and send message from stdin.
 */
func client(server_ip string, server_port string) {
	// necessary variables
	buf := make([]byte, SEND_BUFFER_SIZE)

	// STEP 1: get info about server ip
	// STEP 2: make a socket
	// STEP 3: connect to server
	server_info := server_ip + ":" + server_port
	connect, err := net.Dial("tcp", server_info)
	if err != nil {
		fmt.Println("client: failed to connect:", err)
		return
	}
	defer connect.Close()

	// STEP 4: read & send message
	reader := bufio.NewReader(os.Stdin)
	for true {
		// read in message
		msgsize, err := reader.Read(buf)
		if err != nil {
			if err == io.EOF { // if err is EOF, message is done reading
				return
			} else {
				fmt.Println("client: failed to read message:", err)
				return
			}
		}

		// send it iteratively -> chunks of 2048 characters
		if _, err := connect.Write(buf[0:msgsize]); err != nil {
			fmt.Println("client: failed to send message:", err)
			return
		}
	}

	// STEP 5: close the socket
	connect.Close()
}

// Main parses command-line arguments and calls client function
func main() {
	if len(os.Args) != 3 {
		log.Fatal("Usage: ./client-go [server IP] [server port] < [message file]")
	}
	server_ip := os.Args[1]
	server_port := os.Args[2]
	client(server_ip, server_port)
}
