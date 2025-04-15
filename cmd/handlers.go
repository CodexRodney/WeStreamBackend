package cmd

import (
	"log"
	"net/http"
	"strconv"

	"github.com/CodexRodney/WeStreamBackend/internal/rooms"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

// Pre-configure the upgrader, which is responsible for upgrading
// an HTTP connection to a WebSocket connection.
var (
	websocketUpgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
)

// CreateRoom is used to create a room for other vibers to join
func CreateRoom(w http.ResponseWriter, r *http.Request) {
	newRoom := rooms.NewRoom()
	rooms.AvailableRooms[newRoom.Id] = newRoom
	response := CreateRoomSerializer{Id: strconv.FormatUint(newRoom.Id, 10)}
	respondWithJSON(w, http.StatusCreated, response)
}

// used to create a viber
func GetViberID(w http.ResponseWriter, r *http.Request) {
	// create viber
	viber := rooms.NewViber()
	// add viber to dangling vibers
	rooms.AddViberToDanglingViber(viber)
	respondWithJSON(
		w,
		http.StatusCreated,
		map[string]string{
			"viber_id": strconv.FormatUint(viber.GetViberID(), 10),
		},
	)
}

// JoinRoom is an HTTP handler that upgrades the HTTP connection to a
// WebSocket connection and add the new client to a room
func JoinRoom(w http.ResponseWriter, r *http.Request) {
	log.Println("New Connection")
	vars := mux.Vars(r)

	viberId, err := strconv.ParseUint(vars["viber_id"], 10, 64)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	viber := rooms.GetViberFromDanglingVibers(viberId)
	if viber == nil {
		respondWithError(w, http.StatusBadRequest, "Viber Doesn't Exist")
		return
	}

	roomId, err := strconv.ParseUint(vars["room_id"], 10, 64)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	room, ok := rooms.AvailableRooms[roomId]
	if !ok {
		respondWithError(w, http.StatusBadRequest, "Room Doesn't Exist")
		return
	}

	isAdmin, err := strconv.ParseBool(vars["is_admin"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	conn, err := websocketUpgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		respondWithError(w, http.StatusInternalServerError, err.Error())
	}
	viber.SetViberUnSetProps(conn, room, isAdmin)
	room.AddViberToRoom(viber)

	go viber.ReadMessages()
	go viber.WriteMessages()
}

func JoinMusicStream(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	viberId, err := strconv.ParseUint(vars["viber_id"], 10, 64)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	viber := rooms.GetViberFromDanglingVibers(viberId)
	if viber == nil {
		respondWithError(w, http.StatusBadRequest, "Viber Doesn't Exist")
		return
	}

	if viber.GetViberRoom() == nil {

		respondWithError(w, http.StatusBadRequest, "Cannot Join Stream Not In A Room")
		return
	}

	log.Print(viber.GetViberRoom().Id)
	conn, err := websocketUpgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		respondWithError(w, http.StatusInternalServerError, err.Error())
	}

	viber.JoinMusicStream(conn)
}

// used to upload a music file to a specific room
func UploadMusicFile(w http.ResponseWriter, r *http.Request) {

	roomId, err := strconv.ParseUint(r.FormValue("room_id"), 10, 64)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	room, ok := rooms.AvailableRooms[roomId]
	if !ok {
		respondWithError(w, http.StatusBadRequest, "Room Doesn't Exist")
		return
	}

	file, _, err := r.FormFile("music_file")
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	defer file.Close()

	// getting songs metadata
	durationSecs, err := strconv.ParseInt(
		r.FormValue("duration_seconds"),
		10,
		64,
	)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	title := r.FormValue("title")
	artist := r.FormValue("artist")
	if title == "" || artist == "" {
		respondWithError(w, http.StatusBadRequest, "Missing title or artist")
		return
	}

	room.AddMusicToRoom(
		rooms.SetMusic(title, durationSecs, artist, file),
	)

	respondWithJSON(w, http.StatusCreated, map[string]string{"Good": "good"})
}
