package server

import (
	"groupie-tracker/api"
	"net/http"
)

type AppData struct {
    Artists []api.Artist
}

// ServeMux permet de gérer plusieurs pages dans le même temps. ex /acceuil, /artistes, /localisations,...
func Routes(data *AppData) *http.ServeMux {

	// http.NewServeMux sert de routeur pour rediriger les requêtes HTTP vers les bons Handlers.
	mux := http.NewServeMux()

	// http.FileServer va stocker le CSS dans la variable, que http.Dir va lui indiquer.
	fileserver := http.FileServer(http.Dir("./static"))

	// mux.Handle permet de charger le CSS sur la page html.
	mux.Handle("/static/", http.StripPrefix("/static", fileserver))

 /*  InitArtists() */
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
<<<<<<< Updated upstream
        home(w, r, data)
    })
    mux.HandleFunc("/artist", func(w http.ResponseWriter, r *http.Request) {
        Artist(w, r, data)
    })
    mux.HandleFunc("/search", func(w http.ResponseWriter, r *http.Request) {
        search(w, r)
    })
    mux.HandleFunc("/index", func(w http.ResponseWriter, r *http.Request) {
        IndexPage(w, r)
    })
  
=======
		home(w, r, data)
	})
	mux.HandleFunc("/artist", func(w http.ResponseWriter, r *http.Request) {
		Artist(w, r, data)
	})
	mux.HandleFunc("/search", func(w http.ResponseWriter, r *http.Request) {
		search(w, r)
	})
	mux.HandleFunc("/index", func(w http.ResponseWriter, r *http.Request) {
		IndexPage(w, r,data)
	})
	mux.HandleFunc("/map", func(w http.ResponseWriter, r *http.Request) {
		ArtistMap(w, r, data)
	})
	

>>>>>>> Stashed changes
	return mux
}
