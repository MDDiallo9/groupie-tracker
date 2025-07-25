package server

import (
	"groupie-tracker/api" // Import pour récupérer les structures et intérargir avec les fonctions.
	"html/template"
	"log"
	"net/http"
	"strconv"
)

// Function pour afficher la page d'acueil.
func home(w http.ResponseWriter, r *http.Request) {

	// Récupère les résultats de la fonction GetArtists
	artists := api.GetArtists()

	// Charge le template HMTL du site.
	ts, err := template.ParseFiles("./templates/home.html")
	if err != nil {
		log.Print(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// En réponse à la Request "r", nous renvoyons ResponseWriter "w" avec les data "artists"
	err = ts.Execute(w, artists)
	if err != nil {
		log.Print(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func Artist(w http.ResponseWriter, r *http.Request) {

	// le client web envoit une demande d'information "r" pour chaque artiste lié à une id.
	// Elle est réceptionnée avant d'être convertie en int.
	idstring := r.URL.Query().Get("id")
	id, _ := strconv.Atoi(idstring)

	// artists récupère toutes les informations des artists d'api.GetArtists.
	// Puis il ne prend que l'id identifiée.
	// id +1 car cela commence à 1 et pas 0.
	artists := api.GetArtists()
	artist := artists[id+1]

	// Charge le template HMTL des groupes d'artistes.
	ts, err := template.ParseFiles("./templates/artist.html")
	if err != nil {
		log.Print(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// En réponse à la Request "r", nous renvoyons ResponseWriter "w" avec les data "artists".
	// Par rapport à la fonction précédente, c'est seulement pour une seule ID. Un seul groupe.
	err = ts.Execute(w, artist)
	if err != nil {
		log.Print(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}
