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

/*
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
}*/

func home(w http.ResponseWriter, r *http.Request) {
	var filter api.Filter
	if r.Method == "POST" {
		r.ParseForm()
		fad, _ := strconv.Atoi(r.Form["FirstAlbumDate"][0])
		cdMin, _ := strconv.Atoi(r.Form["creationDate"][0])
		cdMax := 2025
		members := map[int]bool{
			1: false,
			2: false,
			3: false,
			4: false,
			5: false,
			6: false,
			7: false,
		}
		for _, num := range r.Form["members"] {
			n, _ := strconv.Atoi(num)
			members[n] = true
		}
		filter = api.Filter{
			Location:       r.Form["Location"][0],
			FirstAlbumDate: fad,
			Members:        members,
			CreationDate:   []int{cdMin, cdMax},
		} // Besoin de recharger home avec le api.FilterBy(artists,filter)
	}

	data := struct {
		Suggestions []Suggestion
		Artists     []api.Artist
	}{
		Suggestions: SuggestionsGeneration(),
		Artists:     api.GetArtists(),
	}

	/* artists := api.GetArtists() */

	if isFilterFilled(filter) {
		artists = api.FilterBy(artists, filter)
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

	artist := artists[id-1]
	var coords []api.Coordinates
	for _, location := range api.GetLocations(artist) {
		coord, err := api.Geocoding(location)
		if err == nil {
			coords = append(coords, coord)
		}
	}
	mapURL, err := api.GenerateMapUrl(coords)
	if err != nil {
		fmt.Println("error when generating map url :", err)
	}
	artist.MapURL = mapURL

	artist.Locations = api.GetLocations(artist)
	artist.ConcertDates = api.GetConcertDates(artist)
	artist.Relations = api.GetRelations(artist)

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

	unique := make(map[string]bool) // Pour éviter les doublons.
	var results []api.Artist        // Format pour accueillir les multiples informations

	// "answer" parcours chaque sous-ensemble de la structure Artists.
	for _, answer := range api.GetArtists() { // Plutôt que de créer une variable, on exploite directement la struct depuis sa fonction.
		answer.Relations = api.GetRelations(answer)

		nonUnique := false

		// Boucle de recherche pour la date de création, du premier album et le nom du groupe.
		if query == strconv.Itoa(answer.CreationDate) ||
			strings.Contains(strings.ToLower(answer.Name), query) ||
			strings.Contains(strings.ToLower(answer.FirstAlbum), query) {
			nonUnique = true
			// answer.Relations = api.GetRelations(answer)
			// results = append(results, answer)
		}

		// Boucle de recherche pour les noms des artistes.
		for _, response := range answer.Members {
			if strings.Contains(strings.ToLower(response), query) {
				// answer.Relations = api.GetRelations(answer)
				nonUnique = true
				// results = append(results, answer)
				break
			}
		}

		// Boucle de recherche pour les relations entre dates et lieux de concerts.
		for date, cities := range answer.Relations {
			for _, city := range cities {
				if strings.Contains(strings.ToLower(date), query) || strings.Contains(strings.ToLower(city), query) {
					// answer.Relations = api.GetRelations(answer)
					nonUnique = true
					// results = append(results, answer)
					break
				}
			}
		}

		// Condition pour éviter les doublons.
		if nonUnique {
			identifiant := strings.ToLower(answer.Name)
			if !unique[identifiant] {
				unique[identifiant] = true
				results = append(results, answer)
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

func SuggestionsGeneration() []Suggestion {

	var suggestions []Suggestion

	for _, artist := range api.GetArtists() {

		artist.Locations = api.GetLocations(artist)

		// Suggestions pour : Groupe de Musique, Date de création du Groupe,Date de sortie du Premier Album.
		suggestions = append(suggestions,
			Suggestion{
				Texte: artist.Name,
				Label: "Groupe de musique"},
			Suggestion{
				Texte: strconv.Itoa(artist.CreationDate),
				Label: "Date de Fondation du Groupe : " + artist.Name},
			Suggestion{
				Texte: artist.FirstAlbum,
				Label: "Date de sortie du premier Album du groupe " + artist.Name},
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
	if f.Location != "" || f.FirstAlbumDate != 0 {
		return true
	}
	if len(f.CreationDate) == 2 && (f.CreationDate[0] != 0 || f.CreationDate[1] != 2025) {
		return true
	}
	return false
}
