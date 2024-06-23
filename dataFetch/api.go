package dataFetch

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const (
	apiHost = "https://api.worldcubeassociation.org"
)

type WCAApiDTO struct {
	DatabaseId string
	OtherId    string
}

type basicInfoJSONResponse struct {
	Events          []string `json:"event_ids"`
	MainEvent       string   `json:"main_event_id"`
	CompetitorLimit int      `json:"competitor_limit"`
}
type registrationsJSONResponse struct {
	Id int `json:"id"`
}

type WCAApiResponse struct {
	Events          []string
	MainEvent       string
	CompetitorLimit int
	Registered      int
}

type EventsMapResponse struct {
	Items []EventMap `json:"items"`
}

type EventMap struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type EventsMap []EventMap

func (e EventsMap) NameById(id string) (string, bool) {
	for _, event := range e {
		if event.Id == id {
			return event.Name, true
		}
	}

	return "", false
}

func (e EventsMap) IdByName(name string) (string, bool) {
	for _, event := range e {
		if event.Name == name {
			return event.Id, true
		}
	}

	return "", false
}

func InitializeEventsMap() EventsMap {
	log.Info("Trying to fetch events map.")

	res, err := http.Get(fmt.Sprintf("https://raw.githubusercontent.com/robiningelbrecht/wca-rest-api/master/api/events.json"))
	if err != nil {
		log.Error("Couldn't fetch API page", err)

		return EventsMap{}
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Error("Can't fetch Event from API", "Status code", res.StatusCode, "status", res.Status)

		return EventsMap{}
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Error("Couldn't get response body", err)

		return EventsMap{}
	}

	var data EventsMapResponse
	err = json.Unmarshal(body, &data)
	if err != nil {
		log.Error("Couldn't unmarshal JSON", err)
		return EventsMap{}
	}

	log.Info("Found events map")

	return data.Items
}

func GetWCAApiData(ids []WCAApiDTO) map[string]WCAApiResponse {
	var events = make(map[string]WCAApiResponse)
	for _, identifier := range ids {
		basicInfo := fetchBasicInfo(identifier.OtherId)
		time.Sleep(time.Millisecond * 500)
		registered := fetchCompetitors(identifier.OtherId)
		events[identifier.DatabaseId] = WCAApiResponse{
			Events:          basicInfo.Events,
			MainEvent:       basicInfo.MainEvent,
			CompetitorLimit: basicInfo.CompetitorLimit,
			Registered:      registered,
		}
	}

	return events
}

func fetchBasicInfo(id string) basicInfoJSONResponse {
	log.Info("Trying to fetch basic info from WCA Api.", "WCAId", id)

	jsonData, err := fetchApi(fmt.Sprintf("%s/competitions/%s", apiHost, id))
	if err != nil {
		log.Error("Error requesting WCA Api", err)
		return basicInfoJSONResponse{}
	}

	var data basicInfoJSONResponse
	err = json.Unmarshal(jsonData, &data)
	if err != nil {
		log.Error("Couldn't unmarshal JSON", err)
		return basicInfoJSONResponse{}
	}

	log.Info("Found basic info from WCA Api.", "WCAId", id, "Data", data)

	return data
}

func fetchCompetitors(id string) int {
	log.Info("Trying to fetch registered.", "WCAId", id)
	var registered = 0

	jsonData, err := fetchApi(fmt.Sprintf("%s/competitions/%s/registrations", apiHost, id))
	if err != nil {
		log.Error("Error requesting WCA Api", err)
		return registered
	}

	var data []registrationsJSONResponse
	err = json.Unmarshal(jsonData, &data)
	if err != nil {
		log.Error("Couldn't unmarshal JSON", err)
		return registered
	}

	registered = len(data)

	log.Info("Found number of registered.", "WCAId", id, "Registered", registered)

	return registered
}

func fetchApi(url string) ([]byte, error) {
	res, err := http.Get(url)
	if err != nil {
		log.Error("Can't fetch WCA API", err)
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		log.Error("Can't fetch WCA API", "Status code", res.StatusCode, "status", res.Status)
		return nil, err
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Error("Couldn't get response body from WCA API", err)
		return nil, err
	}

	return body, nil
}
