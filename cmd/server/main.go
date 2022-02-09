package main

import (
	"fmt"
	"net"
	"strings"
)

var conns map[string]net.Conn

func init() {
	conns = make(map[string]net.Conn)
}

func handleConnection(conn net.Conn) {
	defer conn.Close()


	for {

		conns[conn.RemoteAddr().String()] = conn

		fmt.Println(conn.RemoteAddr().String())
		buffer := make([]byte, 1024)
		conn.Read(buffer)
		msg := string(buffer)

		msgFormat := strings.Split(msg, ":")

		if len(msgFormat) != 2 {
			fmt.Println("Invalid message format")
			return
		} 

		for _, c := range conns {
			if c.RemoteAddr().String() != conn.RemoteAddr().String() {
				c.Write([]byte(msgFormat[0] + ":" + msgFormat[1]))
			}
		}

	}
}

func main() {
	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		panic(err)
	}

	fmt.Println("🚀 Listening on port 8080 🚀")

	for {
		conn, err := ln.Accept()
		if err != nil {
			panic(err)
		}

		go handleConnection(conn)

	}

	defer ln.Close()
}
