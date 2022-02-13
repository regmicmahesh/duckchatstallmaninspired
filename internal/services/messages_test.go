package services_test

import (
	"bufio"
	"net"
	"testing"

	"github.com/regmicmahesh/term-chat/internal/services"
)

func spawnServer() (net.Conn, net.Conn) {
	r, w := net.Pipe()
	return r, w
}

type MockClient struct {
	Username string
	Conn     net.Conn
}

func (c *MockClient) GetUsername() string {
	return c.Username
}

func (c *MockClient) GetConnection() net.Conn {
	return c.Conn
}

func NewMockClient(username string, conn net.Conn) *MockClient {
	return &MockClient{
		Username: username,
		Conn:     conn,
	}
}

func readMessage(conn net.Conn) string {
	reader := bufio.NewReader(conn)
	message, _ := reader.ReadString('\n')
	return message
}

func TestBroadcastMessage(t *testing.T) {

	t.Run("should receive message", func(t *testing.T) {
		r, w := spawnServer()
		var client = NewMockClient("receiver", w)
		go func() {
			defer w.Close()
			services.BroadcastMessage("sender", "test", []services.BroadcastableClient{client})
		}()

		message := readMessage(r)
		if message != "sender: test\n" {
			t.Errorf("Expected message to be 'sender: test' but got %s", message)
		}
	})

	t.Run("should not receive message", func(t *testing.T) {
		r, w := spawnServer()
		var client = NewMockClient("receiver", w)
		go func() {
			defer w.Close()
			services.BroadcastMessage("receiver", "test", []services.BroadcastableClient{client})
		}()

		message := readMessage(r)
		if message != "" {
			t.Errorf("Expected message to be empty but got %s", message)
		}
	})

	t.Run("should not receive message when sender is empty", func(t *testing.T) {
		r, w := spawnServer()
		var client = NewMockClient("receiver", w)
		go func() {
			defer w.Close()
			services.BroadcastMessage("", "test", []services.BroadcastableClient{client})
		}()

		message := readMessage(r)
		if message != "" {
			t.Errorf("Expected message to be empty but got %s", message)
		}
	})

}

func TestBroadcastServerMessage(t *testing.T) {

	t.Run("should receive message", func(t *testing.T) {
		r, w := spawnServer()
		var client = NewMockClient("receiver", w)
		go func() {
			defer w.Close()
			services.BroadcastServerMessage("test", []services.BroadcastableClient{client})
		}()

		message := readMessage(r)
		if message != "Server: test\n" {
			t.Errorf("Expected message to be 'Server: test' but got %s", message)
		}
	})

}

func TestPrivateServerMessage(t *testing.T) {

	t.Run("should receive message", func(t *testing.T) {
		r, w := spawnServer()
		var client = NewMockClient("receiver", w)
		go func() {
			defer w.Close()
			services.PrivateServerMessage("test", client)
		}()

		message := readMessage(r)
		if message != "Server: test\n" {
			t.Errorf("Expected message to be 'Server: test' but got %s", message)
		}
	})

}
