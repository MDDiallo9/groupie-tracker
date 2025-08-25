package server

import (
	"net/http"
)

// ServeMux permet de gérer plusieurs pages dans le même temps. ex /acceuil, /artistes, /localisations,...
func Routes() *http.ServeMux {

	// http.NewServeMux sert de routeur pour rediriger les requêtes HTTP vers les bons Handlers.
	mux := http.NewServeMux()

	// http.FileServer va stocker le CSS dans la variable, que http.Dir va lui indiquer.
	fileserver := http.FileServer(http.Dir("./static"))

	// mux.Handle permet de charger le CSS sur la page html.
	mux.Handle("/static/", http.StripPrefix("/static", fileserver))


  // InitArtists()
	mux.HandleFunc("/", home)
	mux.HandleFunc("/artist", Artist)
	mux.HandleFunc("/search", search)
	mux.HandleFunc("/index",IndexPage)

	return mux
}
