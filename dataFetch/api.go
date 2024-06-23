package dataFetch

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const (
	apiHost = "https://raw.githubusercontent.com/robiningelbrecht/wca-rest-api/master/api"
)

type jsonResponse struct {
	Events []string `json:"events"`
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

	res, err := http.Get(fmt.Sprintf("%s/events.json", apiHost))
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

func GetEvents(ids []string) map[string][]string {
	var events = make(map[string][]string)
	for _, identifier := range ids {
		events[identifier] = fetchEvents(identifier)
		time.Sleep(time.Millisecond * 500)
	}

	return events
}

// func IncludeEvents(competitions db.Competitions) db.Competitions {
// 	var newCompetitions db.Competitions
// 	for _, competition := range competitions {
// 		if competition.Type == db.CompetitionType.WCA {
// 			var id string
// 			if competition.WCAId != "" {
// 				id = competition.WCAId
// 			} else {
// 				id = competition.TypeSpecificId
// 			}
// 			competition.Events = fetchEvents(id)
// 			time.Sleep(time.Millisecond * 500)
// 		}
// 		newCompetitions = append(newCompetitions, competition)
// 	}

// 	return newCompetitions
// }

func fetchEvents(id string) []string {
	log.Info("Trying to fetch events.", "WCAId", id)

	res, err := http.Get(fmt.Sprintf("%s/competitions/%s.json", apiHost, id))
	if err != nil {
		log.Error("Couldn't fetch API page", err)

		return []string{}
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Error("Can't fetch Event from API", "Status code", res.StatusCode, "status", res.Status)

		return []string{}
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Error("Couldn't get response body", err)

		return []string{}
	}

	var data jsonResponse
	err = json.Unmarshal(body, &data)
	if err != nil {
		log.Error("Couldn't unmarshal JSON", err)
		return []string{}
	}

	log.Info("Found events for.", "WCAId", id, "Events", data.Events)

	return data.Events
}
