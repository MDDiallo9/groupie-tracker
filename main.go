package main

import (
	"fmt"
	"groupie-tracker/server"
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", server.HomeHandler)
	fmt.Println("Démarrage du serveur sur le port 4000")
	log.Fatal(http.ListenAndServe(":4000", mux))
}
