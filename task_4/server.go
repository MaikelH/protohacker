package main

import (
	"fmt"
	"net"
	"strings"
)

const Version = "ckv-1.0"

type DatabaseServer struct {
	store map[string]string
}

func NewDatabaseServer() *DatabaseServer {
	return &DatabaseServer{
		store: make(map[string]string, 64),
	}
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

		go s.handleMessage(addr, buffer[:n], listener)
	}
}

func (s *DatabaseServer) handleMessage(addr net.Addr, buffer []byte, con net.PacketConn) {
	fmt.Printf("Received message from %s - %s\n", addr.String(), string(buffer))

	if s.containsByte(buffer, 61) {
		s.insertValue(buffer)
	} else {
		value, err := s.retrieveValue(buffer)
		if err != nil {
			fmt.Println("Error retrieving value", err.Error())
			return
		}
		err = s.sendResponse(addr, value, con)
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
	before, after, found := strings.Cut(string(buffer), "=")
	if !found {
		return
	}

	s.store[before] = after
}

func (s *DatabaseServer) retrieveValue(buffer []byte) (string, error) {
	key := string(buffer)
	if key == "version" {
		return fmt.Sprintf("version=%s", Version), nil
	}

	if value, exists := s.store[key]; exists {
		return fmt.Sprintf("%s=%s", key, value), nil
	}

	return fmt.Sprintf("%s=", key), nil
}

func (s *DatabaseServer) sendResponse(addr net.Addr, value string, con net.PacketConn) error {
	// Send the message
	_, err := con.WriteTo([]byte(value), addr)
	if err != nil {
		fmt.Println("Error sending message:", err)
		return err
	}

	return nil
}
