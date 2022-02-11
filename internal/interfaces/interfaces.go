package interfaces

import "net"

type Context struct {
	Server ServerInterface
	Client ClientInterface
	Args   map[string]string
}

func (c *Context) Get(key string) string {
	return c.Args[key]
}

type CommandHandlerInterface interface {
	InitCommandHandler() CommandHandlerInterface
	RegisterCommand(pattern string, noOfArgs int, handler func(Context))
	Handle(client ClientInterface, command string) bool
}

type ServerInterface interface {
	BroadcastServerMessage(message string)
	GetClientByUsername(username string) ClientInterface
	GetNumberOfUsers() int
	UpdateUserStatus()
	HandleConn(conn net.Conn)
	RemoveClient(client ClientInterface)
	SendServerPrivateMessage(message string, client ClientInterface)
	RegisterCommandHandler(c CommandHandlerInterface)
	AddClient(client interface{})
}

type ClientInterface interface {
	GetUsername() string
	SetUsername(username string)
	GetIPAddr() string
	GetConnection() net.Conn
}
