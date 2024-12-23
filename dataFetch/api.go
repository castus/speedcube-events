package dataFetch

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

const (
	apiHost               = "https://api.worldcubeassociation.org"
	registrationApiHostV2 = "https://worldcubeassociation.org/api/v1/registrations"
	searchApi             = "https://www.worldcubeassociation.org/api/v0/competitions"
)

type WCAApiDTO struct {
	DatabaseId string
	OtherId    string
}

type SearchWCAApiDTO struct {
	Name string
}

type basicSearchJSONResponse struct {
	Id string `json:"id"`
}

type basicInfoJSONResponse struct {
	Events              []string `json:"event_ids"`
	MainEvent           string   `json:"main_event_id"`
	CompetitorLimit     int      `json:"competitor_limit"`
	RegistrationVersion string   `json:"registration_version"`
	Latitude            float32  `json:"latitude_degrees"`
	Longitude           float32  `json:"longitude_degrees"`
}
type registrationsJSONResponse struct {
	Id int `json:"id"`
}

type WCAApiResponse struct {
	Events          []string
	MainEvent       string
	CompetitorLimit int
	Registered      int
	Longitude       float32
	Latitude        float32
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

	res, err := http.Get("https://raw.githubusercontent.com/robiningelbrecht/wca-rest-api/master/api/events.json")
	if err != nil {
		log.Error("Couldn't fetch API page", "error", err)

		return EventsMap{}
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Error("Can't fetch Event from API", "Status code", res.StatusCode, "status", res.Status)

		return EventsMap{}
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Error("Couldn't get response body", "error", err)

		return EventsMap{}
	}

	var data EventsMapResponse
	err = json.Unmarshal(body, &data)
	if err != nil {
		log.Error("Couldn't unmarshal JSON", "error", err)
		return EventsMap{}
	}

	log.Info("Found events map")

	return data.Items
}

func SearchWCAApiData(item SearchWCAApiDTO) (string, error) {
	return fetchSearchInfo(item.Name)
}

func GetWCAApiData(ids []WCAApiDTO) map[string]WCAApiResponse {
	var events = make(map[string]WCAApiResponse)
	for _, identifier := range ids {
		basicInfo := fetchBasicInfo(identifier.OtherId)
		time.Sleep(time.Millisecond * 500)

		registrationUrl := ""
		if basicInfo.RegistrationVersion == "v1" {
			registrationUrl = fmt.Sprintf("%s/competitions/%s/registrations", apiHost, identifier.OtherId)
		} else {
			registrationUrl = fmt.Sprintf("%s/%s", registrationApiHostV2, identifier.OtherId)
		}

		registered := fetchCompetitors(identifier.OtherId, registrationUrl)

		events[identifier.DatabaseId] = WCAApiResponse{
			Events:          basicInfo.Events,
			MainEvent:       basicInfo.MainEvent,
			CompetitorLimit: basicInfo.CompetitorLimit,
			Registered:      registered,
			Latitude:        float32(basicInfo.Latitude),
			Longitude:       float32(basicInfo.Longitude),
		}
	}

	return events
}

func fetchBasicInfo(id string) basicInfoJSONResponse {
	log.Info("Trying to fetch basic info from WCA Api.", "WCAId", id)

	jsonData, err := fetchApi(fmt.Sprintf("%s/competitions/%s", apiHost, id))
	if err != nil {
		log.Error("Error requesting WCA Api", "error", err)
		return basicInfoJSONResponse{}
	}

	var data basicInfoJSONResponse
	err = json.Unmarshal(jsonData, &data)
	if err != nil {
		log.Error("Couldn't unmarshal JSON", "error", err)
		return basicInfoJSONResponse{}
	}

	log.Info("Found basic info from WCA Api.", "WCAId", id, "Data", data)

	return data
}

func fetchSearchInfo(name string) (string, error) {
	log.Info("Trying to search for WCA Id.", "Event name", name)

	jsonData, err := fetchApi(fmt.Sprintf("%s?q=%s", searchApi, url.QueryEscape(name)))
	if err != nil {
		log.Error("Error requesting WCA Api", "error", err)
		return "", err
	}

	var data []basicSearchJSONResponse
	err = json.Unmarshal(jsonData, &data)
	if err != nil {
		log.Error("Couldn't unmarshal JSON", "error", err)
		return "", err
	}
	if len(data) == 0 {
		return "", fmt.Errorf("search response was empty")
	}

	out := data[0].Id

	log.Info("Found ID for a WCA using Search API", "WCAId", out)

	return out, nil
}

func fetchCompetitors(id string, url string) int {
	log.Info("Trying to fetch registered.", "WCAId", id)
	var registered = 0

	jsonData, err := fetchApi(url)
	if err != nil {
		log.Error("Error requesting WCA Api", "error", err)
		return registered
	}

	var data []registrationsJSONResponse
	err = json.Unmarshal(jsonData, &data)
	if err != nil {
		log.Error("Couldn't unmarshal JSON", "error", err)
		return registered
	}

	registered = len(data)

	log.Info("Found number of registered.", "WCAId", id, "Registered", registered)

	return registered
}

func fetchApi(url string) ([]byte, error) {
	res, err := http.Get(url)
	if err != nil {
		log.Error("Can't fetch WCA API", "error", err)
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		log.Error("Can't fetch WCA API", "Status code", res.StatusCode, "status", res.Status)
		return nil, err
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Error("Couldn't get response body from WCA API", "error", err)
		return nil, err
	}

	return body, nil
}
