package server

import (
	"fmt"
	"groupie-tracker/api"
	"html/template"
	"log"
	"net/http"
	"strconv"
	/* "strconv" */)

func home(w http.ResponseWriter, r *http.Request) {

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

func Artist(w http.ResponseWriter, r *http.Request) {
	idstring := r.URL.Query().Get("id")
	id, _ := strconv.Atoi(idstring)

	artist := artists[id-1]
	var coords []api.Coordinates
	for _, location := range api.GetLocations(artist) {
		coord, err := api.Geocoding(location)
		if err == nil {
			coords = append(coords, coord)
		}
	}
	mapURL, err := api.GenerateMapUrl(coords)
	if err != nil {
		fmt.Println("error when generating map url :", err)
	}
	artist.MapURL = mapURL

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