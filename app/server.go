package main

import (
	"bytes"
	"compress/gzip"
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
	request := NewRequest(buff)
	response := Response{status: NewStatus()}

	switch {
	case request.line.path == "/":
		break
	case strings.HasPrefix(request.line.path, "/echo"):
		handleEcho(&request, &response)
	case strings.HasPrefix(request.line.path, "/user-agent"):
		handleUserAgent(&request, &response)
	case strings.HasPrefix(request.line.path, "/files"):
		handleFiles(&request, &response)
	default:
		response.status.code = 404
		response.status.message = "Not Found"
	}

	conn.Write(response.Bytes())
}

func handleEcho(request *Request, response *Response) {
	echo := strings.TrimPrefix(request.line.path, "/echo/")

	encoding := func() string {
		encodings := strings.Split(request.headers["accept-encoding"], ", ")
		for _, encoding := range encodings {
			switch encoding {
			case "gzip":
				return encoding
			default:
			}
		}
		return ""
	}()

	response.headers = map[string]string{
		"Content-Type": "text/plain",
	}

	switch encoding {
	case "gzip":
		var buffer bytes.Buffer
		gzipWriter := gzip.NewWriter(&buffer)
		gzipWriter.Write([]byte(echo))
		gzipWriter.Close()
		response.body = buffer.Bytes()
		response.headers["Content-Encoding"] = "gzip"
	case "":
		response.body = []byte(echo)
	}
	response.headers["Content-Length"] = fmt.Sprint(len(response.body))
}

func handleUserAgent(request *Request, response *Response) {
	body := request.headers["user-agent"]
	response.body = []byte(body)
	response.headers = map[string]string{
		"Content-Type":   "text/plain",
		"Content-Length": fmt.Sprint(len(body)),
	}
}

func handleFiles(request *Request, response *Response) {
	filePath := strings.TrimPrefix(request.line.path, "/files/")
	absolutePath := fmt.Sprintf("%s/%s", *directoryFlag, filePath)
	switch request.line.method {
	case "GET":
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

		file, _ := os.ReadFile(absolutePath)
		response.body = file
		response.headers = map[string]string{
			"Content-Type":   "application/octet-stream",
			"Content-Length": fmt.Sprint(len(file)),
		}
	case "POST":
		file, _ := os.Create(absolutePath)
		file.Write(request.body)
		response.status.code = 201
		response.status.message = "Created"
	default:
		response.status.code = 405
		response.status.message = "Method Not Allowed"
	}
}
