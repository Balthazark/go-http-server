## Instructions

To run the program make sure go is installed on you machine

To run in development mode 
```bash
 #server
 cd server/
 go run sever.go <port> 
 ```

```bash
 #proxy server
 cd proxy/
 go run proxy.go <port> 
 ```

 To build into binaries 
```bash
 #server
 cd server/
 go build sever.go
 ```

```bash
 #proxy server
 cd proxy/
 go build proxy.go
 ```

To run binaries 

```bash
 #server
 chmod +x server
 ./server <port> 
 ```

 ```bash
 #proxy server
 chmod +x proxy
 ./proxy <port> 
 ```
