package functions

import (
	"fmt"
	"net/http"
)

func ServerGroupieTracker() {

	// http.Dir transforme un chemin vers un répertoire local en un type compatible avec un serveur HTTP.
	// http.FileServer permet de transformer un fichier local en un mini serveur web.
	fileServer := http.FileServer(http.Dir("static"))

	// Http.Handle fait le lien entre une URL et la fonction qui répond à cette URL.
	http.Handle("/static/", http.StripPrefix("/static", fileServer))

	// http.ListAndServe permet de lancer un serveur, d'écouter les requêtes entrantes sur une adresse donnée, puis de déléguer la gestion des requêtes à un "Handler"
	fmt.Println("Serveur démarré sur le port 8081...")
	http.ListenAndServe(":8081", nil)
}
