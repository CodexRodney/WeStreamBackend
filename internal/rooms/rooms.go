package rooms

import (
	"math/rand"
	"sync"
)

// NotifyEvent represents an event that contains a reference
// to the client who initiated the event and the message to be notified.
type NotifyEvent struct {
	viber   *Viber
	message string
}

// Room keeps clients in the same channel together
type Room struct {
	Id     uint64
	vibers VibersList
	sync.RWMutex
	notifyChan chan NotifyEvent
}

// holds a list of created rooms
var AvailableRooms = make(map[uint64]*Room)

// NewManager creates a new Manager instance, initializes the client list,
// and starts the goroutine responsible for notifying other clients.
func NewRoom() *Room {
	r := &Room{
		Id:         rand.Uint64(),
		vibers:     make(VibersList),
		notifyChan: make(chan NotifyEvent),
	}
	go r.notifyOtherVibersInRoom()
	return r
}

// otherVibersInRoom returns a slice of vibers in a specific room excluding the provided .
func (r *Room) otherVibersInRoom(viber *Viber) []*Viber {
	viberList := make([]*Viber, 0)
	for v := range r.vibers {
		if v.id != viber.id {
			viberList = append(viberList, v)
		}
	}
	return viberList
}

// notifyOtherVibersInRoom waits for notify events and broadcasts the message
// to all vibers except the one who sent the message.
func (r *Room) notifyOtherVibersInRoom() {
	for {
		select {
		case e := <-r.notifyChan:
			otherClients := r.otherVibersInRoom(e.viber)
			for _, c := range otherClients {
				c.writeChan <- e.message
			}
		}
	}
}

// addViberToRoom adds a new viber to the room's viber list.
func (r *Room) AddViberToRoom(viber *Viber) {
	r.Lock()
	defer r.Unlock()
	r.vibers[viber] = true
}

// removeViberFromRoom removes a viber from the room's viber list and
// closes the WebSocket connection.
func (r *Room) removeViberFromRoom(viber *Viber) {
	r.Lock()
	defer r.Unlock()
	if _, ok := r.vibers[viber]; ok {
		viber.connection.Close()
		delete(r.vibers, viber)
	}
}
