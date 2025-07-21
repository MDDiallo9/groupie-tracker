package functions

type ArtistsData struct {
	ID           int      // "1"
	Image        string   // "URL"
	Name         string   // "NAME"
	Members      []string // "NAME FIRST-NAME"
	CreationDate string   // "AAAA"
	FirstAlbum   string   // "DD-MM-YYYY"
	Locations    string   // "URL"
	ConcertDate  string   // "URL"
	Relations    string   // "URL"
}

type LocationsData struct {
	ID        int    // "1"
	Locations string // "CityName or CityName_SecondPart" - "CountryName ou CountryName_SecondPart"`
	Dates     string // "URL"
}

type DatesData struct {
	ID    int
	Dates string // "DD-MM-YYYY" Si il y a une "*", c'est une nouvelle ville où se trouve le concert.

}

type RelationsData struct {
	ID             int    // "1"
	DatesLocations string // "CityName or CityName_SecondPart" - "CountryName ou CountryName_SecondPart" ":" "DD-MM-YYYY"
}

func DataRecolt(Links) {
	// Utiliser json.Unmarshal() pour organiser les données

}
