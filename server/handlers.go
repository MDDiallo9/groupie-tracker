package server

import (
	"fmt"
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
	for i := range artists {
		artists[i].ConcertDates = api.GetConcertDates(artists[i])
		artists[i].Relations = api.GetRelations(artists[i])

		// artists[i].Locations = api.GetLocations(artists[i])
		// Boucle pour charger les localisations, avec mise en forme pour l'affichage.
		locations := api.GetLocations(artists[i])
		for j := range locations {
			locations[j] = strings.ReplaceAll(locations[j], "_", " ")
			locations[j] = strings.ReplaceAll(locations[j], "-", " - ")
		}
		artists[i].Locations = locations
	}
	// Charge le template HMTL du site.
	ts, err := template.ParseFiles("./templates/home.html")
	if err != nil {
		log.Print(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// En réponse à la Request "r", nous renvoyons ResponseWriter "w" avec les data "artists"
	ts.Execute(w, artists)

	// J'ai mis cette partie en commentaire et j'ai placé la ligne précédente.
	// Je ne sais plus pourquoi cette partie en commentaire était nécessaire, mais elle génère des erreurs.
	/* 	err = ts.Execute(w, artists)
	   	if err != nil {
	   		log.Print(err.Error())
	   		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	   		return
	   	} */

	// Pour enlever le message d'erreur dans le terminal.
	// http.Error équivaut à response.WriteHeader, et il considére que plusieurs sont actifs si ce dernier return n'est pas là.
	// ERROR   2025/07/29 10:43:06 http: superfluous response.WriteHeader call from groupie-tracker/server.home (handlers.go:37)
	return
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
		return
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

		// Boucle de recherche pour les localisations des concerts.
		for _, response := range answer.Locations {
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

		// Boucle de recherche pour les relations entre dates et lieux de concerts.
		for date, cities := range answer.Relations {
			for _, city := range cities {
				if strings.Contains(strings.ToLower(date), query) || strings.Contains(strings.ToLower(city), query) {
					results = append(results, answer)
					break
				}
			}
		}
	}

	// si la recherche ne correspond à rien.
	if len(results) == 0 {
		fmt.Fprintf(w, "<html><body><p> Nous n'avons pas de correspondances avec votre recherche, veuillez essayer d'autres éléments clés pour tenter de trouver ce que vous voulez. </p></body></html>")
		return
	}

	tmpl := template.Must(template.ParseFiles("templates/home.html"))
	tmpl.Execute(w, results)
}
