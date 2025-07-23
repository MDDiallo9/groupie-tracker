package main

import (
	"encoding/json"
	"net/http"
	"fmt"
	"io"
)

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

func (app *application) GetArtist(id string) Artist {

    var artist Artist

    link := "https://groupietrackers.herokuapp.com/api/artists/" + id
    response, err := http.Get(link)

    if err != nil {
        fmt.Print(err.Error())
        app.errorLog.Printf("%v",err)
    }
    defer response.Body.Close()

    responseData, err := io.ReadAll(response.Body)
    if err != nil {
        app.errorLog.Printf("%v",err)
    }

    err = json.Unmarshal(responseData, &artist)
    if err != nil {
        app.errorLog.Printf("%v",err)
    }

	artist.Locations = app.GetLocations(artist)
	artist.ConcertDates = app.GetConcertDates(artist)
	artist.Relations = app.GetRelations(artist)

    return artist
}

func (app *application) GetLocations(artist Artist) []string {

    var locations LocationData

    response, err := http.Get(artist.LocationsLink)

    if err != nil {
        fmt.Print(err.Error())
        app.errorLog.Printf("%v",err)
    }
    defer response.Body.Close()

    responseData, err := io.ReadAll(response.Body)
    if err != nil {
        app.errorLog.Printf("%v",err)
    }

    err = json.Unmarshal(responseData, &locations)
    if err != nil {
        app.errorLog.Printf("%v",err)
    }
    return locations.Locations
}

func (app *application) GetConcertDates(artist Artist) []string {
    var concertDates ConcertDatesData
    response, err := http.Get(artist.ConcertDatesLink)

    if err != nil {
        fmt.Print(err.Error())
        app.errorLog.Printf("%v",err)
    }
    defer response.Body.Close()

    responseData, err := io.ReadAll(response.Body)
    if err != nil {
        app.errorLog.Printf("%v",err)
    }

    err = json.Unmarshal(responseData, &concertDates)
    if err != nil {
        app.errorLog.Printf("%v",err)
    }
    return concertDates.Dates
}

func (app *application) GetRelations(artist Artist) map[string][]string {
    var relations RelationData
    response, err := http.Get(artist.RelationsLink)

    if err != nil {
        fmt.Print(err.Error())
        app.errorLog.Printf("%v",err)
    }
    defer response.Body.Close()

    responseData, err := io.ReadAll(response.Body)
    if err != nil {
        app.errorLog.Printf("%v",err)
    }

    err = json.Unmarshal(responseData, &relations)
    if err != nil {
        app.errorLog.Printf("%v",err)
    }
    return relations.Relations
}