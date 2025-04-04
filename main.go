package main

import (
	"log"
	"net/http"

	"github.com/CodexRodney/WeStreamBackend/cmd"
)

func main() {
	router := cmd.NewRouter()
	srv := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	err := srv.ListenAndServe()

	if err != nil {
		log.Fatal(err)
	}
}
