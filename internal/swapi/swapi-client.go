package swapi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type SWAPIClient struct {
	endpoint string
	apiKey   string

	httpClient *http.Client
}

type Planet struct {
	Id         string `json:"id"`
	Name       string `json:"name"`
	Diameter   int64  `json:"diameter"`
	Population int64  `json:"population"`
}

func NewSWAPIClient(endpoint string, apiKey string) *SWAPIClient {
	return &SWAPIClient{endpoint, apiKey, http.DefaultClient}
}

func (swapiClient *SWAPIClient) ReadPlanetById(planetID string) (*Planet, error) {
	url := fmt.Sprintf("%s/planets/%s", swapiClient.endpoint, planetID)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("api-key", swapiClient.apiKey)

	resp, httpError := swapiClient.httpClient.Do(req)
	if httpError != nil {
		return nil, httpError
	}
	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("Got HTTP status %s", resp.Status)
	}

	var planet Planet

	jsonError := json.NewDecoder(resp.Body).Decode(&planet)
	if jsonError != nil {
		return nil, jsonError
	}

	return &planet, nil
}

func (swapiClient *SWAPIClient) ReadPlanetByName(planetName string) (*Planet, error) {
	url := fmt.Sprintf("%s/planets?name=%s", swapiClient.endpoint, planetName)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("api-key", swapiClient.apiKey)

	resp, httpError := swapiClient.httpClient.Do(req)
	if httpError != nil {
		return nil, httpError
	}
	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("Got HTTP status %s", resp.Status)
	}

	var planet Planet

	jsonError := json.NewDecoder(resp.Body).Decode(&planet)
	if jsonError != nil {
		return nil, jsonError
	}

	return &planet, nil
}

func (swapiClient *SWAPIClient) CreateOrUpdatePlanet(planet *Planet) (*Planet, error) {
	url := fmt.Sprintf("%s/planets", swapiClient.endpoint)

	jsonData, err := json.Marshal(planet)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("api-key", swapiClient.apiKey)

	resp, httpError := swapiClient.httpClient.Do(req)
	if httpError != nil {

		return nil, httpError
	}
	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("Got HTTP status %s", resp.Status)
	}

	jsonError := json.NewDecoder(resp.Body).Decode(&planet)
	if jsonError != nil {
		return nil, jsonError
	}

	return planet, nil
}

func (swapiClient *SWAPIClient) DeletePlanet(planetId string) error {
	url := fmt.Sprintf("%s/planets/%s", swapiClient.endpoint, planetId)

	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return err
	}

	req.Header.Add("api-key", swapiClient.apiKey)

	resp, err := swapiClient.httpClient.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode >= 400 {
		return fmt.Errorf("Got HTTP status %s", resp.Status)
	}
	return nil
}
