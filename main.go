package main

import (
	"flag"
	"groupie-tracker/server"
	"log"
	"net/http"
	"os"
)

func main() {
	port := flag.String("port", ":8000", "PORT")

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime)

	flag.Parse()

	srv := &http.Server{
		Addr:     *port,
		ErrorLog: errorLog,
		Handler:  server.Routes(),
	}

	infoLog.Println("Starting server on http://localhost" + *port)
	err := srv.ListenAndServe()
	errorLog.Fatal(err)

}
