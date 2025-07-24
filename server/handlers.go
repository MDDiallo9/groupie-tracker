package server

import (
	"groupie-tracker/api"
	"html/template"
	"log"
	"net/http"
	"strconv"
	/* "strconv" */)

func home(w http.ResponseWriter, r *http.Request) {
	/*link := "https://groupietrackers.herokuapp.com/api/artists"
	response, err := http.Get(link)
	if err != nil {
		log.Print(err.Error())
	}

	var artists []api.Artist

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
	}*/
	artists := api.GetArtists()

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

	artists := api.GetArtists()
	artist := artists[id+1]

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
