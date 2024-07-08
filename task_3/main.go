package main

import (
	"fmt"
	"net"
)

func main() {
	fmt.Println("Problem 2 Protohacker")
	fmt.Println("Starting server on port 15000")

	listener, err := net.Listen("tcp", ":15000")
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

	// Handle connection
}
