package server

import (
	"bufio"
	"errors"
	"fmt"
	"net"
	"time"

	cl "github.com/regmicmahesh/term-chat/internal/client"
	"github.com/regmicmahesh/term-chat/internal/common"
	"github.com/regmicmahesh/term-chat/internal/services"
	"github.com/regmicmahesh/term-chat/internal/utils"
)

var UserExistsError = errors.New("User already exists")

type User struct {
	username string
	password string
}

type Server struct {
	Clients         []common.ClientInterface
	CommandHandler  common.CommandHandlerInterface
	RegisteredUsers []*User
}

func NewServer() *Server {
	return &Server{
		Clients:        make([]common.ClientInterface, 0),
		CommandHandler: nil,
	}
}

func (s *Server) IsUserCredentialsValid(username string, password string) bool {
	for _, user := range s.RegisteredUsers {
		if user.username == username && user.password == password {
			return true
		}
	}
	return false
}

func (s *Server) IsUserRegistered(username string) bool {
	for _, user := range s.RegisteredUsers {
		if user.username == username {
			return true
		}
	}
	return false
}

func (s *Server) RegisterUser(username string, password string) error {
	if s.IsUserRegistered(username) {
		return UserExistsError
	}
	s.RegisteredUsers = append(s.RegisteredUsers, &User{username, password})
	return nil
}

func (s *Server) RegisterCommandHandler(c common.CommandHandlerInterface) {
	s.CommandHandler = c.InitCommandHandler()
}

func (s *Server) GetNumberOfUsers() int {
	return len(s.Clients)
}

func (s *Server) broadcastMessage(sender common.ClientInterface, message string) {
	var clients []services.BroadcastableClient = make([]services.BroadcastableClient, 0)
	for _, client := range s.Clients {
		clients = append(clients, client)
	}
	services.BroadcastMessage(sender.GetUsername(), message, clients)
}

func (s *Server) BroadcastServerMessage(message string) {

	var clients []services.BroadcastableClient = make([]services.BroadcastableClient, 0)
	for _, client := range s.Clients {
		clients = append(clients, client)
	}
	services.BroadcastServerMessage(message, clients)
}

func (s *Server) GetClientByUsername(username string) common.ClientInterface {
	for _, client := range s.Clients {
		if client.GetUsername() == username {
			return client
		}
	}
	return nil
}

func (s *Server) getClientByIP(ipAddr string) common.ClientInterface {

	for _, client := range s.Clients {
		if client.GetIPAddr() == ipAddr {
			return client
		}
	}
	return nil
}

func (s *Server) AddClient(client common.ClientInterface) {

	for _, c := range s.Clients {
		if c.GetIPAddr() == client.GetIPAddr() {
			return
		}
	}

	s.Clients = append(s.Clients, client)
}

func (s *Server) SendServerPrivateMessage(message string, client common.ClientInterface) {
	services.PrivateServerMessage(message, client)
}

func (s *Server) RemoveClient(client common.ClientInterface) {
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

func (s *Server) HandleConn(conn net.Conn) {
	defer conn.Close()

	ipAddr := conn.RemoteAddr().String()

	for {
		message, err := utils.ReadMessage(conn)
		if err != nil {
			return
		}

		client := s.getClientByIP(ipAddr)
		if client == nil {
			client = cl.NewClient(conn, ipAddr)
			s.SendServerPrivateMessage(fmt.Sprintf("/join <username> to continue."), client)
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
