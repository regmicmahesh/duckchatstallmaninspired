package main

import (
	"fmt"
	"net"

	"github.com/regmicmahesh/term-chat/internal/server"
)

func main() {
	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		panic(err)
	}
	defer ln.Close()
	server := server.NewServer()

	fmt.Println("🚀 Listening on port 8080 🚀")

	for {
		conn, err := ln.Accept()
		if err != nil {
			panic(err)
		}

		go server.UpdateUserStatus()

		go server.HandleConn(conn)

	}

}
