package functions

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
)

// #MARK: Struct des Localisations.
type LocationsData struct {
	ID               int      `json:"id"`
	Locations        []string `json:"locations"`
	ConcertDatesLink string   `json:"concertDates"`
}

// Pour gérer l'index du json.
type middlemenLocations struct {
	Index []LocationsData `json: "index"`
}

// #MARK: Données des localisations.
func ObtainLocations(links Links) []LocationsData {

	adress, err := http.Get(links.Locations)
	if err != nil {
		log.Fatal(err)
	}
	defer adress.Body.Close()

	adressBytes, err := io.ReadAll(adress.Body)
	if err != nil {
		log.Fatal(err)
	}

	var locationsData middlemenLocations

	err = json.Unmarshal(adressBytes, &locationsData)
	if err != nil {
		log.Fatal(err)
	}

	return locationsData.Index
}
