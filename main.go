package main

import (
	"fmt"
	"log"
	"net"

	"github.com/CodexRodney/WeStreamBackend/internal/rooms"
)

func main() {
	addr, err := net.ResolveTCPAddr("tcp", ":8000")
	if err != nil {
		log.Fatal(err)
	}
	ln, err := net.ListenTCP("tcp", addr)
	if err != nil {
		log.Fatal(err)
	}
	defer ln.Close()
	fmt.Println("Listening on port 8000")
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Fatal(err)
		}
		// add client to list of clients
		rooms.Clients = append(rooms.Clients, conn)

		go rooms.HandleClientConn(conn)
	}
}
