package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
)

func main() {
	first := 0
	conn, err := net.Dial("tcp", "localhost:8000")
	if err != nil {
		log.Fatal(err)
	}
	for {
		reader := bufio.NewReader(os.Stdin)
		fmt.Println("Text to send: ")
		text, _ := reader.ReadString('\n')
		text = strings.Trim(text, "\n")
		fmt.Fprintf(conn, text+"\n")
		servReader := bufio.NewReader(conn)
		if first == 0 {
			first++
			go readFromServ(*servReader)
		}
	}
}

func readFromServ(servReader bufio.Reader) {
	for {
		message, err := servReader.ReadBytes('/')
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Message from server: " + string(message))
	}
}
