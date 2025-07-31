package server

import (
	"fmt"
	"groupie-tracker/api"
	"html/template"
	"log"
	"net/http"
	"strconv"
)

func home(w http.ResponseWriter, r *http.Request) {
	var filter api.Filter
	if r.Method == "POST" {
		r.ParseForm()
		fad, _ := strconv.Atoi(r.Form["FirstAlbumDate"][0])
		cdMin, _ := strconv.Atoi(r.Form["creationDate"][0])
		cdMax := 2025
		members := map[int]bool{
			1: false,
			2: false,
			3: false,
			4: false,
			5: false,
			6: false,
			7: false,
		}
		for _, num := range r.Form["members"] {
			n, _ := strconv.Atoi(num)
			members[n] = true
		}
		filter = api.Filter{
			Location:       r.Form["Location"][0],
			FirstAlbumDate: fad,
			Members:        members,
			CreationDate:   []int{cdMin, cdMax},
		} // Besoin de recharger home avec le api.FilterBy(artists,filter)
	}
	
	artists := api.GetArtists()

	if isFilterFilled(filter) {
		artists = api.FilterBy(artists, filter)
		log.Println(artists,filter)
	}

	ts, err := template.ParseFiles("./templates/home.html", "./templates/partials/base.html", "./templates/partials/footer.html", "./templates/partials/head.html")
	if err != nil {
		log.Print(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	err = ts.ExecuteTemplate(w, "base.html", artists)
	if err != nil {
		log.Print(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func Artist(w http.ResponseWriter, r *http.Request) {
	idstring := r.URL.Query().Get("id")
	id, _ := strconv.Atoi(idstring)

	artist := artists[id-1]
	var coords []api.Coordinates
	for _, location := range api.GetLocations(artist) {
		coord, err := api.Geocoding(location)
		if err == nil {
			coords = append(coords, coord)
		}
	}
	mapURL, err := api.GenerateMapUrl(coords)
	if err != nil {
		fmt.Println("error when generating map url :", err)
	}
	artist.MapURL = mapURL

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

func isFilterFilled(f api.Filter) bool {
    for _, v := range f.Members {
        if v {
            return true
        }
    }
    if f.Location != "" || f.FirstAlbumDate != 0 {
        return true
    }
    if len(f.CreationDate) == 2 && (f.CreationDate[0] != 0 || f.CreationDate[1] != 2025) {
        return true
    }
    return false
}
