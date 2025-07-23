package main

import (
	"groupie-tracker/functions"
)

func main() {
	// #MARK: Lancement du serveur

	links := functions.DataAPI()
	artists := functions.ObtainArtists(links)
	locations := functions.ObtainLocations(links)
	dates := functions.ObtainDates(links)
	relations := functions.ObtainRelations(links)
	// firstPageData := functions.ObtainFirstPageData(artists, locations, dates, relations)
	functions.ServerGroupieTracker(artists, locations, dates, relations)

	// Tests
	/* 	fmt.Println(artists)
	   	fmt.Println(locations)
	   	fmt.Println(dates)
	   	fmt.Println(relations) */
}
