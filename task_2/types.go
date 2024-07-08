package main

import (
	"encoding/binary"
	"math"
	"net"
)

type ConnectionState struct {
	connection net.Conn
	prices     map[int32]int32
}

// GetMeanPrice returns the mean price of the items in the given time range and this connection state.
func (c *ConnectionState) GetMeanPrice(startTime, endTime int32) int32 {
	sum := 0
	itemCount := 0

	for i := startTime; i <= endTime; i++ {
		value, ok := c.prices[i]
		if ok {
			sum += int(value)
			itemCount++
		}
	}

	mean := int32(0)
	if itemCount > 0 {
		mean = int32(math.Round(float64(sum) / float64(itemCount)))
	}

	return mean
}

func (c *ConnectionState) WriteMean(mean int32) error {
	// Send the mean back to the client
	responseBuffer := make([]byte, 4)
	binary.BigEndian.PutUint32(responseBuffer, uint32(mean))
	_, err := c.connection.Write(responseBuffer)

	return err
}

func (c *ConnectionState) Close() {
	c.connection.Close()
}
