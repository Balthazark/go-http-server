package main

import (
	"bufio"
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"
	"strings"
)

func worker(conn net.Conn) {
	handleRequest(conn)
	conn.Close()
}

// Main request handler for incoming HTTP requests
func handleRequest(conn net.Conn) {
	// Implement your request handling logic here.
	// For GET and POST methods, read the request, process it, and send a response.
	// For other methods, send a "Not Implemented" (501) error response.
	// Create a bufio.Reader to read the HTTP request
	reader := bufio.NewReader(conn)
	req, err := http.ReadRequest(reader)

	if err != nil {
		fmt.Println("Error reading request: ", err)
		return
	}

	switch req.Method {
	case http.MethodGet:
		handleServeFile(conn, req.URL.Path)
	case http.MethodPost:
	default:
	}
}

func handleServeFile(conn net.Conn, path string) {
	contentType, err := getContentType(path)

	if err != nil {
		response := http.Response{
			Status:     "400 Bad Request",
			StatusCode: http.StatusBadRequest,
			Proto:      "HTTP/1.1",
			ProtoMajor: 1,
			ProtoMinor: 1,
		}
		response.Write(conn)
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

	// Write the response status and headers to the connection
	response.Write(conn)
	// Close the file after sending it as the response body
	file.Close()
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

	maxConcurrentRequests := 10
	//Set up a chanel to cap the maxmimum amount of child proceses to 10
	requestChannel := make(chan net.Conn, maxConcurrentRequests)

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}
		requestChannel <- conn
		go worker(conn)

	}

}
