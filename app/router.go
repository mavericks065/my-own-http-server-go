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
		Header: header,
	}
}

func HandleGetEcho(request Request) Response {
	header := "HTTP/1.1 200 OK"
	contentType := "Content-Type: text/plain"
	body := strings.TrimPrefix(request.Uri, "/echo/")

	return Response{
		Header:        header,
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
		contentType,
		len(body),
		body,
	}
}

func writeFile(body string, filePath string) {
	fileContent := []byte(body)
	err := os.WriteFile(filePath, fileContent, 0644)
	if err != nil {
		fmt.Println("Error while writing file content: ", err.Error())
		panic(err)
	}
}
