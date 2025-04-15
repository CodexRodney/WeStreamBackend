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
	isAdmin    bool
	connection *websocket.Conn
	musicConn  *websocket.Conn
	room       *Room
	writeChan  chan string
	musicChan  chan []byte
}

// return users Id
func (v *Viber) GetViberID() uint64 {
	return v.id
}

func (v *Viber) GetViberRoom() *Room {
	return v.room
}

// ViberList is a map of clients to keep track of their presence.
type VibersList map[*Viber]bool

// Dangling Vibers - Clients who are yet to be added to a room
type danglingViber map[uint64]*Viber

var danglingVibers danglingViber = make(danglingViber)

func AddViberToDanglingViber(viber *Viber) {
	danglingVibers[viber.id] = viber
}

func RemoveViberFromRoom(viber *Viber) {
	delete(danglingVibers, viber.id)
}

func GetViberFromDanglingVibers(viberId uint64) *Viber {
	viber, ok := danglingVibers[viberId]
	if !ok {
		return nil
	}
	return viber
}

func NewViber() *Viber {
	return &Viber{
		id: rand.Uint64(),
	}
}

func (v *Viber) SetViberUnSetProps(
	conn *websocket.Conn,
	room *Room,
	isAdmin bool,
) {
	v.connection = conn
	v.room = room
	v.isAdmin = isAdmin
	v.writeChan = make(chan string)
}

// Setting up clients music web connection
func (v *Viber) JoinMusicStream(conn *websocket.Conn) {
	v.musicConn = conn
	v.musicChan = make(chan []byte)
	// start channels for listening and playing music
	// if not admin start channel for listening only
	if v.isAdmin {
		go v.readRoomMusicCommands()
	}
	go v.streamMusicToClient()
}

// read commands for music stream
func (v *Viber) readRoomMusicCommands() {
	for {
		_, payload, err := v.musicConn.ReadMessage()
		log.Println("Command I got", string(payload))
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error reading message: %v", err)
			}
			break
		}
		v.room.musicNotifyChan <- NotifyEvent{viber: v, message: string(payload)}
	}
}

func (v *Viber) streamMusicToClient() {
	for {
		var dataLen int = 0
		select {
		case data := <-v.musicChan:
			v.musicConn.WriteMessage(websocket.BinaryMessage, data[dataLen:int(len(data))])
			dataLen = len(data)

		}
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
			err := v.connection.WriteMessage(websocket.TextMessage, []byte(data))
			if err != nil {
				log.Println("Error while writing: ", err.Error())
			}
		}
	}
}
