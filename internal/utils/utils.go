package utils

import (
	"bufio"
	"net"
	"strings"
)

func ReadMessage(conn net.Conn) (string, error) {
	reader := bufio.NewReader(conn)
	message, err := reader.ReadString('\n')

	if err != nil {
		return "", err
	}
	return strings.TrimSpace(message), nil
}

