package main

import (
	"bufio"
	"fmt"
	"net"
	"net/http"
	"os"
)

var badRequestResponse = http.Response{
	Status:     "400 Bad Request",
	StatusCode: http.StatusBadRequest,
	Proto:      "HTTP/1.1",
	ProtoMajor: 1,
	ProtoMinor: 1,
}

func handleRequest(conn net.Conn) {
	defer conn.Close()

	reader := bufio.NewReader(conn)
	req, err := http.ReadRequest(reader)

	if err != nil {
		badRequestResponse.Write(conn)
		return
	}

	if req.Method != http.MethodGet {
		response := http.Response{
			Status:     "501 Not Implemented",
			StatusCode: http.StatusNotImplemented,
			Proto:      "HTTP/1.1",
			ProtoMajor: 1,
			ProtoMinor: 1,
		}
		response.Write(conn)
		return
	}

	targetURL := req.URL.Host

	targetConn, err := net.Dial("tcp", targetURL)
	if err != nil {
		badRequestResponse.Write(conn)
		return
	}
	defer targetConn.Close()

	req.Write(targetConn)

	targetReader := bufio.NewReader(targetConn)
	targetResponse, err := http.ReadResponse(targetReader, req)

	if err != nil {
		badRequestResponse.Write(conn)
		return
	}

	targetResponse.Write(conn)
}

func main() {

	if len(os.Args) != 2 {
		fmt.Printf("No port specified, please rerun the proxy with and additional <port> argument")
		return
	}

	port := os.Args[1]

	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		fmt.Println("Error starting proxy on port:", port, "\n", "Error:", err)
	}

	fmt.Println("Proxy started on port:", port)

	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}
		go handleRequest(conn)

	}

}
