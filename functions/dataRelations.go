package functions

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
)

// #MARK: Struct des Relations
type RelationsData struct {
	ID        int                 `json:"id"`
	Relations map[string][]string `json:"datesLocations"`
}

// Pour gérer l'index du json.
type middlemenRelations struct {
	Index []RelationsData `json: "index"`
}

// #MARK: Données des Relations
func ObtainRelations(links Links) []RelationsData {

	adress, err := http.Get(links.Relations)
	if err != nil {
		log.Fatal(err)
	}
	defer adress.Body.Close()

	adressBytes, err := io.ReadAll(adress.Body)
	if err != nil {
		log.Fatal(err)
	}

	var relationsData middlemenRelations

	err = json.Unmarshal(adressBytes, &relationsData)
	if err != nil {
		log.Fatal(err)
	}

	return relationsData.Index
}
