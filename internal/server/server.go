package server

import (
	"bufio"
	"fmt"
	"log"
	"math/rand"
	"net"
	"strings"
	"time"

	i "github.com/regmicmahesh/term-chat/internal/interfaces"
)

   
type Server struct {
	Clients        []i.ClientInterface
	CommandHandler i.CommandHandlerInterface
}

type Client struct {
	username string
	conn     net.Conn
	ipAddr   string
}


type Context struct {
	Server Server
	Client Client
	args   map[string]string
}

func NewClient(conn net.Conn, ipAddr string) i.ClientInterface {
	randomUsername := rand.Intn(len(USERNAMES))
	username := USERNAMES[randomUsername]
	return &Client{
		conn:     conn,
		ipAddr:   ipAddr,
		username: username,
	}
}

func (c *Client) GetUsername() string {
	return c.username
}

func (c *Client) SetUsername(username string) {
	c.username = username
}

func (c *Client) GetConnection() net.Conn {
	return c.conn
}

func (c *Client) GetIPAddr() string {
	return c.ipAddr
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


func msgCreator(username string, message string) string {
	return fmt.Sprintf("%s: %s", username, message)
}

func (s *Server) broadcastMessage(sender i.ClientInterface, message string) {
	for _, client := range s.Clients {
		fmt.Println("Currently broadcasting to:", client.GetIPAddr())
		if client != sender {
			writer := bufio.NewWriter(client.GetConnection())
			_, err := writer.WriteString(msgCreator(sender.GetUsername(), message) + "\n")
			writer.Flush()
			if err != nil {
				log.Println("Error writing to client:", err)
			}

		}
	}
}

func (s *Server) BroadcastServerMessage(message string) {
	for _, client := range s.Clients {
		writer := bufio.NewWriter(client.GetConnection())
		writer.WriteString(msgCreator("Server", message) + "\n")
		writer.Flush()
	}
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
		if c == client.(*Client) {
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

func (s *Server) AddClient(client interface{}) {
	s.Clients = append(s.Clients, client.(*Client))
}

func (s *Server) SendServerPrivateMessage(message string, client i.ClientInterface) {
	writer := bufio.NewWriter(client.GetConnection())
	writer.WriteString(msgCreator("Server", message) + "\n")
	writer.Flush()
}

func (s *Server) HandleConn(conn net.Conn) {
	defer conn.Close()

	ipAddr := conn.RemoteAddr().String()

	for {
		//TODO: Refactor this to return a string.
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

