package server

import (
	"encoding/json"
	"html/template"
	"log"
	"math/rand"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"
	"io"
	"groupie-tracker/api"
)

type Suggestion struct {
	Texte string // ce qui est envoyé en résultat de recherche.
	Label string // Ce qui est écrit à côté pour préciser l'origine de la donnée.
}

func home(w http.ResponseWriter, r *http.Request, init *AppData) {
	var filter api.Filter
	filter.CreationDate = []int{1950, 2025}   // Définition du Min-Max.
	filter.FirstAlbumDate = []int{1950, 2020} // Définition du Min-Max.

	// Ordre popularité
	popular := make([]api.Artist, 15)
	copy(popular, init.Artists[:15])
	sort.Slice(popular, func(i, j int) bool {
		return popular[i].Popularity > popular[j].Popularity
	})

	if r.Method == "POST" {

		r.ParseForm()

		// Min-Max choisies des dates de sorties du Premier Album par l'utilisateur.
		minFAD, _ := strconv.Atoi(r.Form.Get("minfAd"))
		maxFAD, _ := strconv.Atoi(r.Form.Get("maxfAd"))

		// Min-Max choisies des dates de foncdations du Groupe de Musique par l'utilisateur.
		minCD, _ := strconv.Atoi(r.Form.Get("minCD"))
		maxCD, _ := strconv.Atoi(r.Form.Get("maxCD"))

		// Map pour le filtre de cases à cocher pour le nombre de membres du groupe.
		members := map[int]bool{
			1: false, 2: false, 3: false, 4: false, 5: false, 6: false, 7: false, // Si true, alors sélectionné par l'utilisateur.
		}
		for _, num := range r.Form["members"] { // "r.Form["members"]" contient les valeurs envoyées par le formulaire.
			n, _ := strconv.Atoi(num) // conversion de la réponse en INT.
			members[n] = true
		}

		// Structure du filtre.
		filter = api.Filter{
			Location:       r.Form.Get("Location"),
			FirstAlbumDate: []int{minFAD, maxFAD},
			Members:        members,
			CreationDate:   []int{minCD, maxCD},
		}
	}

	// Structure pour gérer les suggestions et leurs apparitions sur le page HTML.
	data := struct {
		Suggestions []Suggestion
		Artists     []api.Artist
		Playlist20  []api.Artist
		Location    []api.Artist
		Popular     []api.Artist
	}{
		Suggestions: SuggestionsGeneration(init),
		Artists:     init.Artists,
		Playlist20:  api.FilterBy(init.Artists, api.Filter{CreationDate: []int{2000, 2010}}),
		Location:    api.FilterBy(init.Artists, api.Filter{Location: "france", CreationDate: []int{1950, 2025}}),
		Popular:     popular,
	}

	// Vérifie si il y a une demande du client pour obtenir des données au travers du filtre.
	if isFilterFilled(filter) {
		data.Artists = api.FilterBy(init.Artists, filter)
		log.Print(filter)
	}

	// Appel des pages HTML pour afficher les informations.
	ts, err := template.ParseFiles("./templates/home.html", "./templates/partials/base.html", "./templates/partials/footer.html", "./templates/partials/head.html")
	if err != nil {
		log.Print(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Chargement de la base.html avec toutes les options précédentes.
	err = ts.ExecuteTemplate(w, "base.html", data)
	if err != nil {
		log.Print(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

// Fonction qui gère l'affichage des informations pour la page de chaque artiste.
func Artist(w http.ResponseWriter, r *http.Request, init *AppData) {
	idstring := r.URL.Query().Get("id") // Récupération de l'ID du groupe de musique sur lequel l'utilisateur veut aller.
	id, err := strconv.Atoi(idstring)

	/* if err != nil || id < 1 || id > len(artists) {
		http.Error(w, "Invalid artist ID", http.StatusBadRequest)
		return
	} */

	artist := init.Artists[id-1]

	// Appel des pages HTML pour afficher les informations.
	ts, err := template.ParseFiles("./templates/artist.html", "./templates/partials/base.html", "./templates/partials/footer.html", "./templates/partials/head.html")
	if err != nil {
		log.Print(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Chargement de la base.html avec toutes les options précédentes.
	err = ts.ExecuteTemplate(w, "base.html", artist)
	if err != nil {
		log.Print(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func IndexPage(w http.ResponseWriter, r *http.Request, init *AppData) {
	data := struct {
		Artists []api.Artist
	}{
		Artists: api.FilterBy(artists, api.Filter{CreationDate: []int{2000, 2010}})[8:],
	}

	ts, err := template.ParseFiles("./templates/index.html", "./templates/partials/base.html", "./templates/partials/footer.html", "./templates/partials/head.html")
	if err != nil {
		log.Print(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	err = ts.ExecuteTemplate(w, "base.html", data)
	if err != nil {
		log.Print(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func ArtistMap(w http.ResponseWriter, r *http.Request, init *AppData) {
	idstring := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idstring)
	if err != nil || id < 1 {
		http.Error(w, "ID invalide", http.StatusBadRequest)
		return
	}

	artist := init.Artists[id-1]
	artist.TabCoords = GenerateCoordinates(artist)

	coordsJSON, err := json.Marshal(artist.TabCoords)
	if err != nil {
		log.Print("Erreur JSON:", err)
		http.Error(w, "Erreur interne", http.StatusInternalServerError)
		return
	}

	artist.CoordsJSON = template.JS(coordsJSON)

	tmpl := template.Must(template.ParseFiles(
		"./templates/map.html",
		"./templates/partials/base.html",
		"./templates/partials/head.html",
		"./templates/partials/footer.html",
	))

	err = tmpl.ExecuteTemplate(w, "base.html", artist)
	if err != nil {
		log.Print("Erreur template:", err)
		http.Error(w, "Erreur interne", http.StatusInternalServerError)
	}
}

func search(w http.ResponseWriter, r *http.Request, init *AppData) {
	var query string

	// Vérification que c'est une méthode POST
	if r.Method == "POST" {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Error reading request body", http.StatusBadRequest)
			return
		}
		query = strings.ToLower(string(body))
		log.Print(query)
	} else {
		// Si ce n'est pas la méthode POST.
		http.Error(w, "Erreur d'envoi de données", http.StatusMethodNotAllowed)
		return
	}

	// Vérification de la taille de la recherche.
	if len(query) == 0 {
		http.Error(w, "Veuillez écrire au moins un caractère dans la barre de recherche.", http.StatusBadRequest)
		return
	}

	unique := make(map[string]bool) // Pour éviter les doublons.
	var results []api.Artist        // Format pour accueillir les multiples informations

	// "answer" parcours chaque sous-ensemble de la structure Artists.
	for _, answer := range init.Artists { // Plutôt que de créer une variable, on exploite directement la struct depuis sa fonction.
		nonUnique := false

		// Boucle de recherche pour la date de création, du premier album et le nom du groupe.
		if query == strconv.Itoa(answer.CreationDate) ||
			strings.Contains(strings.ToLower(answer.Name), query) ||
			strings.Contains(strings.ToLower(answer.FirstAlbum), query) {
			nonUnique = true
			// answer.Relations = api.GetRelations(answer)
			// results = append(results, answer)
		}

		// Boucle de recherche pour les noms des artistes.
		for _, response := range answer.Members {
			if strings.Contains(strings.ToLower(response), query) {
				// answer.Relations = api.GetRelations(answer)
				nonUnique = true
				// results = append(results, answer)
				break
			}
		}

		// Boucle de recherche pour les relations entre dates et lieux de concerts.
		for date, cities := range answer.Relations {
			for _, city := range cities {
				if strings.Contains(strings.ToLower(date), query) || strings.Contains(strings.ToLower(city), query) {
					// answer.Relations = api.GetRelations(answer)
					nonUnique = true
					// results = append(results, answer)
					break
				}
			}
		}

		// Condition pour éviter les doublons.
		if nonUnique {
			identifiant := strings.ToLower(answer.Name)
			if !unique[identifiant] {
				unique[identifiant] = true
				results = append(results, answer)
			}
		}
	}

	// si la recherche ne correspond à rien.
	if len(results) == 0 {
		http.Redirect(w, r, "/404", http.StatusSeeOther)
		return
	}

	data := struct {
		Suggestions []Suggestion
		Artists     []api.Artist
	}{
		Suggestions: SuggestionsGeneration(init),
		Artists:     results,
	}

	// Exécution du template
	/* ts, err := template.ParseFiles("./templates/index.html", "./templates/partials/base.html", "./templates/partials/footer.html", "./templates/partials/head.html")
	if err != nil {
		log.Print("Erreur template:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	err = ts.ExecuteTemplate(w, "base.html", data)
	if err != nil {
		log.Print("Erreur exécution template:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	} */
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data.Artists)
}

func NotFound(w http.ResponseWriter, r *http.Request) {
	ts, err := template.ParseFiles("./templates/404.html", "./templates/partials/base.html", "./templates/partials/footer.html", "./templates/partials/head.html")
	if err != nil {
		log.Print("Erreur template:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	err = ts.ExecuteTemplate(w, "base.html", nil)
	if err != nil {
		log.Print("Erreur exécution template:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

// Construction des suggestions
func SuggestionsGeneration(init *AppData) []Suggestion {
	var suggestions []Suggestion

	for _, artist := range init.Artists {

		// Suggestions pour : Groupe de Musique, Date de création du Groupe,Date de sortie du Premier Album.
		suggestions = append(suggestions,
			Suggestion{
				Texte: artist.Name,
				Label: "Groupe de musique",
			},
			Suggestion{
				Texte: strconv.Itoa(artist.CreationDate),
				Label: "Date de Fondation du Groupe : " + artist.Name,
			},
			Suggestion{
				Texte: artist.FirstAlbum,
				Label: "Date de sortie du premier Album du groupe " + artist.Name,
			},
		)

		// Suggestion pour : Artistes composant un groupe de musique.
		for _, member := range artist.Members {
			suggestions = append(suggestions, Suggestion{
				Texte: member,
				Label: "Membre du groupe : " + artist.Name,
			})
		}

		// Suggestion pour : les villes des concerts.
		for _, city := range artist.Locations {
			suggestions = append(suggestions, Suggestion{
				Texte: city,
				Label: "Ville où joue le groupe " + artist.Name,
			})
		}
	}
	return suggestions
}

// Vérification des filtres de recherches.
func isFilterFilled(f api.Filter) bool {
	for _, v := range f.Members {
		if v {
			return true
		}
	}
	if f.Location != "" || (f.FirstAlbumDate[0] != 1950 || f.FirstAlbumDate[1] != 2020) {
		return true
	}
	if len(f.CreationDate) == 2 && (f.CreationDate[0] != 1950 || f.CreationDate[1] != 2025) {
		return true
	}
	return false
}

func GetRandomArtist(artists []api.Artist) api.Artist {
	if len(artists) == 0 {
		return api.Artist{}
	}
	rand.Seed(time.Now().UnixNano())
	return artists[rand.Intn(len(artists))]
}

func FilterRefresh(w http.ResponseWriter, r *http.Request, init *AppData) {
	decades := map[string][]int{
		"60s":   {1960, 1969},
		"70s":   {1970, 1979},
		"80s":   {1980, 1989},
		"90s":   {1990, 1999},
		"2000s": {2000, 2009},
		"2010s": {2010, 2019},
		"2020s": {2020, 2029},
	}
	locations := map[string]string{
		"fr":   "france",
		"us":   "usa",
		"uk":   "uk",
		"ger":  "germany",
		"den:": "denmark",
		"swe":  "sweden",
		"aus":  "australia",
		"indo": "indonesia",
	}

	dec := r.URL.Query().Get("dec")
	loc := r.URL.Query().Get("loc")

	var filter api.Filter

	if years, ok := decades[dec]; ok {
		filter.CreationDate = years
	}
	if location, ok := locations[loc]; ok {
		filter.Location = location
		log.Print(filter)
	}

	filtered := api.FilterBy(init.Artists, filter)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(filtered)
}
