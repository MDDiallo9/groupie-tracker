package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Coordinates struct {
	Name string
	Lat  string `json:"lat"`
	Lon  string `json:"lon"`
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
	results.Name = location

	return results, nil
}

/*
func haversine(lat1, lon1, lat2, lon2 float64) float64 {
	const R = 6371 // Rayon terrestre en km
	dLat := (lat2 - lat1) * math.Pi / 180
	dLon := (lon2 - lon1) * math.Pi / 180
	lat1R := lat1 * math.Pi / 180
	lat2R := lat2 * math.Pi / 180

	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Sin(dLon/2)*math.Sin(dLon/2)*math.Cos(lat1R)*math.Cos(lat2R)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	return R * c
}

func determineZoom(maxDistance float64) int {
	switch {
	case maxDistance > 10000:
		return 1
	case maxDistance > 5000:
		return 2
	case maxDistance > 2000:
		return 3
	case maxDistance > 1000:
		return 4
	case maxDistance > 500:
		return 5
	default:
		return 6
	}
}

func GenerateMapUrl(locations []Coordinates) (string, error) {
	var sumLat, sumLon float64
	var validCoords []struct {
		lat  float64
		lon  float64
		name string
	}

	var markerParams string

	for _, coords := range locations {
		lat, errLat := strconv.ParseFloat(coords.Lat, 64)
		lon, errLon := strconv.ParseFloat(coords.Lon, 64)
		if errLat == nil && errLon == nil {
			sumLat += lat
			sumLon += lon
			validCoords = append(validCoords, struct {
				lat  float64
				lon  float64
				name string
			}{lat, lon, coords.Name})
		}
	}

	if len(validCoords) == 0 {
		return "", fmt.Errorf("no valid coordinates")
	}

	// Calcul centre
	centerLat := sumLat / float64(len(validCoords))
	centerLon := sumLon / float64(len(validCoords))

	// Calcul distance maximale entre paires
	var maxDistance float64
	for i := 0; i < len(validCoords); i++ {
		for j := i + 1; j < len(validCoords); j++ {
			dist := haversine(validCoords[i].lat, validCoords[i].lon, validCoords[j].lat, validCoords[j].lon)
			if dist > maxDistance {
				maxDistance = dist
			}
		}
	}

	// Zoom adapté
	zoom := determineZoom(maxDistance)

	// Création des marqueurs enrichis
	for _, loc := range validCoords {
		label := url.QueryEscape(loc.name)
		marker := fmt.Sprintf("&marker=lonlat:%f,%f;type:circle;size:48;text:%s;color:%%23ff0000", loc.lon, loc.lat, label)
		markerParams += marker
	}

	mapImageUrl := fmt.Sprintf(
		"https://maps.geoapify.com/v1/staticmap?style=osm-bright&width=800&height=600&center=lonlat:%f,%f&zoom=%d%s&apiKey=%s",
		centerLon, centerLat,
		zoom,
		markerParams,
		"90bfb2c4beb645718fdc2c925fe235a1",
	)

	return mapImageUrl, nil
}

//90bfb2c4beb645718fdc2c925fe235a1
*/
