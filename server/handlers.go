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

// Structure pour les suggestions de la barre de recherche.
type Suggestion struct {
	Texte string // ce qui est envoyé en recherche.
	Label string // Ce qui sera envoyé à côté dans la liste.
}

// Function pour afficher la page d'acueil.
func home(w http.ResponseWriter, r *http.Request) {

	// Tentative pour gérer la barre de suggestion.
	data := struct {
		Suggestions []Suggestion
		Artists     []api.Artist
	}{
		Suggestions: SuggestionsGeneration(),
		Artists:     api.GetArtists(),
	}

	tmpl := template.Must(template.ParseFiles("templates/home.html"))
	tmpl.Execute(w, data)

	// Récupère les résultats de la fonction GetArtists
	artists := api.GetArtists()
	for i := range artists {

		// Charge les dates des concerts.
		artists[i].ConcertDates = api.GetConcertDates(artists[i])

		// artists[i].Locations = api.GetLocations(artists[i])
		// Boucle pour charger les localisations, avec mise en forme pour l'affichage.
		locations := api.GetLocations(artists[i])
		for j := range locations {
			locations[j] = strings.ReplaceAll(locations[j], "_", " ")
			locations[j] = strings.ReplaceAll(locations[j], "-", " - ")
		}
		artists[i].Locations = locations

		//#MARK: A revoir, ça ne marche pas
		// artists[i].Relations = api.GetRelations(artists[i])
		// Boucle pour charger les relations, avec mise en forme pour l'affichage.
		relations := api.GetRelations(artists[i])

		for date, locations := range relations {
			for j, value := range locations {
				value = strings.ReplaceAll(value, "_", " ")
				value = strings.ReplaceAll(value, "-", " - ")
				locations[j] = value
			}
			relations[date] = locations
		}
		artists[i].Relations = relations
	}

	// Charge le template HMTL du site.
	/* ts, err := template.ParseFiles("./templates/home.html")
	if err != nil {
		log.Print(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// En réponse à la Request "r", nous renvoyons ResponseWriter "w" avec les data "artists"
	ts.Execute(w, artists) */

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

		answer.Relations = api.GetRelations(answer)

		// Boucle de recherche pour la date de création, du premier album et le nom du groupe.
		if query == strconv.Itoa(answer.CreationDate) ||
			strings.Contains(strings.ToLower(answer.Name), query) ||
			strings.Contains(strings.ToLower(answer.FirstAlbum), query) {

			// answer.Relations = api.GetRelations(answer)
			results = append(results, answer)
		}

		// Boucle de recherche pour les noms des artistes.
		for _, response := range answer.Members {
			if strings.Contains(strings.ToLower(response), query) {
				// answer.Relations = api.GetRelations(answer)
				results = append(results, answer)
				break
			}
		}

		/* 		// Boucle de recherche pour les localisations des concerts.
		   		for _, response := range answer.Locations {
		   			if strings.Contains(strings.ToLower(response), query) {
		   				answer.Relations = api.GetRelations(answer)
		   				results = append(results, answer)
		   				break
		   			}
		   		}

		   		// Boucle de recherche pour les dates des concerts.
		   		for _, response := range answer.ConcertDates {
		   			if strings.Contains(strings.ToLower(response), query) {
		   				answer.Relations = api.GetRelations(answer)
		   				results = append(results, answer)
		   				break
		   			}
		   		} */

		// Boucle de recherche pour les relations entre dates et lieux de concerts.
		// Sa fonctionne, donc je mets les recherches sur les dates et lieux individuels en commentaires, car doublons.
		for date, cities := range answer.Relations {
			for _, city := range cities {
				if strings.Contains(strings.ToLower(date), query) || strings.Contains(strings.ToLower(city), query) {
					// answer.Relations = api.GetRelations(answer)
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

	data := struct {
		Suggestions []Suggestion
		Artists     []api.Artist
	}{
		Suggestions: SuggestionsGeneration(),
		Artists:     results,
	}

	tmpl := template.Must(template.ParseFiles("templates/home.html"))
	tmpl.Execute(w, data)
}

func SuggestionsGeneration() []Suggestion {

	var suggestions []Suggestion
	

	for _, artist := range api.GetArtists() {

		artist.Locations = api.GetLocations(artist)

		// Suggestions pour : Groupe de Musique, Date de création du Groupe,Date de sortie du Premier Album.
		suggestions = append(suggestions,
			Suggestion{Texte: artist.Name, Label: "Groupe de musique"},
			Suggestion{Texte: strconv.Itoa(artist.CreationDate), Label: "Date de Fondation du Groupe de musique"},
			Suggestion{Texte: artist.FirstAlbum, Label: "Date de sortie du premier Album"},
		)

		// Suggestion pour : Artistes composant un groupe de musique.
		for _, member := range artist.Members {
			suggestions = append(suggestions, Suggestion{Texte: member, Label: "Membre d'un groupe de musique"})
		}

		// Suggestion pour : les villes des concerts.
		for _, city := range artist.Locations {
			suggestions = append(suggestions, Suggestion{Texte: city, Label: "Ville"})
		}

		/* 		// Suggestion pour : les dates des concerts.
		   		// Non demandé dans les instructions.
		   		for _, dates := range artist.Relations {
		   			for _, date := range dates {
		   				suggestions = append(suggestions, Suggestion{Texte: date, Label: "Date de concert"})
		   			}
		   		} */
	}
	return suggestions
}
