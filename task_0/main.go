package main

import (
	"fmt"
	"io"
	"net"
)

func main() {
	fmt.Println("Starting TCP echo server")
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

	if _, err := io.Copy(connection, connection); err != nil {
		fmt.Println("Error copying data", err.Error())
	}
}
