package services

import (
	"bufio"
	"fmt"

	"github.com/regmicmahesh/term-chat/internal/common"
)

func msgCreator(username string, message string) string {
	return fmt.Sprintf("%s: %s\n", username, message)
}

func BroadcastMessage(sender string, message string, clients []common.ClientInterface) {
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

func BroadcastServerMessage(message string, clients []common.ClientInterface) {
	BroadcastMessage("Server", message, clients)
}

func PrivateServerMessage(message string, client common.ClientInterface) {
	BroadcastServerMessage(message, []common.ClientInterface{client})
}
