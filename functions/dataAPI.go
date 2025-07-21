package functions

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

// L'objectif est d'appeler l'API centralisant les autres URL.
// Puis de récupérer les URL dedans pour les utiliser par les suite.

// #MARK: Déclaration de la Struct
type Links struct {
	Artists   string
	Locations string
	Dates     string
	Relation  string
}

// #MARK: Lecture des URL.
func DataAPI() {
	dataAPIdata, err := http.Get("https://groupietrackers.herokuapp.com/api")
	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}

	// Convertit le dataAPI en Bytes.
	dataAPIBytes, err := io.ReadAll(dataAPIdata.Body)
	if err != nil {
		log.Fatal(err)
	}

	if len(dataAPIBytes) == 0 {
		fmt.Println("API Empty")
		return
	}

	// Déclaration de cette variable pour amener la Struct à obtenir les données par ce lien.
	var links Links

	err = json.Unmarshal(dataAPIBytes, &links)
	if err != nil {
		log.Fatal(err)
	}

	/*
		 	// test pour savoir si ça s'affiche bien.
			fmt.Println(links.Artists)
			fmt.Println(links.Locations)
			fmt.Println(links.Dates)
			fmt.Println(links.Relation)
	*/
}
