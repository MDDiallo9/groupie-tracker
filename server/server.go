package server

import (
	"net/http"
)

func Routes() *http.ServeMux {
	mux := http.NewServeMux()

	fileserver := http.FileServer(http.Dir("./ui/static"))
	mux.Handle("/static/", http.StripPrefix("/static", fileserver))

	mux.HandleFunc("/{$}", home)
	mux.HandleFunc("/artist", Artist)

	return mux
}
