package server

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
)

type Artist struct {
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

type LocationData struct {
	ID               int      `json:"id"`
	Locations        []string `json:"locations"`
	ConcertDatesLink string   `json:"concertDates"`
}

type ConcertDatesData struct {
	ID    int      `json:"id"`
	Dates []string `json:"dates"`
}

type RelationData struct {
	ID        int                 `json:"id"`
	Relations map[string][]string `json:"datesLocations"`
}

func GetArtists() []Artist {

	var artists []Artist

	link := "https://groupietrackers.herokuapp.com/api/artists"
	response, err := http.Get(link)

	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}
	defer response.Body.Close()

	responseData, err := io.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	err = json.Unmarshal(responseData, &artists)
	if err != nil {
		log.Fatal(err)
	}
	return artists
}

func GetLocations(artist Artist) []string {

	var locations LocationData

	response, err := http.Get(artist.LocationsLink)

	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}
	defer response.Body.Close()

	responseData, err := io.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	err = json.Unmarshal(responseData, &locations)
	if err != nil {
		log.Fatal(err)
	}
	return locations.Locations
}

func GetConcertDates(artist Artist) []string {
	var concertDates ConcertDatesData
	response, err := http.Get(artist.ConcertDatesLink)

	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}
	defer response.Body.Close()

	responseData, err := io.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	err = json.Unmarshal(responseData, &concertDates)
	if err != nil {
		log.Fatal(err)
	}
	return concertDates.Dates
}

func GetRelations(artist Artist) map[string][]string {
	var relations RelationData
	response, err := http.Get(artist.RelationsLink)

	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}
	defer response.Body.Close()

	responseData, err := io.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	err = json.Unmarshal(responseData, &relations)
	if err != nil {
		log.Fatal(err)
	}
	return relations.Relations
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	tmpl, err := template.ParseFiles("templates/welcome.html")
	if err != nil {
		http.Error(w, "Error 500 : Internal error", http.StatusInternalServerError)
		log.Printf("Error when loading the template: %v", err)
		return
	}
	err = tmpl.Execute(w, nil)
	if err != nil {
		http.Error(w, "Error 5 : Internal error", http.StatusInternalServerError)
		log.Printf("Error when executing the template: %v", err)
	}
}
