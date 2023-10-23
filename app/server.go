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
	fmt.Println("SERVER REQUEST DECODED: ", req)

	responseContent := handleRequest(req)

	fmt.Println("RESPONSE CONTENT: ", strings.Join(responseContent, ""))
	serverResponse := []byte(strings.Join(responseContent, ""))
	_, err = connection.Write(serverResponse)
	if err != nil {
		fmt.Println("Error writing response: ", err.Error())
		os.Exit(1)
	}
	fmt.Println("data sent")
	connection.Close()
}

func handleRequest(req string) []string {
	fmt.Println("IDEEEEEEEEEEEEAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAS")
	requestRows := strings.Split(req, "\r\n")

	requestPath := strings.Split(requestRows[0], " ")[1]
	uriParts := strings.Split(requestPath, "/")
	header := "HTTP/1.1 200 OK\r\n"
	contentType := "Content-Type: text/plain\r\n"
	var responseContent []string

	if requestPath == "/" {
		responseContent = []string{header, "\r\n"}
	} else if uriParts[1] == "echo" {
		body := strings.TrimPrefix(requestPath, "/echo/")
		contentLength := fmt.Sprintf("Content-Length: %d\r\n\r\n", len(body))
		responseContent = []string{header, contentType, contentLength, body, "\r\n"}
	} else if uriParts[1] == "user-agent" {
		var body string
		for _, row := range requestRows {
			if strings.Contains(row, "User-Agent") {
				body = strings.Split(row, " ")[1]
				continue
			}
		}
		contentLength := fmt.Sprintf("Content-Length: %d\r\n\r\n", len(body))
		responseContent = []string{header, contentType, contentLength, body, "\r\n"}
	} else {
		header = "HTTP/1.1 404 Not Found response\r\n\r\n"
		responseContent = []string{header}
	}
	return responseContent
}
