package main

import (
	"encoding/json"
	"html/template"
	"io"
	"log"
	"net/http"
	/* "strconv" */)

type Artist struct {
	ID               int      `json:"id"`
	Image            string   `json:"image"`
	Name             string   `json:"name"`
	Members          []string `json:"members"`
	CreationDate     int      `json:"creationDate"`
	FirstAlbum       string   `json:"firstAlbum"`
	LocationsLink    string   `json:"locations"`
	Locations        []string
	ConcertDatesLink string `json:"concertDates"`
	ConcertDates     []string
	RelationsLink    string `json:"relations"`
	Relations        map[string][]string
}

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	link := "https://groupietrackers.herokuapp.com/api/artists"
	response, err := http.Get(link)
	if err != nil {
		log.Print(err.Error())
	}

	var artists []Artist

	if err != nil {
		log.Print(err.Error())
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

	ts, err := template.ParseFiles("./templates/home.html")
	if err != nil {
		log.Print(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	err = ts.Execute(w, artists)
	if err != nil {
		log.Print(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func (app *application) Artist(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")

	artist := app.GetArtist(id)

	ts, err := template.ParseFiles("./templates/artist.html")
	if err != nil {
		log.Print(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	err = ts.Execute(w, artist)
	if err != nil {
		log.Print(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}
