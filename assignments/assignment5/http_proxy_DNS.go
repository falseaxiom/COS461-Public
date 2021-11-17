/*****************************************************************************
 * http_proxy_DNS.go
 * Names:  Monique Legaspi
 * NetIds: mlegaspi
 *****************************************************************************/

// TODO: implement an HTTP proxy with DNS Prefetching

// Note: it is highly recommended to complete http_proxy.go first, then copy it
// with the name http_proxy_DNS.go, thus overwriting this file, then edit it
// to add DNS prefetching (don't forget to change the filename in the header
// to http_proxy_DNS.go in the copy of http_proxy.go)

package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"net"
	"net/html"
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
	// coded in OH with help from Sebastian Guzman
	buf_s := bufio.NewReader(connect_s)
	var buf bytes.Buffer
	tee := io.TeeReader(buf_s, &buf)
	io.Copy(connect_c, tee)
	new_reader := bytes.NewReader(buf.Bytes())

	/*** OLD CODE--would only pass 4-6 out of 8 tests in test_proxy_conc.py ***/
	// // response, err := http.ReadResponse(buf_s, request)
	// if err != nil {
	// 	connect_c.Write([]byte("500 'Internal Error': failed to receive response from server\n"))
	// 	connect_s.Close()
	// 	connect_c.Close()
	// 	return
	// }
	// // send response message AS-IS to client via appropriate socket
	// if err := response.Write(connect_c); err != nil {
	// 	connect_c.Write([]byte("500 'Internal Error': failed to write response to client\n"))
	// 	connect_s.Close()
	// 	connect_c.Close()
	// 	return
	// }

	// tokenize & parse html, do DNS prefetchin
	go parse(new_reader)

	// close connection to client
	connect_s.Close()
	connect_c.Close()
	return
}

// parsing
func parse(r io.Reader) {
	doc, err := html.Parse(r)
	if err != nil {
		// connect_c.Write([]byte("500 'Internal Error': failed to parse response\n"))
		// connect_s.Close()
		// connect_c.Close()
		return
	}
	for n := doc; n != nil; n = n.NextSibling {
		dns(n)
	}
}

// DNS prefetching recursive function
// (code referenced from net/html documentation)
func dns(n *html.Node) {
	if n.Type == html.ElementNode && n.Data == "a" {
		for _, a := range n.Attr {
			if a.Key == "href" {
				net.LookupHost(a.Val)
				break
			}
		}
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		dns(c)
	}
}

func main() {
	if len(os.Args) != 2 {
		log.Fatal("Usage: ./http_proxy [port]")
	}
	port := os.Args[1]
	conc(port)
}
