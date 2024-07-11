package main

import (
	"fmt"
	"net"
)

type DatabaseServer struct {
}

func NewDatabaseServer() *DatabaseServer {
	return &DatabaseServer{}
}

func (s *DatabaseServer) Start() {
	fmt.Println("Starting server on port 15000")

	listener, err := net.ListenPacket("udp", ":15000")
	if err != nil {
		fmt.Println("Error starting UDP server", err.Error())
		return
	}
	defer listener.Close()

	for {
		buffer := make([]byte, 1000)
		n, addr, err := listener.ReadFrom(buffer)
		if err != nil {
			fmt.Println("Error reading from connection", err.Error())
			return
		}

		go s.handleMessage(addr, buffer[:n])
	}
}

func (s *DatabaseServer) handleMessage(addr net.Addr, buffer []byte) {
	fmt.Printf("Received message from %s\n", addr.String())

	if s.containsByte(buffer, 61) {
		s.insertValue(buffer)
	} else {
		value, err := s.retrieveValue(buffer)
		if err != nil {
			fmt.Println("Error retrieving value", err.Error())
			return
		}
		err = s.sendResponse(addr, value)
		if err != nil {
			fmt.Println("Error sending response", err.Error())
		}
	}
}

func (s *DatabaseServer) containsByte(buffer []byte, target byte) bool {
	for _, b := range buffer {
		if b == target {
			return true
		}
	}
	return false
}

func (s *DatabaseServer) insertValue(buffer []byte) {

}

func (s *DatabaseServer) retrieveValue(buffer []byte) (string, error) {
	return "", nil
}

func (s *DatabaseServer) sendResponse(addr net.Addr, value string) error {
	return nil
}
