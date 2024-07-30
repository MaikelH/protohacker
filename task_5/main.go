package main

import (
	"bufio"
	"io"
	"log/slog"
	"net"
	"regexp"
	"strings"
)

var regex = regexp.MustCompile(`(?:^|\s)(7[a-zA-Z0-9]{25,34})(?:\s|$)`)

func main() {
	slog.Info("Problem 5 Protohacker")
	slog.Info("Starting server on port 15000")

	listener, err := net.Listen("tcp", ":15000")
	if err != nil {
		slog.Error("Error starting TCP server", "error", err.Error())
	}

	for {
		connection, err := listener.Accept()
		defer connection.Close()
		if err != nil {
			slog.Info("Error accepting connection", "error", err.Error())
			return
		}

		go handleConnection(connection)
	}
}

func handleConnection(connection net.Conn) {
	slog.Info("Handling new connection ", "address", connection.RemoteAddr())

	outgoingConnection, err := net.Dial("tcp", "chat.protohackers.com:16963")
	defer outgoingConnection.Close()
	if err != nil {
		slog.Error("Error connecting to chat server", "error", err.Error())
		return
	}

	go func() {
		replaceAndCopyMessage(connection, outgoingConnection)
	}()

	replaceAndCopyMessage(outgoingConnection, connection)
}

func replaceAndCopyMessage(destination io.Writer, source io.Reader) {
	bufferedReader := bufio.NewReader(source)
	for {
		nextLine, err := bufferedReader.ReadString('\n')
		if err != nil {
			slog.Error("Error reading line", "error", err.Error())
			return
		}

		outputMessage := replaceBoguscoinAddresses(nextLine)

		if _, err := io.WriteString(destination, outputMessage); err != nil {
			slog.Error("Error writing to chat server", "error", err.Error())
			return
		}
	}
}

func replaceBoguscoinAddresses(input string) string {
	// Replace all matches with the specified text
	replacement := "7YWHMfk9JZe0LM0g1ZauHuiSxhI"
	result := regex.ReplaceAllStringFunc(input, func(match string) string {
		// Preserve leading/trailing spaces
		return strings.Replace(match, match[1:len(match)-1], replacement, -1)
	})

	// Ensure the message ends with a newline character
	if !strings.HasSuffix(result, "\n") {
		result += "\n"
	}

	return result
}
