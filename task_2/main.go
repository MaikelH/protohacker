package main

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"io"
	"math"
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

	assets := make(map[int32]int32)
	c := bufio.NewReader(connection)
	buffer := make([]byte, 9)
	// read the full message, or return an error
	for {
		read, err := io.ReadFull(c, buffer)
		if err != nil {
			fmt.Println("Error reading from connection", err.Error())
			return
		}

		// According to the specification all messages must be 9 bytes long
		if read != 9 {
			fmt.Printf("Received %d bytes instead of 9\n", read)
			continue
		}

		messageType := buffer[0]
		// Check if the message is a query or insert
		if messageType == 'I' {
			timestamp := convertNumber(buffer[1:5])
			assets[timestamp] = convertNumber(buffer[5:])
		} else if messageType == 'Q' {
			mintime := convertNumber(buffer[1:5])
			maxtime := convertNumber(buffer[5:])
			fmt.Printf("Querying for assets between %d and %d\n", mintime, maxtime)

			sum := 0
			itemCount := 0
			// TODO: Make this more efficient
			for i := mintime; i <= maxtime; i++ {
				value, ok := assets[i]
				if ok {
					sum += int(value)
					itemCount++
				}
			}
			mean := int32(0)
			if itemCount > 0 {
				mean = int32(math.Round(float64(sum) / float64(itemCount)))
			}

			responseBuffer := make([]byte, 4)
			binary.BigEndian.PutUint32(responseBuffer, uint32(mean))
			_, err = connection.Write(responseBuffer)
			if err != nil {
				fmt.Println("Error writing to connection", err.Error())
				return
			}
		} else {
			// Unknown message type (disconnect)
			fmt.Printf("Received unknown message type: %c\n", buffer[0])
			continue
		}

		// Reset buffer
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
