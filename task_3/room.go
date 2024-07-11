package main

import (
	"fmt"
	"strings"
)

// Room represents a chat room
type Room struct {
	Users map[string]UserConnection
}

// UserNameExists checks if a username already exists in the room.
func (r *Room) UserNameExists(userName string) bool {
	_, exists := r.Users[userName]
	return exists
}

// AddUser adds a user to the room.
// If the username already exists, an error is returned. Otherwise, the user is added to the room.
// Also sends a message to all users in the room about the new user and the new user a message about the current
// participants in the room.
func (r *Room) AddUser(userConnection UserConnection) error {
	if r.UserNameExists(userConnection.Username) {
		return fmt.Errorf("username %s already exists", userConnection.Username)
	}

	r.SendRoomParticipants(userConnection)
	r.Users[userConnection.Username] = userConnection
	r.SendJoinMessage(userConnection.Username)

	return nil
}

// SendJoinMessage broadcasts a message to all users in the room announcing a new user has joined.
// The message format is defined by the JoinMessage constant.
func (r *Room) SendJoinMessage(userName string) {
	r.BroadcastMessage(userName, fmt.Sprintf(JoinMessage, userName))
}

// BroadcastMessage sends a given message to all users in the room except the sender.
// It appends a newline character to each message before sending.
func (r *Room) BroadcastMessage(userName string, message string) {
	for _, userConnection := range r.Users {
		if userConnection.Username != userName {
			userConnection.Connection.Write([]byte(message + "\n"))
		}
	}
}

// SendRoomParticipants sends a message to a newly connected user listing all current participants in the room.
// The message format is defined by the ParticipantsMessage constant.
func (r *Room) SendRoomParticipants(connection UserConnection) {
	var participants []string
	for userName := range r.Users {
		participants = append(participants, userName)
	}

	connection.Connection.Write([]byte(fmt.Sprintf(ParticipantsMessage, strings.Join(participants, ", ")) + "\n"))
}

// RemoveUser removes a user from the room's Users map and broadcasts a disconnect message.
// The disconnect message format is defined by the DisconnectMessage constant.
func (r *Room) RemoveUser(username string) {
	delete(r.Users, username)
	r.SendDisconnectMessage(username)
}

// SendDisconnectMessage broadcasts a message to all users in the room announcing a user has left.
// The message format is defined by the DisconnectMessage constant.
func (r *Room) SendDisconnectMessage(username string) {
	r.BroadcastMessage(username, fmt.Sprintf(DisconnectMessage, username))
}

const WelcomeMessage = "Welcome to budgetchat! What shall I call you?"
const JoinMessage = "* %s has entered the room"
const ParticipantsMessage = "* The room contains: %s"
const DisconnectMessage = "* %s has left the room"
