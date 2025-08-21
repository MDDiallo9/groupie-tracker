package server

import (
	"encoding/json"
	"groupie-tracker/api"
	"html/template"
	"log"
	"net/http"
	"strconv"
	/* "strconv" */)

type Suggestion struct {
	Texte string // ce qui est envoyé en recherche.
	Label string // Ce qui sera envoyé à côté dans la liste.
}

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
	id, err := strconv.Atoi(idstring)
	if err != nil || id < 1 || id > len(artists) {
		http.Error(w, "Invalid artist ID", http.StatusBadRequest)
		return
	}

	artist := artists[id-1]

	artist.Locations = api.GetLocations(artist)
	artist.ConcertDates = api.GetConcertDates(artist)
	artist.Relations = api.GetRelations(artist)

	artist.TabCoords = GenerateCoordinates(artist)

	coordsJSON, err := json.Marshal(artist.TabCoords)
	if err != nil {
		log.Print("Erreur JSON:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	artist.CoordsJSON = template.JS(coordsJSON)

	ts, err := template.ParseFiles(
		"./templates/artist.html",
		"./templates/partials/base.html",
		"./templates/partials/footer.html",
		"./templates/partials/head.html",
	)
	if err != nil {
		log.Print("Erreur template:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	err = ts.ExecuteTemplate(w, "base.html", artist)
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
