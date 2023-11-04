package main

import (
	"fmt"
	"os"
	"strings"
)

type Request struct {
	Uri      string
	HttpVerb string
	Headers  map[string]string
	Body     string
}

type Response struct {
	Header        string
	ContentType   string
	ContentLength int
	Body          string
}

var CRLF = "\r\n"

func createRequest(req string) Request {
	requestRows := strings.Split(req, CRLF)
	httpVerb := strings.Split(requestRows[0], " ")[0]
	uri := strings.Split(requestRows[0], " ")[1]
	headers := extractHeaders(httpVerb, requestRows)
	var body string
	if httpVerb == "POST" {
		body = requestRows[len(requestRows)-1]
	}
	return Request{Uri: uri, HttpVerb: httpVerb, Headers: headers, Body: body}
}

func extractHeaders(httpVerb string, requestRows []string) map[string]string {
	headers := make(map[string]string)
	var lastHeader int
	if httpVerb == "GET" {
		lastHeader = len(requestRows)
	} else {
		lastHeader = len(requestRows) - 1
	}
	for _, str := range requestRows[1:lastHeader] {
		parts := strings.SplitN(str, ":", 2)
		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])
			headers[key] = value
		}
	}
	return headers
}

func HandleRequest(req string) []byte {
	request := createRequest(req)
	var responseContent []string

	if request.HttpVerb == "GET" && request.Uri == "/" {
		response := HandleGet()
		responseContent = []string{response.Header, CRLF, CRLF}
	} else if request.HttpVerb == "GET" && strings.Split(request.Uri, "/")[1] == "echo" {
		response := HandleGetEcho(request)
		contentLength := fmt.Sprintf("Content-Length: %d", response.ContentLength)
		responseContent = []string{response.Header, CRLF, response.ContentType, CRLF, contentLength, CRLF, CRLF, response.Body, CRLF}
	} else if request.HttpVerb == "GET" && strings.Split(request.Uri, "/")[1] == "user-agent" {
		response := HandleGetUserAgent(request)
		contentLength := fmt.Sprintf("Content-Length: %d\r\n\r\n", response.ContentLength)
		responseContent = []string{response.Header, CRLF, response.ContentType, CRLF, contentLength, response.Body, CRLF}
	} else if request.HttpVerb == "GET" && strings.Split(request.Uri, "/")[1] == "files" {
		responseContent = HandleGetRequest(req)
	} else if request.HttpVerb == "POST" && strings.Split(request.Uri, "/")[1] == "files" {
		response := HandlePostFiles(request)
		responseContent = []string{response.Header, CRLF, CRLF, response.ContentType}
	} else {
		responseContent = HandleGetRequest(req)
	}

	return []byte(strings.Join(responseContent, ""))
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
