package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
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
