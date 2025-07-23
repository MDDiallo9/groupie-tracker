package functions

import (
	"log"
	"net/http"
	"os"
)

func ServerGroupieTracker(artists []ArtistsData, locations []LocationsData, dates []ConcertDatesData, relations []RelationsData) {
	// Chargement du template HTML) {

	// si une requête arrive dont le chemin commence par "/", on appellera la fonction (handlerFunc).
	// Si l’URL n’est pas exactement "/", on renvoie l'erreur 404 ( vérifie le chemin pour éviter de traiter autre chose)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		HomeHandler(w, r, artists, locations, dates, relations)
		// Chargement du template HTML)
	})

	// définis un "serveur de fichiers" pointant sur le dossier ./template
	fs := http.FileServer(http.Dir("./statics"))
	http.Handle("/statics/", http.StripPrefix("/statics/", fs)) // StripPrefix enlève le /static/ de l'URL pour retrouver le.css

	// Log du répertoire de travail
	// sert uniquement à indiquer, au lancement, le chemin absolu du dossier depuis lequel le code tourne
	// Renvoie une erreur si le répertoire courant a été supprimé ou est inaccessible.
	if wd, err := os.Getwd(); err == nil {
		log.Println("Répertoire de travail :", wd)
	} else {
		log.Fatal("Erreur pour obtenir le répertoire de travail :", err)
	}

	log.Println("Serveur démarré sur http://localhost:8081")
	if err := http.ListenAndServe(":8081", nil); err != nil {
		log.Fatal("Erreur lors du démarrage du serveur :", err)
	}
}
