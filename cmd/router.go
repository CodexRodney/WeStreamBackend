package cmd

import "github.com/gorilla/mux"

func NewRouter() *mux.Router {
	router := mux.NewRouter()

	// middlewares
	router.Use(LoggingMiddleware)

	router.HandleFunc("/create-room", CreateRoom).Methods("POST")
	router.HandleFunc("/get-viber", GetViberID).Methods("GET")
	router.HandleFunc("/join-room/{room_id}/{viber_id}/{is_admin}", JoinRoom)
	router.HandleFunc("/stream-room-music/{viber_id}", JoinMusicStream)
	router.HandleFunc("/upload-music", UploadMusicFile).Methods("POST")

	return router
}
