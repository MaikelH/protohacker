package main

import (
	"bufio"
	"encoding/binary"
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
	defer connection.Close()

	c := bufio.NewReader(connection)
	buffer := make([]byte, 9)
	// read the full message, or return an error
	for {
		read, err := c.Read(buffer)
		if err != nil {
			fmt.Println("Error reading from connection", err.Error())
			return
		}

		// According to the specification all messages must be 9 bytes long
		if read != 9 {
			fmt.Printf("Received %d bytes instead of 9\n", read)
			continue
		}

		// Check if the message is a query or insert
		if buffer[8] == 'I' {
			timestamp := convertNumber(buffer[4:8])
			fmt.Println(timestamp)
		} else if buffer[8] == 'Q' {

		} else {
			// Unknown message type (disconnect)
			fmt.Printf("Received unknown message type: %c\n", buffer[0])
			continue
		}

		buffer = make([]byte, 9)
	}
}

func convertNumber(bytes []byte) int32 {
	// Convert byte slice to uint32
	unsignedInt := binary.BigEndian.Uint32(bytes)

	// Convert uint32 to int32 (assuming the integer is int32)
	signedInt := int32(unsignedInt)

	return signedInt
}
