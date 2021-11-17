/*****************************************************************************
 * http_proxy.go
 * Names:  Monique Legaspi
 * NetIds: mlegaspi
 *****************************************************************************/

package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
)

// TODO: implement an HTTP proxy
// handles clients concurrently
func conc(port string) {
	/*** I. GETTING REQUESTS FROM CLIENTS ***/

	// listen for client on specified port
	listen, err := net.Listen("tcp", ":"+port) //does it use http instead of tcp??
	if err != nil {
		fmt.Println("500 'Internal Error': failed to listen. Error statement:", err)
		return
	}

	// while loop, so clients can be handled concurrently
	for true {
		// connect to client
		connect_c, err := listen.Accept()
		if err != nil {
			connect_c.Write([]byte("500 'Internal Error': failed to accept client connection\n"))
			connect_c.Close()
			return
		}

		// open new goroutine proxy for each client
		go proxy(port, connect_c)
	}
}

// HTTP proxy function
func proxy(port string, connect_c net.Conn) {

	// read data from client
	buf_c := bufio.NewReader(connect_c)
	request, err := http.ReadRequest(buf_c)
	if err != nil {
		connect_c.Write([]byte("500 'Internal Error': failed to read data\n"))
		connect_c.Close()
		return
	}

	// check for properly-formatted HTTP request
	if request.Method != "GET" {
		connect_c.Write([]byte("500 'Internal Error': failed to read data\n"))
		connect_c.Close()
		return
	}

	// reformat request
	request.Header.Set("Connection", "close")

	/*** II. SENDING REQUESTS TO SERVERS ***/

	// make connection to requested host (using appropriate remote port, or 80 if unspecified)
	server_info := request.Host + ":http"
	connect_s, err := net.Dial("tcp", server_info)
	if err != nil {
		connect_c.Write([]byte("500 'Internal Error': failed to connect to server\n"))
		connect_c.Close()
		return
	}

	// send HTTP request for appropriate resource
	// (send request in relative URL + Host reader format -- see assignment specs)
	if err := request.Write(connect_s); err != nil {
		connect_c.Write([]byte("500 'Internal Error': failed to write request to server\n"))
		connect_s.Close()
		connect_c.Close()
		return
	}

	/*** III. RETURNING RESPONSE TO CLIENTS ***/

	// receive response from remote server
	buf_s := bufio.NewReader(connect_s)
	response, err := http.ReadResponse(buf_s, request)
	if err != nil {
		connect_c.Write([]byte("500 'Internal Error': failed to receive response from server\n"))
		connect_s.Close()
		connect_c.Close()
		return
	}

	// send response message AS-IS to client via appropriate socket
	if err := response.Write(connect_c); err != nil {
		connect_c.Write([]byte("500 'Internal Error': failed to write response to client\n"))
		connect_s.Close()
		connect_c.Close()
		return
	}

	// close connection to client
	connect_s.Close()
	connect_c.Close()
	return
}

func main() {
	if len(os.Args) != 2 {
		log.Fatal("Usage: ./http_proxy [port]")
	}
	port := os.Args[1]
	conc(port)
}
