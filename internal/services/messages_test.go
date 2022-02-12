package services_test

import (
	"net"
)

func spawnServer() (net.Conn, net.Conn) {
	r, w := net.Pipe()
	return r, w
}

