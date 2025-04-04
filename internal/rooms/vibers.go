// A viber is basically a normal client in a network
package rooms

import (
	"log"
	"math/rand"

	"github.com/gorilla/websocket"
)

// Viber represents a single WebSocket connection.
// It holds the client's ID, the WebSocket connection itself, and
// the manager that controls all clients.
type Viber struct {
	id         uint64
	connection *websocket.Conn
	room       *Room
	writeChan  chan string
}

// ClientList is a map of clients to keep track of their presence.
type VibersList map[*Viber]bool

// NewClient creates a new Client instance with a unique ID, its connection,
// and a reference to the Manager.
func NewViber(conn *websocket.Conn, room *Room) *Viber {
	return &Viber{
		id:         rand.Uint64(),
		connection: conn,
		room:       room,
		writeChan:  make(chan string),
	}
}

// readMessages continuously reads messages from the WebSocket connection.
// It will send any received messages to the manager's notification channel.
func (v *Viber) ReadMessages() {
	defer func() {
		v.room.removeViberFromRoom(v)
	}()
	for {
		messageType, payload, err := v.connection.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error reading message: %v", err)
			}
			break
		}
		v.room.notifyChan <- NotifyEvent{viber: v, message: string(payload)}
		log.Println("MessageType: ", messageType)
		log.Println("Payload: ", string(payload))
	}
}

// writeMessages listens on the client's write channel for messages
// and writes any received messages to the WebSocket connection.
func (v *Viber) WriteMessages() {
	defer func() {
		v.room.removeViberFromRoom(v)
	}()
	for {
		select {
		case data := <-v.writeChan:
			v.connection.WriteMessage(websocket.TextMessage, []byte(data))
		}
	}
}
