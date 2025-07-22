package main

import (
	"log"
	"flag"
	"net/http"
	"os"
)

type application struct {
	errorLog *log.Logger
	infoLog  *log.Logger
}

func main() {
	port := flag.String("port",":8000","PORT")

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime)

	flag.Parse()

	app := &application{
		errorLog : errorLog,
		infoLog : infoLog,
	}

	srv := &http.Server{
		Addr: *port,
		ErrorLog: errorLog,
		Handler: app.routes(),
	}

	infoLog.Println("Starting server on http://localhost" + *port)
	err := srv.ListenAndServe()
	errorLog.Fatal(err)

}
