package functions

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
)

// #MARK: Struct des Artistes
type ArtistsData struct {
	ID               int      `json:"id"`
	Image            string   `json:"image"`
	Name             string   `json:"name"`
	Members          []string `json:"members"`
	CreationDate     int      `json:"creationDate"`
	FirstAlbum       string   `json:"firstAlbum"`
	LocationsLink    string   `json:"locations"`
	ConcertDatesLink string   `json:"concertDates"`
	RelationsLink    string   `json:"relations"`
}

// #MARK: Données des Artistes
func ObtainArtists(links Links) []ArtistsData {

	adress, err := http.Get(links.Artists)
	if err != nil {
		log.Fatal(err)
	}
	defer adress.Body.Close()

	adressBytes, err := io.ReadAll(adress.Body)
	if err != nil {
		log.Fatal(err)
	}

	var artistsData []ArtistsData

	err = json.Unmarshal(adressBytes, &artistsData)
	if err != nil {
		log.Fatal(err)
	}

	return artistsData
}
