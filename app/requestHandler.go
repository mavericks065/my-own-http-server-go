package main

import (
	"fmt"
	"os"
	"strings"
)

var CRLF = "\r\n"

func HandlePostRequest(req string) []string {
	requestRows := strings.Split(req, CRLF)
	requestPath := strings.Split(requestRows[0], " ")[1]
	uriParts := strings.Split(requestPath, "/")
	header := "HTTP/1.1 201 OK"
	contentType := "Content-Type: text/plain"
	var responseContent []string
	if uriParts[1] == "files" {
		writeFile(requestRows, directory+uriParts[2])
		responseContent = []string{header, CRLF, CRLF, contentType}
	}
	return responseContent
}

func writeFile(requestRows []string, filePath string) {
	requestBody := strings.Join(requestRows[6:], " ")
	fileContent := []byte(requestBody)
	err := os.WriteFile(filePath, fileContent, 0644)
	if err != nil {
		fmt.Println("Error while writing file content: ", err.Error())
		panic(err)
	}
}

func HandleGetRequest(req string) []string {
	requestRows := strings.Split(req, CRLF)
	requestPath := strings.Split(requestRows[0], " ")[1]
	uriParts := strings.Split(requestPath, "/")
	header := "HTTP/1.1 200 OK"
	contentType := "Content-Type: text/plain"
	var responseContent []string

	if requestPath == "/" {
		responseContent = []string{header, CRLF, CRLF}
	} else if uriParts[1] == "echo" {
		body := strings.TrimPrefix(requestPath, "/echo/")
		contentLength := fmt.Sprintf("Content-Length: %d", len(body))
		responseContent = []string{header, CRLF, contentType, CRLF, contentLength, CRLF, CRLF, body, CRLF}
	} else if uriParts[1] == "user-agent" {
		var body string
		for _, row := range requestRows {
			if strings.Contains(row, "User-Agent") {
				body = strings.Split(row, " ")[1]
				continue
			}
		}
		contentLength := fmt.Sprintf("Content-Length: %d\r\n\r\n", len(body))
		responseContent = []string{header, CRLF, contentType, CRLF, contentLength, body, CRLF}
	} else if uriParts[1] == "files" {
		if _, statErr := os.Stat(directory + uriParts[2]); statErr == nil {
			contentType = "Content-Type: application/octet-stream"
			body := readFileContent(uriParts[2])
			contentLength := fmt.Sprintf("Content-Length: %d\r\n\r\n", len(body))
			responseContent = []string{header, CRLF, contentType, CRLF, contentLength, body, CRLF}
		} else {
			responseContent = build404ResponseContent()
		}
	} else {
		responseContent = build404ResponseContent()
	}
	return responseContent
}

func readFileContent(filename string) string {
	fileContent, err := os.ReadFile(directory + filename)
	if err != nil {
		fmt.Println("Error while reading file content: ", err.Error())
		panic(err)
	}
	body := string(fileContent)
	return body
}

func build404ResponseContent() []string {
	header := "HTTP/1.1 404 Not Found response"
	responseContent := []string{header, CRLF, CRLF}
	return responseContent
}
