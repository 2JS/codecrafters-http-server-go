package main

import "strings"

type RequestLine struct {
	method  string
	path    string
	version string
}

type Request struct {
	line    RequestLine
	headers map[string]string
	body    string
}

func NewRequest(requestString string) Request {
	requestSegments := strings.Split(requestString, "\r\n\r\n")
	lineHeader := strings.Split(requestSegments[0], "\r\n")
	statusStrings := strings.SplitN(lineHeader[0], " ", 3)
	status := RequestLine{
		method:  statusStrings[0],
		path:    statusStrings[1],
		version: statusStrings[2],
	}
	return Request{
		line:    status,
		headers: nil,
		body:    requestSegments[1],
	}
}
