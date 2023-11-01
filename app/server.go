package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"strings"
)

var directory string

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("Logs from your program will appear here!")

	flag.StringVar(&directory, "directory", ".", "directory")
	flag.Parse()

	l, err := net.Listen("tcp", "0.0.0.0:4221")
	handleError("Failed to bind to port 4221", err)

	for {
		connection, err := l.Accept()
		handleError("Error accepting connection: ", err)
		go handleConnection(connection)
	}
}

func handleError(log string, err error) {
	if err != nil {
		fmt.Println(log)
		os.Exit(1)
	}
}

func handleConnection(connection net.Conn) {
	defer connection.Close()
	buf := make([]byte, 1024)
	requestBytes, err := connection.Read(buf)
	handleError("Error reading data: ", err)

	req := string(buf[:requestBytes])

	var responseContent []string

	if strings.HasPrefix(req, "GET") {
		responseContent = HandleGetRequest(req)
	} else if strings.HasPrefix(req, "POST") {
		responseContent = HandlePostRequest(req)
	}

	serverResponse := []byte(strings.Join(responseContent, ""))
	_, err = connection.Write(serverResponse)
	handleError("Error writing response: ", err)
	connection.Close()
}
