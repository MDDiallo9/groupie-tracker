package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"groupie-tracker/api"
	/* "strconv" */)

var artists []api.Artist

// Fonction qui parcours les artistes et les affilies à leurs informations de l'API.Artist.
func InitArtists(key string) []api.Artist {
	artists = api.GetArtists()
	for i := range artists {
		// Spotify search
		client := &http.Client{}
		req, err := http.NewRequest("GET", fmt.Sprintf("https://api.spotify.com/v1/search?q=%v&type=artist", url.QueryEscape(artists[i].Name)), nil)
		if err != nil {
			// handle error
		}
		req.Header.Set("Authorization", "Bearer "+key)
		resp, err := client.Do(req)
		if err != nil {
			// handle error
		}
		defer resp.Body.Close()
		var result map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&result)
		if err != nil {
			// handle error
		}
		artistsObj, ok := result["artists"].(map[string]interface{})
		if !ok {
			// handle error
		}
		items, ok := artistsObj["items"].([]interface{})
		if !ok || len(items) == 0 {
			// handle error
		}
		var first map[string]interface{}
		for _, item := range items {
			artistMap, ok := item.(map[string]interface{})
			if !ok {
				continue
			}
			name, ok := artistMap["name"].(string)
			if ok && strings.EqualFold(name, artists[i].Name) {
				first = artistMap
				break
			}
		}
		// If not found, fallback to first item
		if first == nil && len(items) > 0 {
			first, _ = items[0].(map[string]interface{})
		}
		if first == nil {
			// handle error
		}
		genres, _ := first["genres"].([]interface{})
		images, _ := first["images"].([]interface{})
		id, _ := first["id"].(string)
		popularity, _ := first["popularity"].(float64) // Spotify returns float64 for numbers
		// Convert genres []interface{} to []string
		var genreList []string
		for _, g := range genres {
			if gs, ok := g.(string); ok {
				genreList = append(genreList, gs)
			}
		}
		artists[i].Genres = genreList
		// Convert images []interface{} to []string
		var imageList []string
		for _, img := range images {
			if imgMap, ok := img.(map[string]interface{}); ok {
				if url, ok := imgMap["url"].(string); ok {
					imageList = append(imageList, url)
				}
			}
		}
		artists[i].SpotifyImages = imageList

		artists[i].SpotifyID = id
		artists[i].Popularity = int(popularity)

		artists[i].Locations = api.GetLocations(artists[i])
		// artists[i].ConcertDates = api.GetConcertDates(artists[i]) */
		artists[i].Relations = api.GetRelations(artists[i])
	}
	return artists
}

// Fonction en charge d'appeler "Geocoding" pour obtenir les coordonnées des concerts du groupe de musique sondé.
func GenerateCoordinates(artist api.Artist) []api.Coordinates {
	var tabcoords []api.Coordinates
	for _, place := range artist.Locations {
		coords, err := api.Geocoding(place)
		if err == nil {
			tabcoords = append(tabcoords, coords)
		} else {
			fmt.Println("Couldn't complete geocoding")
		}
	}
	return tabcoords
}
