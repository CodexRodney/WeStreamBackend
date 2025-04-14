package cmd

import "github.com/gorilla/mux"

func NewRouter() *mux.Router {
	router := mux.NewRouter()

	// middlewares
	router.Use(LoggingMiddleware)

	router.HandleFunc("/create-room", CreateRoom).Methods("POST")
	router.HandleFunc("/join-room/{room}", JoinRoom)
	router.HandleFunc("/upload-music", UploadMusicFile).Methods("POST")

	return router
}
