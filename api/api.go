package api

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"
)

// Structure du json artists et réunissant toutes les composantes des informations des autres fichiers json.
type Artist struct {
	ID               int                 `json:"id"`
	Image            string              `json:"image"` // sous forme d'URL
	Name             string              `json:"name"`
	Members          []string            `json:"members"` // Liste de membres
	CreationDate     int                 `json:"creationDate"`
	FirstAlbum       string              `json:"firstAlbum"`
	LocationsLink    string              `json:"locations"` // sous forme d'URL
	Locations        []string            // Liste de villes
	ConcertDatesLink string              `json:"concertDates"` // sous forme d'URL
	ConcertDates     []string            // Liste de dates
	RelationsLink    string              `json:"relations"` // sous forme d'URL
	Relations        map[string][]string // Assemble les localisations et dates de concerts.
	TabCoords        []Coordinates       // Stock les coordonnées
	CoordsJSON       template.JS
}

// Structure du json locations
type LocationData struct {
	ID               int      `json:"id"`
	Locations        []string `json:"locations"`    // Liste de villes
	ConcertDatesLink string   `json:"concertDates"` // sous forme d'URL
}

// Structure du json dates
type ConcertDatesData struct {
	ID    int      `json:"id"`
	Dates []string `json:"dates"` // Liste de dates
}

// Structure du json relation
type RelationData struct {
	ID        int                 `json:"id"`
	Relations map[string][]string `json:"datesLocations"` // Assemble les localisations et dates de concerts.
}

type Filter struct {
	CreationDate   []int
	FirstAlbumDate []int
	Members        map[int]bool
	Location       string
}

func GetArtists() []Artist {
	// Création d'une variable artists de type []Artist pour pouvoir appeler la structure dans la fonction.
	var artists []Artist

	// Envoi une requête GET pour obtenir les informations du fichier json.
	link := "https://groupietrackers.herokuapp.com/api/artists"
	response, err := http.Get(link)
	if err != nil {
		fmt.Print(err.Error())
	}
	defer response.Body.Close() // Sert à éviter de garder la connexion avec l'API ouvert après usage !

	// Lecture des données obtenues.

	responseData, err := io.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	// Convertit grâce à json.Unmarshal le json en Bytes.
	err = json.Unmarshal(responseData, &artists)
	if err != nil {
		log.Fatal(err)
	}
	return artists
}

func GetLocations(artist Artist) []string {
	// Création d'une variable locations de type LocationData pour pouvoir appeler la structure dans la fonction.
	var locations LocationData

	// Envoi une requête GET pour obtenir les informations du fichier json.
	// L'URL est obtenue à partir de la Structure artist.LocationsLink.
	response, err := http.Get(artist.LocationsLink)
	if err != nil {
		fmt.Print(err.Error())
		fmt.Printf("%v", err)
	}

	defer response.Body.Close() // Sert à éviter de garder la connexion avec l'API ouvert après usage !

	// Lecture des données obtenues.
	responseData, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Printf("%v", err)
	}

	// Convertit grâce à json.Unmarshal le json en Bytes.

	err = json.Unmarshal(responseData, &locations)
	if err != nil {
		fmt.Printf("%v", err)
	}

	return FormatLocations(locations.Locations)
}

func GetConcertDates(artist Artist) []string {
	// Création d'une variable concertDates de type ConcertDatesData pour pouvoir appeler la structure dans la fonction.
	var concertDates ConcertDatesData

	// Envoi une requête GET pour obtenir les informations du fichier json.
	// L'URL est obtenue à partir de la Structure artist.ConcertDatesLink.
	response, err := http.Get(artist.ConcertDatesLink)
	if err != nil {
		fmt.Print(err.Error())
		fmt.Printf("%v", err)
	}
	defer response.Body.Close() // Sert à éviter de garder la connexion avec l'API ouvert après usage !

	// Lecture des données obtenues.
	responseData, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Printf("%v", err)
	}

	// Convertit grâce à json.Unmarshal le json en Bytes.
	err = json.Unmarshal(responseData, &concertDates)
	if err != nil {
		fmt.Printf("%v", err)
	}
	return concertDates.Dates
}

func GetRelations(artist Artist) map[string][]string {
	// Création d'une variable relations de type RelationData pour pouvoir appeler la structure dans la fonction.
	var relations RelationData

	// Envoi une requête GET pour obtenir les informations du fichier json.
	// L'URL est obtenue à partir de la Structure artist.RelationsLink.
	response, err := http.Get(artist.RelationsLink)
	if err != nil {
		fmt.Print(err.Error())
		fmt.Printf("%v", err)
	}
	defer response.Body.Close() // Sert à éviter de garder la connexion avec l'API ouvert après usage !

	// Lecture des données obtenues.

	responseData, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Printf("%v", err)
	}

	// Convertit grâce à json.Unmarshal le json en Bytes.

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
		if len(filter.FirstAlbumDate) == 2 {
			if firstAlbumDate < filter.FirstAlbumDate[0] || firstAlbumDate > filter.FirstAlbumDate[1] {
				match = false
			}
		}

		// Members filter
		if anyMemberSelected(filter.Members) {
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
			if strings.Contains(part, "_") {
				parts[i] = ""
				end := ""
				for strings.Contains(part, "_") {
					part1, part2, _ := strings.Cut(part, "_")
					parts[i] += capitalize(part1) + " "
					part = part2
					end = part2
				}
				parts[i] += capitalize(end)
			} else {
				parts[i] = capitalize(part)
			}
		}
		formatted = append(formatted, strings.Join(parts, ", "))
	}
	return formatted
}

func anyMemberSelected(members map[int]bool) bool {
	for _, v := range members {
		if v {
			return true
		}
	}
	return false
}

// Example: Fetch locations for all artists concurrently
func FetchAllLocations(artists []Artist) [][]string {
	results := make([][]string, len(artists))
	var wg sync.WaitGroup
	for i, artist := range artists {
		wg.Add(1)
		go func(i int, artist Artist) {
			defer wg.Done()
			results[i] = GetLocations(artist)
		}(i, artist)
	}
	wg.Wait()
	return results
}
