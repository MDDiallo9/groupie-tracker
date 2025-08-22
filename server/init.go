package server

import (
	"fmt"
	"groupie-tracker/api"
	/* "strconv" */)

var artists []api.Artist

/* func InitArtists() {
	artists = api.GetArtists()
	for _, artist := range artists {
		artist.Locations = api.GetLocations(artist)
		artist.ConcertDates = api.GetConcertDates(artist)
		artist.Relations = api.GetRelations(artist)
	}
} */

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
