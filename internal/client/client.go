package client

import (
	"math/rand"
	"net"

	i "github.com/regmicmahesh/term-chat/internal/common"
)

type Client struct {
	username string
	conn     net.Conn
	ipAddr   string
}

func NewClient(conn net.Conn, ipAddr string) *Client {
	randomUsername := rand.Intn(len(i.USERNAMES))
	username := i.USERNAMES[randomUsername]
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
