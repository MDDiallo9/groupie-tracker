package server

import (
	"groupie-tracker/api"
	/* "strconv" */)

var artists []api.Artist

func InitArtists() {
	artists = api.GetArtists()
	for _, artist := range artists {
		artist.Locations = api.GetLocations(artist)
		artist.ConcertDates = api.GetConcertDates(artist)
		artist.Relations = api.GetRelations(artist)
	}
}
