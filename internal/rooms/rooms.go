package rooms

import (
	"io"
	"log"
	"math/rand"
	"sync"
)

// NotifyEvent represents an event that contains a reference
// to the client who initiated the event and the message to be notified.
type NotifyEvent struct {
	viber   *Viber
	message string
}

type EventNotify struct {
	eventType string
	viber     *Viber
	payload   interface{}
}

type EventNotifyPayload struct {
	eventType string
	payload   interface{}
}

// Room keeps clients in the same channel together
type Room struct {
	Id     uint64
	vibers VibersList
	// change in future to pointers
	musics []*Music
	sync.RWMutex
	notifyChan      chan NotifyEvent
	musicNotifyChan chan NotifyEvent
	eventNotifyChan chan EventNotify
}

// holds a list of created rooms
var AvailableRooms = make(map[uint64]*Room)

// NewManager creates a new Manager instance, initializes the client list,
// and starts the goroutine responsible for notifying other clients.
func NewRoom() *Room {
	r := &Room{
		Id:              rand.Uint64(),
		vibers:          make(VibersList),
		musics:          make([]*Music, 0),
		notifyChan:      make(chan NotifyEvent),
		musicNotifyChan: make(chan NotifyEvent),
		eventNotifyChan: make(chan EventNotify),
	}
	go r.notifyOtherVibersInRoom()
	go r.streamMusicToAll()
	go r.listenToNotifyevents()
	return r
}

// add music to a room
func (r *Room) AddMusicToRoom(music *Music, viber *Viber) {
	r.Lock()
	defer r.Unlock()
	r.musics = append(r.musics, music)
	// notify other people about music being added
	r.eventNotifyChan <- EventNotify{
		eventType: "added_music",
		viber:     viber,
		payload:   music,
	}

}

// return musics in a room
func (r *Room) GetMusicsFromRoom() []*Music {
	return r.musics
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

func (r *Room) streamMusicToAll() {
	for {
		select {
		case e := <-r.musicNotifyChan:
			log.Println("Got Message: ", e.message)
			// check commands like play pause skip
			if e.message == "play" {
				// looping through all songs
				for _, music := range r.musics {
					log.Print("Gor here")
					buffer := make([]byte, 4096)
					for {
						n, err := music.GetMusicFile().Read(buffer)
						if err == io.EOF {
							break
						}
						if err != nil {
							log.Println(err)
							continue
						}
						if n > 0 {
							// send to all clients read bytes
							for v, c := range r.vibers {
								if c {
									// v.musicConn.WriteMessage(websocket.BinaryMessage, buffer[:n])
									v.musicChan <- buffer[:n]

								}
							}
						}
					}
					// notify users that there is a new song
				}
			}
		}
	}
}

// listens for events such as music being added users joining etc
func (r *Room) listenToNotifyevents() {
	for {
		select {
		case e := <-r.eventNotifyChan:
			if e.eventType == "joining" {
				log.Println("I am in joining")
				// notify other vibers about them joining
				otherVibers := r.otherVibersInRoom(e.viber)
				go notifyOthersVibersAboutViberJoining(otherVibers, e.viber)
				// send viber other vibers in room
				e.viber.eventChan <- EventNotifyPayload{
					eventType: "vibers",
					payload:   otherVibers,
				}

			}

			if e.eventType == "added_music" {
				log.Println("Sending added music to other vibers")
				otherVibers := r.otherVibersInRoom(e.viber)
				go notifyOthersVibersAboutMusicAdded(otherVibers, e.payload.(*Music))
			}

		}
	}
}

// used to notify other vibers about recently added music
func notifyOthersVibersAboutMusicAdded(otherVibers []*Viber, music *Music) {
	log.Println("I am here")
	for _, v := range otherVibers {
		log.Println("Sending in loop")
		v.eventChan <- EventNotifyPayload{
			eventType: "added_music",
			payload:   music.extractMusicInfo(),
		}
	}
}

// used to notify other vibers about users joining
func notifyOthersVibersAboutViberJoining(otherVibers []*Viber, viber *Viber) {
	for _, v := range otherVibers {
		v.eventChan <- EventNotifyPayload{
			eventType: "viber_join",
			payload: map[string]interface{}{
				"username": viber.username,
				"viber_id": viber.id,
				"is_admin": viber.isAdmin,
			},
		}
	}
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
