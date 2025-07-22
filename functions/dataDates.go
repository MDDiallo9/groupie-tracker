package functions

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
)

// #MARK: Struct des dates
type ConcertDatesData struct {
	ID    int      `json:"id"`
	Dates []string `json:"dates"`
}

// Pour gérer l'index du json.
type middlemenDates struct {
	Index []ConcertDatesData `json: "index"`
}

// #MARK: Données des dates
func ObtainDates(links Links) []ConcertDatesData {

	adress, err := http.Get(links.Dates)
	if err != nil {
		log.Fatal(err)
	}
	defer adress.Body.Close()

	adressBytes, err := io.ReadAll(adress.Body)
	if err != nil {
		log.Fatal(err)
	}

	var datesData middlemenDates

	err = json.Unmarshal(adressBytes, &datesData)
	if err != nil {
		log.Fatal(err)
	}

	return datesData.Index
}
