package client_test

import (
	"net"
	"testing"

	"github.com/regmicmahesh/term-chat/internal/client"
)

func TestClient(t *testing.T) {

	t.Run("should create a client", func(t *testing.T) {
		var conn net.Conn = &net.TCPConn{}
		client := client.NewClient(conn, "1.1.1.1")
		if client == nil {
			t.Errorf("expected a client to be created")
		}
	})

	t.Run("should return nil", func(t *testing.T) {
		var conn net.Conn = &net.TCPConn{}
		client := client.NewClient(conn, "")
		if client != nil {
			t.Errorf("expected a nil client")
		}
	})

	t.Run("should return nil", func(t *testing.T) {
		client := client.NewClient(nil, "")
		if client != nil {
			t.Errorf("expected a nil client")
		}
	})

	t.Run("should get/set username", func(t *testing.T) {
		var conn net.Conn = &net.TCPConn{}
		client := client.NewClient(conn, "1.1.1.1")
		client.SetUsername("test")
		if client.GetUsername() != "test" {
			t.Errorf("expected a username")
		}
	})

	t.Run("shoud return conn", func(t *testing.T) {
		var conn net.Conn = &net.TCPConn{}
		client := client.NewClient(conn, "1.1.1.1")
		if client.GetConnection() != conn {
			t.Errorf("expected a conn")
		}
	})

	t.Run("should return ip", func(t *testing.T) {
		var conn net.Conn = &net.TCPConn{}
		client := client.NewClient(conn, "1.1.1.1")
		if client.GetIPAddr() != "1.1.1.1" {
			t.Errorf("expected ip: 1.1.1.1")
		}
	})


}
