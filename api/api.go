package api

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
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

func GetArtists() []Artist {

	// Création d'une variable artists de type []Artist pour pouvoir appeler la structure dans la fonction.
	var artists []Artist

	// Envoi une requête GET pour obtenir les informations du fichier json.
	link := "https://groupietrackers.herokuapp.com/api/artists"
	response, err := http.Get(link)

	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
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
	return locations.Locations
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
		log.Fatalf("Échec de la récupération de l'API des Relations : %v", err)
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
