package server

import (
	"encoding/json"
	"fmt"
	"groupie-tracker/api"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings" // Import pour récupérer les structures et intérargir avec les fonctions.
)

type Suggestion struct {
	Texte string // ce qui est envoyé en recherche.
	Label string // Ce qui sera envoyé à côté dans la liste.
}

func home(w http.ResponseWriter, r *http.Request) {
	var filter api.Filter
	filter.CreationDate = []int{1950, 2025}
	filter.FirstAlbumDate = []int{1970, 2020}
	if r.Method == "POST" {
		r.ParseForm()
		// First Album Date
		minFAD, _ := strconv.Atoi(r.Form.Get("minfAd"))
		maxFAD, _ := strconv.Atoi(r.Form.Get("maxfAd"))
		// Creation Date
		minCD, _ := strconv.Atoi(r.Form.Get("minCD"))
		maxCD, _ := strconv.Atoi(r.Form.Get("maxCD"))
		members := map[int]bool{
			1: false, 2: false, 3: false, 4: false, 5: false, 6: false, 7: false,
		}
		for _, num := range r.Form["members"] {
			n, _ := strconv.Atoi(num)
			members[n] = true
		}
		filter = api.Filter{
			Location:       r.Form.Get("Location"),
			FirstAlbumDate: []int{minFAD, maxFAD}, // or use both min/max if your filter supports a range
			Members:        members,
			CreationDate:   []int{minCD, maxCD},
		} // Besoin de recharger home avec le api.FilterBy(artists,filter)
	}

	data := struct {
		Suggestions []Suggestion
		Artists     []api.Artist
	}{
		Suggestions: SuggestionsGeneration(),
		Artists:     api.GetArtists(),
	}

	artists := api.GetArtists()

	if isFilterFilled(filter) {
		data.Artists = api.FilterBy(artists, filter)
		log.Print(filter)

	}

	ts, err := template.ParseFiles("./templates/home.html", "./templates/partials/base.html", "./templates/partials/footer.html", "./templates/partials/head.html")
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

func Artist(w http.ResponseWriter, r *http.Request) {
	idstring := r.URL.Query().Get("id")
	id, _ := strconv.Atoi(idstring)
	/*
		if err != nil || id < 1 || id > len(artists) {
			http.Error(w, "Invalid artist ID", http.StatusBadRequest)
			return
		}*/

	artist := api.GetArtists()[id-1]

	artist.Locations = api.GetLocations(artist)
	artist.ConcertDates = api.GetConcertDates(artist)
	artist.Relations = api.GetRelations(artist)

	//artist.TabCoords = GenerateCoordinates(artist)

	coordsJSON, err := json.Marshal(artist.TabCoords)
	if err != nil {
		log.Print("Erreur JSON:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	artist.CoordsJSON = template.JS(coordsJSON)

	ts, err := template.ParseFiles("./templates/artist.html", "./templates/partials/base.html", "./templates/partials/footer.html", "./templates/partials/head.html")
	if err != nil {
		log.Print(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	err = ts.ExecuteTemplate(w, "base.html", artist)
	if err != nil {
		log.Print(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func IndexPage(w http.ResponseWriter, r *http.Request) {
	data := struct {
		Playlist20 []api.Artist
	}{

		Playlist20: api.FilterBy(artists, api.Filter{CreationDate: []int{2000, 2010}})[8:],
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

func ArtistMap(w http.ResponseWriter, r *http.Request) {
	idstring := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idstring)
	if err != nil || id < 1 {
		http.Error(w, "ID invalide", http.StatusBadRequest)
		return
	}

	artist := api.GetArtists()[id-1]
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

		// Boucle de recherche pour les relations entre dates et lieux de concerts.
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

	// Exécution du template
	ts, err := template.ParseFiles("./templates/home.html", "./templates/partials/base.html", "./templates/partials/footer.html", "./templates/partials/head.html")
	if err != nil {
		log.Print("Erreur template:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	err = ts.ExecuteTemplate(w, "base.html", data)
	if err != nil {
		log.Print("Erreur exécution template:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

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

func isFilterFilled(f api.Filter) bool {
	for _, v := range f.Members {
		if v {
			return true
		}
	}
	if f.Location != "" || (f.FirstAlbumDate[0] != 1970 || f.FirstAlbumDate[1] != 2020) {
		return true
	}
	if len(f.CreationDate) == 2 && (f.CreationDate[0] != 1950 || f.CreationDate[1] != 2025) {
		return true
	}
	return false
}
