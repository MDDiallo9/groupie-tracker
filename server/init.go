package server

import (
	"fmt"
	"groupie-tracker/api"
	/* "strconv" */)

var artists []api.Artist

// Fonction qui parcours les artistes et les affilies à leurs informations de l'API.Artist.
func InitArtists() []api.Artist {
	artists = api.GetArtists()
	for i := range artists {
		artists[i].Locations = api.GetLocations(artists[i])
		//artists[i].ConcertDates = api.GetConcertDates(artists[i]) */
		artists[i].Relations = api.GetRelations(artists[i])
	}
	return artists
}

// Fonction en charge d'appeler "Geocoding" pour obtenir les coordonnées des concerts du groupe de musique sondé.
func GenerateCoordinates(artist api.Artist) []api.Coordinates {
	var tabcoords []api.Coordinates
	for _, place := range artist.Locations {
		coords, err := api.Geocoding(place)
		if err == nil {
			tabcoords = append(tabcoords, coords)
		} else {
			fmt.Println("Couldn't complete geocoding")
		}
	}
	return tabcoords
}
