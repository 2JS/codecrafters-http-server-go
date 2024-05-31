package main

import "fmt"

type Response struct {
	status  Status
	headers map[string]string
	body    string
}

type Status struct {
	version string
	code    int
	message string
}

func NewStatus() Status {
	return Status{
		version: "HTTP/1.1",
		code:    200,
		message: "OK",
	}
}

func (response Response) String() string {
	headerString := ""
	for key, value := range response.headers {
		headerString += fmt.Sprintf("%s: %s\r\n", key, value)
	}
	return fmt.Sprintf(
		"%s\r\n%s\r\n%s",
		response.status,
		headerString,
		response.body,
	)
}

func (status Status) String() string {
	return fmt.Sprintf(
		"%s %d %s",
		status.version,
		status.code,
		status.message,
	)
}
