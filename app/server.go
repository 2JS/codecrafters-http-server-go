package main

import (
	"fmt"
	"net"
	"os"
	"strings"
)

func main() {
	l, err := net.Listen("tcp", "0.0.0.0:4221")
	if err != nil {
		fmt.Println("Failed to bind to port 4221")
		os.Exit(1)
	}

	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}
		go respond(conn)
	}
}

func respond(conn net.Conn) {
	buff := make([]byte, 1024)
	conn.Read(buff)
	request := NewRequest(string(buff))

	response := Response{status: NewStatus()}

	switch {
	case request.line.path == "/":
		break
	case strings.HasPrefix(request.line.path, "/echo"):
		echo := strings.TrimPrefix(request.line.path, "/echo/")
		response.headers = map[string]string{
			"Content-Type":   "text/plain",
			"Content-Length": fmt.Sprint(len(echo)),
		}
		response.body = echo
	case strings.HasPrefix(request.line.path, "/user-agent"):
		body := request.headers["User-Agent"]
		response.body = body
		response.headers = map[string]string{
			"Content-Type":   "text/plain",
			"Content-Length": fmt.Sprint(len(body)),
		}
	default:
		response.status.code = 404
		response.status.message = "Not Found"
	}

	conn.Write([]byte(response.String()))
}
