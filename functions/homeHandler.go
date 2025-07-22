package functions

import (
	"html/template"
	"log"
	"net/http"
)

// Structure globale pour passer plusieurs types de données au template
type PageData struct {
	Artists   []ArtistsData
	Locations []LocationsData
	Dates     []ConcertDatesData
	Relations []RelationsData
}

func HomeHandler(w http.ResponseWriter, r *http.Request, artists []ArtistsData, locations []LocationsData, dates []ConcertDatesData, relations []RelationsData) {
	// Chargement du template HTML
	tmpl, err := template.ParseFiles("template/presentation.html")
	if err != nil {
		log.Print(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Création de l'objet PageData pour passer toutes les infos au template
	data := PageData{
		Artists:   artists,
		Locations: locations,
		Dates:     dates,
		Relations: relations,
	}

	// Exécution du template avec les données
	err = tmpl.Execute(w, data)
	if err != nil {
		log.Print(err.Error())
		http.Error(w, "Erreur lors de l'affichage de la page", http.StatusInternalServerError)
	}
}
