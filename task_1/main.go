package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
)

func main() {
	fmt.Println("Problem 1 Protohacker")
	fmt.Println("Starting server on port 8080")

	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println("Error starting TCP server", err.Error())
	}

	for {
		connection, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection", err.Error())
			return
		}

		go handleConnection(connection)
	}
}

func handleConnection(connection net.Conn) {
	fmt.Println("Handling new connection ", connection.RemoteAddr())
	defer connection.Close()

	c := bufio.NewReader(connection)
	// read the full message, or return an error
	buffer, err := io.ReadAll(c)
	if err != nil {
		fmt.Println("Error reading body")
		return
	}

	fmt.Printf("received %x\n", buffer)
}

type Response struct {
	Method string `json:"method"`
	Prime  bool   `json:"prime"`
}

type Request struct {
	Method string `json:"method"`
	Number int    `json:"number"`
}
