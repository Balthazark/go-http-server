package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"strings"
)

const maxConcurrentRequests = 10
	//Set up a chanel to cap the maxmimum amount of child proceses to 10
var requestChannel chan int 

var badRequestResponse = http.Response{
	Status:     "400 Bad Request",
	StatusCode: http.StatusBadRequest,
	Proto:      "HTTP/1.1",
	ProtoMajor: 1,
	ProtoMinor: 1,
}

func worker(conn net.Conn) {
	handleRequest(conn)
	defer conn.Close()
}

// Main request handler for incoming HTTP requests
func handleRequest(conn net.Conn) {
	reader := bufio.NewReader(conn)
	req, err := http.ReadRequest(reader)

	//Send error if http request is badly formatted
	if err != nil {
		badRequestResponse.Write(conn)
		return
	}

	switch req.Method {
	case http.MethodGet:
		handleServeFile(conn, req.URL.Path)
	case http.MethodPost:
		handleWriteFile(conn, req)
	default:
		response := http.Response{
			Status:     "501 Not Implemented",
			StatusCode: http.StatusNotImplemented,
			Proto:      "HTTP/1.1",
			ProtoMajor: 1,
			ProtoMinor: 1,
		}
		response.Write(conn)
	}
}

func handleServeFile(conn net.Conn, path string) {
	contentType, err := getContentType(path)

	if err != nil {
		badRequestResponse.Write(conn)
		return
	}

	file, err := os.Open("./content" + "/file" + strings.Replace(path, "/", ".", 1))

	fmt.Println("./content" + "/file" + strings.Replace(path, "/", ".", 1))

	if err != nil {
		response := http.Response{
			Status:     "404 Not Found",
			StatusCode: http.StatusNotFound,
			Proto:      "HTTP/1.1",
			ProtoMajor: 1,
			ProtoMinor: 1,
		}
		response.Write(conn)
		return
	}

	response := http.Response{
		Status:     "200 OK",
		StatusCode: http.StatusOK,
		Proto:      "HTTP/1.1",
		ProtoMajor: 1,
		ProtoMinor: 1,
		Header:     http.Header{"Content-Type": []string{contentType}},
		Body:       file,
	}
	response.Write(conn)
	file.Close()
}

func handleWriteFile(conn net.Conn, req *http.Request) {
	path := req.URL.Path

	contentType, contentError := getContentType(path)
	if contentError != nil {
		badRequestResponse.Write(conn)
		return
	}

	content, readError := io.ReadAll(req.Body)
	if readError != nil {
		badRequestResponse.Write(conn)
		return
	}

	if contentType != req.Header["Content-Type"][0] {
		badRequestResponse.Write(conn)
		return
	}

	fileName := "./content/file" + strings.Replace(path, "/", ".", 1)

	writeError := os.WriteFile(fileName, content, os.ModePerm)
	if writeError != nil {
		response := http.Response{
			Status:     "500 Internal Server Error",
			StatusCode: http.StatusInternalServerError,
			Proto:      "HTTP/1.1",
			ProtoMajor: 1,
			ProtoMinor: 1,
		}
		response.Write(conn)
	}

    response := http.Response{
		Status:     "200 OK",
		StatusCode: http.StatusOK,
		Proto:      "HTTP/1.1",
		ProtoMajor: 1,
		ProtoMinor: 1,
	}
	response.Write(conn)
}

func getContentType(path string) (string, error) {
	switch path {
	case "/html":
		return "text/html", nil
	case "/txt":
		return "text/plain", nil
	case "/gif":
		return "image/gif", nil
	case "/jpeg":
		return "image/jpeg", nil
	case "/jpg":
		return "image/jpeg", nil
	case "/css":
		return "text/css", nil
	default:
		return "", errors.New("TEST")
	}
}

func main() {

	requestChannel = make(chan int, maxConcurrentRequests)

	if len(os.Args) != 2 {
		//TODO better message
		fmt.Printf("No port specified, please rerun the server with <port> argument")
		return
	}

	port := os.Args[1]

	listener, err := net.Listen("tcp", ":"+port)

	if err != nil {
		fmt.Println("Error starting server on port:", port, "\n", "Error:", err)
	}

	//Cleanup, close connection if main were to return
	defer listener.Close()


	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
		}
		
		requestChannel <- 1
		go func (conn net.Conn)  {
			fmt.Println(len(requestChannel))
			worker(conn)
			<- requestChannel
		}(conn)
		
	}

}
