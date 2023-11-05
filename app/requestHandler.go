package main

import (
	"fmt"
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
	StatusCode    int
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

func HandleRequest(req string) string {
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
		response := HandleGetFiles(request)
		contentLength := fmt.Sprintf("Content-Length: %d\r\n\r\n", len(response.Body))
		if response.StatusCode == 200 {
			responseContent = []string{response.Header, CRLF, response.ContentType, CRLF, contentLength, response.Body, CRLF}
		} else {
			responseContent = []string{response.Header, CRLF, CRLF}
		}
	} else if request.HttpVerb == "POST" && strings.Split(request.Uri, "/")[1] == "files" {
		response := HandlePostFiles(request)
		responseContent = []string{response.Header, CRLF, CRLF, response.ContentType}
	} else {
		responseContent = build404ResponseContent()
	}

	return strings.Join(responseContent, "")
}

func build404ResponseContent() []string {
	header := "HTTP/1.1 404 Not Found response"
	responseContent := []string{header, CRLF, CRLF}
	return responseContent
}
