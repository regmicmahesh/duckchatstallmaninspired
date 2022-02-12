package utils_test

import (
	"bufio"
	"net"
	"testing"

	"github.com/regmicmahesh/term-chat/internal/utils"
)

func spawnServer() (net.Conn, net.Conn) {
	r, w := net.Pipe()
	return r, w
}

func writeMessage(c net.Conn, message string) {
	writer := bufio.NewWriter(c)
	writer.WriteString(message)
	writer.Flush()
}

func TestReadMessage(t *testing.T) {

	t.Run("should return hello world", func(t *testing.T) {
		r, w := spawnServer()
		go writeMessage(w, "Hello World\n")

		expected := "Hello World"
		actual, _ := utils.ReadMessage(r)

		if actual != expected {
			t.Errorf("Expected %s but got %s", expected, actual)
		}

	})

	t.Run("should show blank string", func(t *testing.T) {
		r, w := spawnServer()
		go writeMessage(w, "\n")

		expected := ""
		actual, _ := utils.ReadMessage(r)

		if actual != expected {
			t.Errorf("Expected %s but got %s", expected, actual)
		}

	})

	t.Run("show show error", func(t *testing.T) {
		r, w := spawnServer()
		r.Close()
		w.Close()

		_, err := utils.ReadMessage(r)

		if err.Error() != "io: read/write on closed pipe" {
			t.Errorf("Expected %s but got %s", "io: read/write on closed pipe", err.Error())
		}

	})

}
