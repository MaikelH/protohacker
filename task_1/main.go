package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"math"
	"net"
	"net/textproto"
)

func main() {
	fmt.Println("Problem 1 Protohacker")
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
	textReader := textproto.NewReader(c)
	// read the full message, or return an error
	for {
		line, err := textReader.ReadLine()
		if err != nil {
			fmt.Println("Error reading body: ", err.Error())
			return
		}

		fmt.Println("Received message: ", line)
		var request Request
		err = json.Unmarshal([]byte(line), &request)
		if err != nil || request.Number == nil || request.Method != "isPrime" {
			response := Response{}
			sendResponse(connection, response)
			return
		}

		isPrime := IsPrime(int64(*request.Number))
		response := Response{
			Method: "isPrime",
		}
		if isPrime {
			fmt.Println("Number is prime")
			response.Prime = true
		}

		sendResponse(connection, response)
	}
}

func sendResponse(connection net.Conn, response Response) {
	bytes, err := json.Marshal(response)
	if err != nil {
		fmt.Println("Error marshalling body", err.Error())
		return
	}

	// Append newline character
	bytes = append(bytes, 10)
	written, err := connection.Write(bytes)
	if err != nil {
		fmt.Println("Error responding to client", err.Error())
		return
	}

	fmt.Println("Written ", written)
}

func IsPrime(number int64) bool {
	if number <= 1 {
		return false
	}
	if number == 2 {
		return true
	}
	if number%2 == 0 {
		return false
	}

	boundary := int(math.Floor(math.Sqrt(float64(number))))

	for i := 3; i <= boundary; i++ {
		if number%int64(i) == 0 {
			return false
		}
	}

	return true
}

type Response struct {
	Method string `json:"method"`
	Prime  bool   `json:"prime"`
}

type Request struct {
	Method string   `json:"method"`
	Number *float64 `json:"number"`
}
