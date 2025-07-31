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
	Members        map[int]bool
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
	return FormatLocations(locations.Locations)
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
	search := normalize(filter.Location)

	for _, artist := range artists {
		artist.Locations = GetLocations(artist)
		firstAlbumDate, _ := strconv.Atoi(artist.FirstAlbum[6:])
		match := true

		// CreationDate filter
		if len(filter.CreationDate) == 2 {
			if artist.CreationDate < filter.CreationDate[0] || artist.CreationDate > filter.CreationDate[1] {
				match = false
			}
		}

		// FirstAlbum filter
		if filter.FirstAlbumDate != 0 {
			if firstAlbumDate <= filter.FirstAlbumDate {
				match = false
			}
		}

		// Members filter
		if len(filter.Members) > 0 {
			if !filter.Members[len(artist.Members)] {
				match = false
			}
		}

		// Locations filter
		if len(search) > 2 {
			found := false
			for _, location := range artist.Locations {
				if strings.Contains(normalize(location), search) {
					found = true
					break
				}
			}
			if !found {
				match = false
			}
		}

		if match && !containsArtist(results, artist.ID) {
			results = append(results, artist)
		}
	}

	return results
}
// Vérifier si des artistes sont déjà présents dans la slice artists pour éviter les doublons
func containsArtist(results []Artist, id int) bool {
	for _, a := range results {
		if a.ID == id {
			return true
		}
	}
	return false
}

// Pour enlever les caratères spéciaux dans le filtre de recherche
func normalize(s string) string {
	s = strings.ToLower(s)
	s = strings.ReplaceAll(s, " ", "")
	s = strings.ReplaceAll(s, ",", "")
	s = strings.ReplaceAll(s, "-", "")
	return s
}

func capitalize(word string) string {
    if len(word) == 0 {
        return word
    }
    return strings.ToUpper(string(word[0])) + strings.ToLower(word[1:])
}

// Pour formatter "california-usa" en "California, USA"
func FormatLocations(locations []string) []string {
    var formatted []string
    for _, loc := range locations {
        parts := strings.Split(loc, "-")
        for i, part := range parts {
            parts[i] = capitalize(part)
        }
        formatted = append(formatted, strings.Join(parts, ", "))
    }
    return formatted
}