package main

import (
	"fmt"
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
	defer connection.Close()

	// Read incoming data
	buffer := make([]byte, 1024)

	for {

	}
}

type Response struct {
	Method string `json:"method"`
	Prime  bool   `json:"prime"`
}

type Request struct {
	Method string `json:"method"`
	Number int    `json:"number"`
}
