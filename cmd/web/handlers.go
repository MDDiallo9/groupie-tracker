package main

import (
	"net/http"
	"html/template"
	"log"
	"encoding/json"
	"io"
)

type Artist struct {
    ID           int      `json:"id"`
    Image        string   `json:"image"`
    Name         string   `json:"name"`
    Members      []string `json:"members"`
    CreationDate int      `json:"creationDate"`
    FirstAlbum   string   `json:"firstAlbum"`
    Locations    string   `json:"locations"`
    ConcertDates string   `json:"concertDates"`
    Relations    string   `json:"relations"`
}

func (app *application) home(w http.ResponseWriter,r *http.Request){
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
	log.Println(artists)
	err = ts.Execute(w, artists)
	if err != nil {
		log.Print(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}
