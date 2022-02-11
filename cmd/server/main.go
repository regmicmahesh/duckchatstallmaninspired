package main

import (
	"fmt"
	"net"

	"github.com/regmicmahesh/term-chat/internal/handler"
	"github.com/regmicmahesh/term-chat/internal/server"
)

func main() {
	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		panic(err)
	}
	defer ln.Close()

	srv := server.NewServer()


	handler.NewCommandHandler(srv)

	fmt.Println("🚀 Listening on port 8080 🚀")

	for {
		conn, err := ln.Accept()
		if err != nil {
			panic(err)
		}

		go srv.UpdateUserStatus()

		go srv.HandleConn(conn)

	}

}
