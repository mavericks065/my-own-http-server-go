package main

import (
	"fmt"
	"os"
	"strings"
)

func HandlePostFiles(request Request) Response {
	uriParts := strings.Split(request.Uri, "/")
	header := "HTTP/1.1 201 OK"
	contentType := "Content-Type: text/plain"
	writeFile(request.Body, directory+uriParts[2])
	return Response{
		Header:      header,
		ContentType: contentType,
	}
}

func HandleGet() Response {
	header := "HTTP/1.1 200 OK"
	return Response{
		Header:     header,
		StatusCode: 200,
	}
}

func HandleGetEcho(request Request) Response {
	header := "HTTP/1.1 200 OK"
	contentType := "Content-Type: text/plain"
	body := strings.TrimPrefix(request.Uri, "/echo/")

	return Response{
		Header:        header,
		StatusCode:    200,
		ContentType:   contentType,
		ContentLength: len(body),
		Body:          body,
	}
}

func HandleGetUserAgent(request Request) Response {
	header := "HTTP/1.1 200 OK"
	contentType := "Content-Type: text/plain"

	body := request.Headers["User-Agent"]
	return Response{
		header,
		200,
		contentType,
		len(body),
		body,
	}
}

func HandleGetFiles(request Request) Response {
	uriParts := strings.Split(request.Uri, "/")
	header := "HTTP/1.1 200 OK"
	contentType := "Content-Type: text/plain"
	var body string
	var statusCode int
	if _, statErr := os.Stat(directory + uriParts[2]); statErr == nil {
		contentType = "Content-Type: application/octet-stream"
		body = readFileContent(uriParts[2])
		statusCode = 200
	} else {
		header = "HTTP/1.1 404 Not Found response"
		body = ""
		statusCode = 404
	}
	return Response{
		Header:        header,
		StatusCode:    statusCode,
		ContentType:   contentType,
		ContentLength: len(body),
		Body:          body,
	}
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

func writeFile(body string, filePath string) {
	fileContent := []byte(body)
	err := os.WriteFile(filePath, fileContent, 0644)
	if err != nil {
		fmt.Println("Error while writing file content: ", err.Error())
		panic(err)
	}
}
