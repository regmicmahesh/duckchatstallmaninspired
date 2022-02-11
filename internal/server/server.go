package server

import (
	"bufio"
	"fmt"
	"log"
	"math/rand"
	"net"
	"strings"
	"time"
)

type Server struct {
	Clients []*Client
}

type Client struct {
	username string
	conn     net.Conn
	ipAddr   string
}

func NewServer() *Server {
	return &Server{
		Clients: make([]*Client, 0),
	}
}

func NewClient(conn net.Conn, ipAddr string) *Client {

	randomUsername := rand.Intn(len(USERNAMES))

	username := USERNAMES[randomUsername]

	return &Client{
		conn:     conn,
		ipAddr:   ipAddr,
		username: username,
	}
}

func msgCreator(username string, message string) string {
	return fmt.Sprintf("%s: %s", username, message)
}

func (s *Server) broadcastMessage(sender *Client, message string) {
	for _, client := range s.Clients {
		fmt.Println("Currently broadcasting to:", client.ipAddr)
		if client != sender {
			writer := bufio.NewWriter(client.conn)
			_, err := writer.WriteString(msgCreator(sender.username, message) + "\n")
			writer.Flush()
			if err != nil {
				log.Println("Error writing to client:", err)
			}

		}
	}
}

func (s *Server) BroadcastServerMessage(message string) {
	for _, client := range s.Clients {
		writer := bufio.NewWriter(client.conn)
		writer.WriteString(msgCreator("Server", message) + "\n")
		writer.Flush()
	}
}

func (s *Server) getClientByUsername(username string) *Client {
	for _, client := range s.Clients {
		if client.username == username {
			return client
		}
	}
	return nil
}

func (s *Server) getClientByIP(ipAddr string) *Client {

	for _, client := range s.Clients {
		if client.ipAddr == ipAddr {
			return client
		}
	}
	return nil
}

func (s *Server) RemoveClient(client *Client) {
	for i, c := range s.Clients {
		if c == client {
			s.Clients = append(s.Clients[:i], s.Clients[i+1:]...)
			return
		}
	}
}

func (s *Server) UpdateUserStatus() {

	ticker := time.NewTicker(time.Second * 5)

	for {
		select {
		case <-ticker.C:
			for _, client := range s.Clients {
				writer := bufio.NewWriter(client.conn)
				writer.WriteString("1\n")
				err := writer.Flush()

				if err != nil {
					s.BroadcastServerMessage(fmt.Sprintf("%s disconnected.", client.username))
					s.RemoveClient(client)
				}

			}
		}
	}
}

func (s *Server) SendServerPrivateMessage(message string, client *Client) {
	writer := bufio.NewWriter(client.conn)
	writer.WriteString(msgCreator("Server", message) + "\n")
	writer.Flush()
}

func (s *Server) HandleConn(conn net.Conn) {
	defer conn.Close()

	ipAddr := conn.RemoteAddr().String()

	for {
		reader := bufio.NewReader(conn)
		message, err := reader.ReadString('\n')
		message = strings.TrimSpace(message)
		if err != nil {
			return
		}

		client := s.getClientByIP(ipAddr)
		if client == nil {
			client = NewClient(conn, ipAddr)
			s.SendServerPrivateMessage("You are not in the chat. Enter /join server to join the chat.", client)
		}

		if message[0] == '/' {
			cmd := strings.SplitN(message, " ", 2)
			if len(cmd) == 2 {
				s.handleCommand(client, cmd[0], cmd[1])
			} else {
				s.handleCommand(client, cmd[0], "")
			}

		} else {
			s.broadcastMessage(client, message)
		}

	}

}

func (s *Server) handleCommand(client *Client, cmd string, arg string) {
	switch cmd {
	case "/nick":
		s.BroadcastServerMessage(fmt.Sprintf("%s changed their name to %s", client.username, arg))
		client.username = arg
	case "/join":
		client = NewClient(client.conn, client.ipAddr)
		s.Clients = append(s.Clients, client)
		s.BroadcastServerMessage(fmt.Sprintf("%s joined the chat", client.username))
  case "/users":
    users := len(s.Clients)
    s.BroadcastServerMessage(fmt.Sprintf("connected users :%d" , users))
    
	case "/quit":
		s.BroadcastServerMessage(fmt.Sprintf("%s left the chat", client.username))
		s.RemoveClient(client)
	case "/whisper":
		args := strings.SplitN(arg, " ", 2)
		if len(args) == 2 {
			target := s.getClientByUsername(args[0])
			if target != nil {
				s.SendServerPrivateMessage(fmt.Sprintf("%s whispered to you: %s", client.username, args[1]), target)
			} else {
				s.SendServerPrivateMessage(fmt.Sprintf("%s is not in the chat.", args[0]), client)
			}
		} else {
			s.SendServerPrivateMessage("Usage: /whisper <username> <message>", client)
		}
	default:
		fmt.Println("Invalid command.")
	}
}
