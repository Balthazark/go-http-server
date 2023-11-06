package main

import (
	"fmt"
	"net"
    "errors"
	//"net"
	//"net/http"
	"os"
    "sync"
    "strconv"
)



type clientManager struct {
    mutex sync.Mutex
    clientQueue []net.Conn
    clients map[net.Conn]bool
}


//TODO possibly look into status codes
func addClient(connection net.Conn, clientManager *clientManager)(string){
    clientManager.mutex.Lock()
    defer clientManager.mutex.Unlock()
    if (len(clientManager.clients) >= 10){
        clientManager.clientQueue = append(clientManager.clientQueue, connection)
        return "Connection pool full, you are in queue" + strconv.Itoa(len(clientManager.clientQueue))
    }
    //Check if user entered the queue
    if (clientManager.clients[connection] != false) {
        return "Already in queue"
    }
    clientManager.clients[connection] = true
}



func main() {

    if (len(os.Args) != 2){
        //TODO better message
        fmt.Printf("No port specified, please rerun the server with <port> argument")
        return
    }
    port := os.Args[1]

    listener, err := net.Listen("tcp", ":"+port)
    
    if err != nil {
        fmt.Println("Error starting server on port:", port, "\n", err)
    }

    //Cleanup, close connection when main returns
    defer listener.Close()


    //Init empty clientManager
    clientManager := clientManager{clients: map[net.Conn]bool{}}

    fmt.Println("Server up and running and listening on port", port)

    //Loop for handling incoming requests
    for {
        conn, err := listener.Accept()
        if err != nil {
            fmt.Println("Error on accepting your connection")
        }
        //clientManager.

    }


}

