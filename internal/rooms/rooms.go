package rooms

import (
	"bufio"
	"fmt"
	"log"
	"net"
)

var Clients []net.Conn

func HandleClientConn(conn net.Conn) {
	reader := bufio.NewReader(conn)
	for {
		message, err := reader.ReadBytes('/')
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Message received is:", string(message))
		// broad cast messages to all clients decouple this
		for i := range Clients {
			Clients[i].Write(message)
		}
	}
}
