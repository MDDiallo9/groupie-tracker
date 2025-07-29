package api

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

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

type Filter struct {
	CreationDate   []int
	FirstAlbumDate int
	Members        int
	Location       string
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
		fmt.Printf("%v", err)
	}
	defer response.Body.Close()

	responseData, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Printf("%v", err)
	}

	err = json.Unmarshal(responseData, &locations)
	if err != nil {
		fmt.Printf("%v", err)
	}
	return locations.Locations
}

func GetConcertDates(artist Artist) []string {
	var concertDates ConcertDatesData
	response, err := http.Get(artist.ConcertDatesLink)
	if err != nil {
		fmt.Print(err.Error())
		fmt.Printf("%v", err)
	}
	defer response.Body.Close()

	responseData, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Printf("%v", err)
	}

	err = json.Unmarshal(responseData, &concertDates)
	if err != nil {
		fmt.Printf("%v", err)
	}
	return concertDates.Dates
}

func GetRelations(artist Artist) map[string][]string {
	var relations RelationData
	response, err := http.Get(artist.RelationsLink)
	if err != nil {
		fmt.Print(err.Error())
		fmt.Printf("%v", err)
	}
	defer response.Body.Close()

	responseData, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Printf("%v", err)
	}

	err = json.Unmarshal(responseData, &relations)
	if err != nil {
		fmt.Printf("%v", err)
	}
	return relations.Relations
}

func FilterBy(artists []Artist, filter Filter) []Artist {
	var results []Artist
	filter.Location = strings.Replace(filter.Location, "-", "", -1)

	for _, artist := range artists {
		artist.Locations = GetLocations(artist)
		/* artist.ConcertDates = GetConcertDates(artist)
		artist.Relations = GetRelations(artist) */
		firstAlbumDate, _ := strconv.Atoi(artist.FirstAlbum[6:])

		// creationDate filter
		if len(filter.CreationDate) == 2 {
			if artist.CreationDate >= filter.CreationDate[0] && artist.CreationDate <= filter.CreationDate[1] {
				if !containsArtist(results, artist.ID) {
					results = append(results, artist)
				}
			}
		}
		//  firstAlbum
		if filter.FirstAlbumDate != 0 {
			if firstAlbumDate > filter.FirstAlbumDate {
				if !containsArtist(results, artist.ID) {
					results = append(results, artist)
				}
			}
		}
		// Members
		if filter.Members != 0 {
			if len(artist.Members) == filter.Members {
				if !containsArtist(results, artist.ID) {
					results = append(results, artist)
				}
			}
		}
		// Locations
		if len(filter.Location) > 2 {
			for _, location := range artist.Locations {
				if strings.Contains(location, filter.Location) {
					if !containsArtist(results, artist.ID) {
						results = append(results, artist)
					}
					break
				}
			}
		}

	}

	return results
}

func containsArtist(results []Artist, id int) bool {
	for _, a := range results {
		if a.ID == id {
			return true
		}
	}
	return false
}
