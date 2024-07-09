package main

import (
	"fmt"
	"net"
)

type UserConnection struct {
	Username   string
	Connection net.Conn
}

func (c *UserConnection) SendWelcomeMessage() error {
	_, err := c.Connection.Write([]byte(WelcomeMessage + "\n"))
	return err
}

func (c *UserConnection) GetUsername() (string, error) {
	buf := make([]byte, 64)
	n, err := c.Connection.Read(buf)
	if err != nil {
		return "", err
	}

	// TODO: Validate username
	username := string(buf[:n])
	return username, nil
}

type Room struct {
	Users map[string]UserConnection
}

func (r *Room) UserNameExists(userName string) bool {
	_, exists := r.Users[userName]
	return exists
}

func (r *Room) AddUser(userConnection UserConnection) error {
	if r.UserNameExists(userConnection.Username) {
		return fmt.Errorf("username %s already exists", userConnection.Username)
	}

	r.Users[userConnection.Username] = userConnection
	r.SendJoinMessage(userConnection.Username)

	return nil
}

func (r *Room) SendJoinMessage(userName string) {
	r.BroadcastMessage(userName, fmt.Sprintf(JoinMessage, userName))
}

func (r *Room) BroadcastMessage(userName string, message string) {
	for _, userConnection := range r.Users {
		if userConnection.Username != userName {
			userConnection.Connection.Write([]byte(message + "\n"))
		}
	}
}

const WelcomeMessage = "Welcome to budgetchat! What shall I call you?"
const JoinMessage = "* %s has entered the room"
