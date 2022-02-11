package server

import (
	"bufio"
	"fmt"
	"net"
	"strings"
	"time"

	cl "github.com/regmicmahesh/term-chat/internal/client"
	i "github.com/regmicmahesh/term-chat/internal/common"
	"github.com/regmicmahesh/term-chat/internal/services"
)

type Server struct {
	Clients        []i.ClientInterface
	CommandHandler i.CommandHandlerInterface
}

func NewServer() *Server {
	return &Server{
		Clients:        make([]i.ClientInterface, 0),
		CommandHandler: nil,
	}
}

func (s *Server) RegisterCommandHandler(c i.CommandHandlerInterface) {
	s.CommandHandler = c.InitCommandHandler()
}

func (s *Server) GetNumberOfUsers() int {
	return len(s.Clients)
}

func (s *Server) broadcastMessage(sender i.ClientInterface, message string) {
	services.BroadcastMessage(sender.GetUsername(), message, s.Clients)
}

func (s *Server) BroadcastServerMessage(message string) {
	services.BroadcastServerMessage(message, s.Clients)
}

func (s *Server) GetClientByUsername(username string) i.ClientInterface {
	for _, client := range s.Clients {
		if client.GetUsername() == username {
			return client
		}
	}
	return nil
}

func (s *Server) getClientByIP(ipAddr string) i.ClientInterface {

	for _, client := range s.Clients {
		if client.GetIPAddr() == ipAddr {
			return client
		}
	}
	return nil
}

func (s *Server) RemoveClient(client i.ClientInterface) {
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
				writer := bufio.NewWriter(client.GetConnection())
				writer.WriteString("1\n")
				err := writer.Flush()

				if err != nil {
					s.BroadcastServerMessage(fmt.Sprintf("%s disconnected.", client.GetUsername()))
					s.RemoveClient(client)
				}

			}
		}
	}
}

func (s *Server) AddClient(client i.ClientInterface) {
	s.Clients = append(s.Clients, client)
}

func (s *Server) SendServerPrivateMessage(message string, client i.ClientInterface) {
	services.PrivateServerMessage(message, client)
}

func readMessage(conn net.Conn) (string, error) {
	reader := bufio.NewReader(conn)
	message, err := reader.ReadString('\n')

	if err != nil {
		return "", err
	}
	return strings.TrimSpace(message), nil
}

func (s *Server) HandleConn(conn net.Conn) {
	defer conn.Close()

	ipAddr := conn.RemoteAddr().String()

	for {
		message, err := readMessage(conn)
		if err != nil {
			return
		}

		client := s.getClientByIP(ipAddr)
		if client == nil {
			client = cl.NewClient(conn, ipAddr)
			s.SendServerPrivateMessage(fmt.Sprintf("%s connected.", client.GetUsername()), client)
			s.CommandHandler.Handle(client, "/join")
		}
		if len(message) == 0 {
			s.SendServerPrivateMessage("Please enter a message.", client)
		}

		if message[0] == '/' {
			if status := s.CommandHandler.Handle(client, message); !status {
				s.SendServerPrivateMessage("Command not found.", client)
			}
		} else {
			s.broadcastMessage(client, message)
		}

	}
}
