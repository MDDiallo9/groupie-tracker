package server

import (
	"encoding/json"
	"groupie-tracker/api"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type Suggestion struct {
	Texte string // ce qui est envoyé en résultat de recherche.
	Label string // Ce qui est écrit à côté pour préciser l'origine de la donnée.
}

func home(w http.ResponseWriter, r *http.Request, init *AppData) {
	var filter api.Filter
	filter.CreationDate = []int{1950, 2025}   // Définition du Min-Max.
	filter.FirstAlbumDate = []int{1950, 2020} // Définition du Min-Max.

	if r.Method == "POST" {

		r.ParseForm()

		// Min-Max choisies des dates de sorties du Premier Album par l'utilisateur.
		minFAD, _ := strconv.Atoi(r.Form.Get("minfAd"))
		maxFAD, _ := strconv.Atoi(r.Form.Get("maxfAd"))

		// Min-Max choisies des dates de foncdations du Groupe de Musique par l'utilisateur.
		minCD, _ := strconv.Atoi(r.Form.Get("minCD"))
		maxCD, _ := strconv.Atoi(r.Form.Get("maxCD"))

		// Map pour le filtre de cases à cocher pour le nombre de membres du groupe.
		members := map[int]bool{
			1: false, 2: false, 3: false, 4: false, 5: false, 6: false, 7: false, // Si true, alors sélectionné par l'utilisateur.
		}
		for _, num := range r.Form["members"] { // "r.Form["members"]" contient les valeurs envoyées par le formulaire.
			n, _ := strconv.Atoi(num) // conversion de la réponse en INT.
			members[n] = true
		}

		// Structure du filtre.
		filter = api.Filter{
			Location:       r.Form.Get("Location"),
			FirstAlbumDate: []int{minFAD, maxFAD},
			Members:        members,
			CreationDate:   []int{minCD, maxCD},
		}
	}

	// Structure pour gérer les suggestions et leurs apparitions sur le page HTML.
	data := struct {
		Suggestions []Suggestion
		Artists     []api.Artist
		Playlist20 []api.Artist
		Location []api.Artist
	}{
		Suggestions: SuggestionsGeneration(),
		Artists:     init.Artists,
		Playlist20: api.FilterBy(init.Artists, api.Filter{CreationDate: []int{2000, 2010}}),
		Location: api.FilterBy(init.Artists, api.Filter{Location: "france",CreationDate: []int{1950, 2025}}),
	}

	// Vérifie si il y a une demande du client pour obtenir des données au travers du filtre.
	if isFilterFilled(filter) {
		data.Artists = api.FilterBy(init.Artists, filter)
		log.Print(filter)
	}

	// Appel des pages HTML pour afficher les informations.
	ts, err := template.ParseFiles("./templates/home.html", "./templates/partials/base.html", "./templates/partials/footer.html", "./templates/partials/head.html")
	if err != nil {
		log.Print(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Chargement de la base.html avec toutes les options précédentes.
	err = ts.ExecuteTemplate(w, "base.html", data)
	if err != nil {
		log.Print(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

// Fonction qui gère l'affichage des informations pour la page de chaque artiste.
func Artist(w http.ResponseWriter, r *http.Request, init *AppData) {
	idstring := r.URL.Query().Get("id") // Récupération de l'ID du groupe de musique sur lequel l'utilisateur veut aller.
	id, err := strconv.Atoi(idstring)

	/* if err != nil || id < 1 || id > len(artists) {
		http.Error(w, "Invalid artist ID", http.StatusBadRequest)
		return
	} */

	artist := init.Artists[id-1]

	// Appel des pages HTML pour afficher les informations.
	ts, err := template.ParseFiles("./templates/artist.html", "./templates/partials/base.html", "./templates/partials/footer.html", "./templates/partials/head.html")
	if err != nil {
		log.Print(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Chargement de la base.html avec toutes les options précédentes.
	err = ts.ExecuteTemplate(w, "base.html", artist)
	if err != nil {
		log.Print(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func IndexPage(w http.ResponseWriter, r *http.Request,init *AppData) {
	data := struct {
		Artists []api.Artist
	}{
		Artists: api.FilterBy(artists, api.Filter{CreationDate: []int{2000, 2010}})[8:],
	}

	ts, err := template.ParseFiles("./templates/index.html", "./templates/partials/base.html", "./templates/partials/footer.html", "./templates/partials/head.html")
	if err != nil {
		log.Print(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	err = ts.ExecuteTemplate(w, "base.html", data)
	if err != nil {
		log.Print(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func ArtistMap(w http.ResponseWriter, r *http.Request, init *AppData) {
	idstring := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idstring)
	if err != nil || id < 1 {
		http.Error(w, "ID invalide", http.StatusBadRequest)
		return
	}

	artist := init.Artists[id-1]
	artist.TabCoords = GenerateCoordinates(artist)

	coordsJSON, err := json.Marshal(artist.TabCoords)
	if err != nil {
		log.Print("Erreur JSON:", err)
		http.Error(w, "Erreur interne", http.StatusInternalServerError)
		return
	}

	artist.CoordsJSON = template.JS(coordsJSON)

	tmpl := template.Must(template.ParseFiles(
		"./templates/map.html",
		"./templates/partials/base.html",
		"./templates/partials/head.html",
		"./templates/partials/footer.html",
	))

	err = tmpl.ExecuteTemplate(w, "base.html", artist)
	if err != nil {
		log.Print("Erreur template:", err)
		http.Error(w, "Erreur interne", http.StatusInternalServerError)
	}
}

func search(w http.ResponseWriter, r *http.Request) {
	var query string

	if r.Method == "POST" {
		r.ParseForm()
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
			nonUnique = true
		}

		// Boucle de recherche pour les noms des artistes.
		for _, response := range answer.Members {
			if strings.Contains(strings.ToLower(response), query) {
				nonUnique = true
				break
			}
		}

		// Boucle de recherche pour les relations entre dates et lieux de concerts.
		for date, cities := range answer.Relations {
			for _, city := range cities {
				if strings.Contains(strings.ToLower(date), query) || strings.Contains(strings.ToLower(city), query) {
					nonUnique = true
					break
				}
			}
		}
	}

	// si la recherche ne correspond à rien.
	if len(results) == 0 {
		http.Redirect(w,r,"/404",404)
		return
	}

	// Structure anonyme pour gérer l'envoi des suggestions.
	data := struct {
		Suggestions []Suggestion
		Artists     []api.Artist
	}{
		Suggestions: SuggestionsGeneration(),
		Artists:     results,
	}

	// Chargement des pages HTML.
	ts, err := template.ParseFiles("./templates/home.html", "./templates/partials/base.html", "./templates/partials/footer.html", "./templates/partials/head.html")
	if err != nil {
		log.Print("Erreur template:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Exécution du template
	err = ts.ExecuteTemplate(w, "base.html", data)
	if err != nil {
		log.Print("Erreur exécution template:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

// Construction des suggestions
func SuggestionsGeneration() []Suggestion {
	var suggestions []Suggestion

	for _, artist := range api.GetArtists() {

		artist.Locations = api.GetLocations(artist)

		// Suggestions pour : Groupe de Musique, Date de création du Groupe,Date de sortie du Premier Album.
		suggestions = append(suggestions,
			Suggestion{
				Texte: artist.Name,
				Label: "Groupe de musique",
			},
			Suggestion{
				Texte: strconv.Itoa(artist.CreationDate),
				Label: "Date de Fondation du Groupe : " + artist.Name,
			},
			Suggestion{
				Texte: artist.FirstAlbum,
				Label: "Date de sortie du premier Album du groupe " + artist.Name,
			},
		)

		// Suggestion pour : Artistes composant un groupe de musique.
		for _, member := range artist.Members {
			suggestions = append(suggestions, Suggestion{
				Texte: member,
				Label: "Membre du groupe : " + artist.Name,
			})
		}

		// Suggestion pour : les villes des concerts.
		for _, city := range artist.Locations {
			suggestions = append(suggestions, Suggestion{
				Texte: city,
				Label: "Ville où joue le groupe " + artist.Name,
			})
		}
	}
	return suggestions
}

// Vérification des filtres de recherches.
func isFilterFilled(f api.Filter) bool {
	for _, v := range f.Members {
		if v {
			return true
		}
	}
	if f.Location != "" || (f.FirstAlbumDate[0] != 1950 || f.FirstAlbumDate[1] != 2020) {
		return true
	}
	if len(f.CreationDate) == 2 && (f.CreationDate[0] != 1950 || f.CreationDate[1] != 2025) {
		return true
	}
	return false
}
