package main

import (
	"fmt"
)

func main() {
	fmt.Println("Problem 4 Protohacker")

	server := NewDatabaseServer()
	server.Start()
}
