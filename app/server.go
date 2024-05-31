package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"strings"
)

var directoryFlag = flag.String("directory", ".", "The directory to serve files from")

func main() {
	flag.Parse()

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
		response.body = []byte(echo)
	case strings.HasPrefix(request.line.path, "/user-agent"):
		body := request.headers["User-Agent"]
		response.body = []byte(body)
		response.headers = map[string]string{
			"Content-Type":   "text/plain",
			"Content-Length": fmt.Sprint(len(body)),
		}
	case strings.HasPrefix(request.line.path, "/files"):
		filePath := strings.TrimPrefix(request.line.path, "/files/")
		dirEntries, _ := os.ReadDir(*directoryFlag)

		dirEntry := func() os.DirEntry {
			for _, dirEntry := range dirEntries {
				if dirEntry.Name() == filePath {
					return dirEntry
				}
			}
			return nil
		}()

		if dirEntry == nil {
			response.status.code = 404
			response.status.message = "Not Found"
			break
		}

		file, _ := os.ReadFile(fmt.Sprintf("%s/%s", *directoryFlag, filePath))
		response.body = file
		response.headers = map[string]string{
			"Content-Type":   "application/octet-stream",
			"Content-Length": fmt.Sprint(len(file)),
		}

	default:
		response.status.code = 404
		response.status.message = "Not Found"
	}

	conn.Write(response.Bytes())
}
