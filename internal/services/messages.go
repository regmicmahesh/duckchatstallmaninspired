package services

import (
	"bufio"
	"fmt"
	"net"
)

func msgCreator(username string, message string) string {
	return fmt.Sprintf("%s: %s\n", username, message)
}

type BroadcastableClient interface {
	GetUsername() string
	GetConnection() net.Conn
}

func BroadcastMessage(sender string, message string, clients []BroadcastableClient) {
	if sender == "" {
		return
	}
	for _, client := range clients {
		if client.GetUsername() == sender {
			continue
		}
		writer := bufio.NewWriter(client.GetConnection())
		writer.WriteString(msgCreator(sender, message))
		writer.Flush()
	}
}

func BroadcastServerMessage(message string, clients []BroadcastableClient) {
	BroadcastMessage("Server", message, clients)
}

func PrivateServerMessage(message string, client BroadcastableClient) {
	BroadcastServerMessage(message, []BroadcastableClient{client})
}
