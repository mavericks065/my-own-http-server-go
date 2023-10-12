package main

import (
	"fmt"
	"net"
	"os"
	"strings"
)

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("Logs from your program will appear here!")

	l, err := net.Listen("tcp", "0.0.0.0:4221")
	if err != nil {
		fmt.Println("Failed to bind to port 4221")
		os.Exit(1)
	}

	connection, err := l.Accept()
	defer connection.Close()
	if err != nil {
		fmt.Println("Error accepting connection: ", err.Error())
		os.Exit(1)
	}

	fmt.Println("data received")

	buf := make([]byte, 1024)
	requestBytes, err := connection.Read(buf)
	if err != nil {
		fmt.Println("Error reading data: ", err.Error())
		os.Exit(1)
	}

	req := string(buf[:requestBytes])
	fmt.Println("Server Request decoded: ", req)

	requestParts := strings.Split(req, " ")
	fmt.Println("Server Request parts split: ", requestParts)

	serverResponse := []byte("HTTP/1.1 200 OK\r\n\r\n")
	_, err = connection.Write(serverResponse)
	if err != nil {
		fmt.Println("Error writing response: ", err.Error())
		os.Exit(1)
	}
	fmt.Println("data sent")
	connection.Close()
}
