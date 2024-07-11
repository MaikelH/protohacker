package main

import (
	"bufio"
	"fmt"
	"net"
	"net/textproto"
	"regexp"
	"strings"
)

type UserConnection struct {
	Username       string
	Connection     net.Conn
	isAlphanumeric *regexp.Regexp
}

// NewUserConnection initializes a new UserConnection instance with the provided net.Conn connection.
// It precompiles the regular expression for alphanumeric validation of usernames.
func NewUserConnection(connection net.Conn) UserConnection {
	return UserConnection{
		Connection:     connection,
		isAlphanumeric: regexp.MustCompile(`^[a-zA-Z0-9]+$`),
	}
}

// SendWelcomeMessage sends a welcome message to the user upon connecting to the chat server.
// It appends a newline character to the message for proper formatting.
func (c *UserConnection) SendWelcomeMessage() error {
	_, err := c.Connection.Write([]byte(WelcomeMessage + "\n"))
	return err
}

// GetUsername reads the username from the connection, trims whitespace, and validates it.
// It returns an error if the username is invalid or if there's an issue reading from the connection.
func (c *UserConnection) GetUsername() (string, error) {
	buf := make([]byte, 64)
	n, err := c.Connection.Read(buf)
	if err != nil {
		return "", err
	}

	username := strings.TrimSpace(string(buf[:n]))
	if !c.isValidUsername(username) {
		return "", fmt.Errorf("invalid username")
	}

	return strings.TrimSpace(username), nil
}

// Close closes the user's connection to the server.
func (c *UserConnection) Close() {
	c.Connection.Close()
}

// HandleMessages continuously reads messages from the user's connection,
// broadcasts them to the chat room, and handles disconnection.
// It stops reading and removes the user from the room upon encountering an error.
func (c *UserConnection) HandleMessages(room *Room) {
	bufferedReader := bufio.NewReader(c.Connection)
	textReader := textproto.NewReader(bufferedReader)

	for {
		nextLine, err := textReader.ReadLine()
		if err != nil {
			fmt.Println("Disconnecting: ", c.Username)
			room.RemoveUser(c.Username)
			return
		}
		fmt.Println(nextLine)

		room.BroadcastMessage(c.Username, fmt.Sprintf("[%s] %s", c.Username, nextLine))
	}
}

// isValidUsername checks if the username is valid
//   - The username is a string of alphanumeric characters (a-z, A-Z, 0-9)
//   - The username is at least 1 character long
func (c *UserConnection) isValidUsername(s string) bool {
	if len(s) < 1 {
		return false
	}

	return c.isAlphanumeric.MatchString(s)
}
