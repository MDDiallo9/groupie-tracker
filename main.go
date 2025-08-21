package main

import (
	"flag"
	"groupie-tracker/server"
	"log"
	"net/http"
	"os"
)

func main() {
	// flag.String permet de créer des arguments/options à utiliser lors du lancement du programme.
	// flag.String(Nom du Flag, valeur par défaut si rien n'est spécifié, description de l'attendu ou de l'usage).
	// Ici, cela nous permet de nous assurer d'utiliser le Port 8000 par défaut, mais aussi d'en prendre un autre si s'est spécifié dans l'argument de commande.
	port := flag.String("port", ":8000", "PORT")

	// log.New(où écrire le log, texte écrit en début de chaque ligne, option à afficher tel que l'heure).
	// os.Stdout est la contraction de "Standard Out", sortie standard. Ici, le terminal.
	// Chaque ligne va afficher dans le terminal, un début de texte, avec la date et l'heure dexécution. "|" sert à combiner les informations et éviter l'erreur d'avoir trop de variables qu'attendu.
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime)

	// flg.Parse permet d'utiliser les options définies précedement, et de pouvoir les gérer dans l'invite de commande du terminal.
	flag.Parse()

	// Création du serveur GO.
	srv := &http.Server{
		Addr:     *port,           // définition du port à utiliser.
		ErrorLog: errorLog,        // Loggage des erreurs.
		Handler:  server.Routes(), // Reçois les requêtes pour les diriger vers le bon fichier.
	}

	// Récupère les informations de infoLog, puis inscrit le texte suivit du numéro du port utilisé

	infoLog.Println("Starting server on http://localhost" + *port)
	err := srv.ListenAndServe() // Démarre le serveur et la capacité de recevoir des requêtes.
	errorLog.Fatal(err)

}
