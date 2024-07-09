package main

import (
	"fmt"
	"net"
)

/*
Problem 3: Protohacker

Protocol notes:
- Each message is a line of text terminated by a newline character (\n or ASCII 10)
- Clients can send multiple messages per connection
- All messages are raw ASCII text
- Trailing whitespace or return character '\r' should be ignored

- After a client connects, the server should send a welcome message: "Welcome to budgetchat! What shall I call you?"
- The first message from the client is the client's username
  - The username is a string of alphanumeric characters (a-z, A-Z, 0-9)
  - The username is at least 16 characters long
  - Duplicate usernames are not allowed
  - If the user requests an illegal name, the server disconnects the client

- After the username is accepted, the server sends a message: "* bob has entered the room" to all other users in the room
- The server sends a list of users in the channel to the new user: "* The room contains: bob, charlie, dave"
- All other messages from the client are chat messages

- All chat messages are broadcast to all users in the room (except the sender)
  - Chat messages send to other users are in the form of: "[sender] message"
  - Chat messages are not sent to the sender
  - Chat messages are not sent to users who have left the room
  - Chat messages are not sent to users who have not yet entered the room
  - Chat messages are not sent to users who have been disconnected
  - Chat messages must allow at least 1000 characters

- If a user disconnects, the server sends a message: "* bob has left the room" to all other users in the room
*/
func main() {
	fmt.Println("Problem 3 Protohacker")
	fmt.Println("Starting server on port 15000")

	listener, err := net.Listen("tcp", ":15000")
	if err != nil {
		fmt.Println("Error starting TCP server", err.Error())
	}

	room := Room{Users: map[string]UserConnection{}}
	for {
		connection, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection", err.Error())
			return
		}

		go handleConnection(connection, &room)
	}
}

func handleConnection(connection net.Conn, room *Room) {
	fmt.Println("Handling new connection ", connection.RemoteAddr())
	userCon := UserConnection{Connection: connection}

	err := userCon.SendWelcomeMessage()
	if err != nil {
		fmt.Println("Error sending welcome message", err.Error())
		return
	}

	userName, err := userCon.GetUsername()
	if err != nil {
		fmt.Println("Error getting username", err.Error())
		return
	}
	userCon.Username = userName

	err := room.AddUser(userCon)
	if err != nil {
		fmt.Println("Error adding user to room", err.Error())
		return
	}

}
