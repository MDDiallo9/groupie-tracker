package main

import (
    "fmt"
    "io"
    "log"
    "net/http"
    "os"
	"encoding/json"
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

func main() {

    response, err := http.Get("https://groupietrackers.herokuapp.com/api/artists/1")

    if err != nil {
        fmt.Print(err.Error())
        os.Exit(1)
    }

    responseData, err := io.ReadAll(response.Body)
    if err != nil {
        log.Fatal(err)
    }
    
	
	var artist Artist
    err = json.Unmarshal(responseData, &artist)
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("L'artiste #%v s'appelle %v, son premier date du : %v",artist.ID, artist.Name, artist.FirstAlbum)

}