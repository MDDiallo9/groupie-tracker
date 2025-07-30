package server

import (
	"html/template"
	"log"
	"net/http"
	"strconv"

	"groupie-tracker/api"
	/* "strconv" */)

func home(w http.ResponseWriter, r *http.Request) {
	testFilter := api.Filter{ /* CreationDate: []int{1970,2000}, */ /* Members: map[int]bool{3:true,4:true,5:true,6:true,7:true}, */ Location: "USA"}
	artists := api.GetArtists()

	test := api.FilterBy(artists, testFilter)

	ts, err := template.ParseFiles("./templates/home.html", "./templates/partials/base.html", "./templates/partials/footer.html", "./templates/partials/head.html")
	if err != nil {
		log.Print(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	err = ts.ExecuteTemplate(w, "base.html", test)
	if err != nil {
		log.Print(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func Artist(w http.ResponseWriter, r *http.Request) {
	idstring := r.URL.Query().Get("id")
	id, _ := strconv.Atoi(idstring)

	artists := api.GetArtists()
	artist := artists[id-1]

	artist.Locations = api.GetLocations(artist)
	artist.ConcertDates = api.GetConcertDates(artist)
	artist.Relations = api.GetRelations(artist)

	ts, err := template.ParseFiles("./templates/artist.html", "./templates/partials/base.html", "./templates/partials/footer.html", "./templates/partials/head.html")
	if err != nil {
		log.Print(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	err = ts.ExecuteTemplate(w, "base.html", artist)
	if err != nil {
		log.Print(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}
