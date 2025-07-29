package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
)

type Coordinates struct {
	Lat string `json:"lat"`
	Lon string `json:"lon"`
}

func Geocoding(location string) (Coordinates, error) {
	baseUrl := "https://nominatim.openstreetmap.org/search"
	client := &http.Client{}
	var results Coordinates

	// Building the URL
	req, err := http.NewRequest("GET", baseUrl, nil)
	if err != nil {
		return results, err
	}

	q := req.URL.Query()
	q.Add("q", location)
	q.Add("format", "json")
	req.URL.RawQuery = q.Encode()

	// Necessary by API policies
	req.Header.Set("User-Agent", "GroupieTracker/1.0 (nathpacc@gmail.com)")

	// Sending the request
	resp, err := client.Do(req)
	if err != nil {
		return results, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return results, fmt.Errorf("HTTP error %d", resp.StatusCode)
	}

	// Recovering datas from JSON

	bodyBytes, _ := io.ReadAll(resp.Body)
	var rawResults []Coordinates
	json.Unmarshal(bodyBytes, &rawResults)

	if len(rawResults) > 0 {
		results = rawResults[0]
	}

	return results, nil
}

func GenerateMapUrl(locations []Coordinates) (string, error) {
	var sumLat, sumLon float64
	var count int
	var markers []string

	if len(locations) == 0 {
		return "", fmt.Errorf("no coordinates given")
	}

	for _, coords := range locations {
		lat, errorLat := strconv.ParseFloat(coords.Lat, 64)
		lon, errorLon := strconv.ParseFloat(coords.Lon, 64)
		if errorLat != nil && errorLon != nil {
			sumLat += lat
			sumLon += lon
			count++
			markers = append(markers, fmt.Sprintf("%s,%s,red-pushpin", coords.Lat, coords.Lon))
		} else {
			fmt.Printf("map_generator: Warning: Impossible to convert convertir Lat/Lon (%s, %s) for center estimation: %v, %v\n", coords.Lat, coords.Lon, errorLat, errorLon)
		}
	}

	var centerLat, centerLon string
	if count > 0 {
		centerLat = fmt.Sprintf("%f", sumLat/float64(count))
		centerLon = fmt.Sprintf("%f", sumLon/float64(count))
	} else {
		return "", fmt.Errorf("no valid coordinates found for map center")
	}

	if len(markers) == 0 {
		return "", fmt.Errorf("no valid marker")
	}

	mapImageUrl := fmt.Sprintf(
		"https://staticmap.openstreetmap.de/staticmap.php?center=%s,%s&zoom=3&size=800x600&maptype=mapnik&markers=%s",
		centerLat, centerLon, strings.Join(markers, "|"),
	)

	return mapImageUrl, nil

}
