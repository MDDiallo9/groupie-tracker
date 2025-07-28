package server

import (
	"groupie-tracker/api" // Import pour récupérer les structures et intérargir avec les fonctions.
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"
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

func search(w http.ResponseWriter, r *http.Request) {

	var query string

	// Vérification que c'est une méthode POST
	if r.Method == "POST" {
		r.ParseForm() // Mise en forme des données de la requête "r".
		// la requête obtenue est mise en minuscule pour quelle ne soit pas sensible à la casse.
		query = strings.ToLower(r.FormValue("search"))
	} else {
		// Si ce n'est pas la méthode POST.
		http.Error(w, "Erreur d'envoi de données", http.StatusMethodNotAllowed)
		return
	}

	// Vérification de la taille de la recherche.
	if len(query) == 0 {
		http.Error(w, "Veuillez écrire au moins un caractère dans la barre de recherche.", http.StatusMethodNotAllowed)
		return
	}

	var results []api.Artist // Format pour accueillir les multiples informations

	// "answer" parcours chaque sous-ensemble de la structure Artists.
	for _, answer := range api.GetArtists() { // Plutôt que de créer une variable, on exploite directement la struct depuis sa fonction.

		// Boucle de recherche pour la date de création, du premier album et le nom du groupe.
		if query == strconv.Itoa(answer.CreationDate) ||
			strings.Contains(strings.ToLower(answer.Name), query) ||
			strings.Contains(strings.ToLower(answer.FirstAlbum), query) {

			results = append(results, answer)
		}

		// Boucle de recherche pour les noms des artistes.
		for _, response := range answer.Members {
			if strings.Contains(strings.ToLower(response), query) {
				results = append(results, answer)
				break
			}
		}

		// Boucle de recherche pour les dates des concerts.
		for _, response := range answer.ConcertDates {
			if strings.Contains(strings.ToLower(response), query) {
				results = append(results, answer)
				break
			}
		}

		// Boucle de recherche pour les localisations des concerts.
		for _, response := range answer.Locations {
			if strings.Contains(strings.ToLower(response), query) {
				results = append(results, answer)
				break
			}
		}
	}

	tmpl := template.Must(template.ParseFiles("templates/home.html"))
	tmpl.Execute(w, results)
}
