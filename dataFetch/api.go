package dataFetch

import (
	"encoding/json"
	"fmt"
	"github.com/castus/speedcube-events/db"
	"io"
	"net/http"
	"time"
)

const (
	apiHost = "https://raw.githubusercontent.com/robiningelbrecht/wca-rest-api/master/api/competitions"
)

type jsonResponse struct {
	Events []string `json:"events"`
}

func IncludeEvents(competitions db.Competitions) db.Competitions {
	var newCompetitions db.Competitions
	for _, competition := range competitions {
		if competition.HasWCAPage() {
			competition.Events = events(competition.WCAId)
			time.Sleep(time.Millisecond * 500)
		}
		newCompetitions = append(newCompetitions, competition)
	}

	return newCompetitions
}

func events(id string) []string {
	log.Info("Trying to fetch events.", "WCAId", id)

	res, err := http.Get(fmt.Sprintf("%s/%s.json", apiHost, id))
	if err != nil {
		log.Error("Couldn't fetch API page", err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Error("Can't fetch Event from API", "Status code", res.StatusCode, "status", res.Status)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Error("Couldn't get response body", err)
	}

	var data jsonResponse
	err = json.Unmarshal(body, &data)
	if err != nil {
		log.Error("Couldn't unmarshal JSON", err)
	}

	log.Info("Found events for.", "WCAId", id, "Events", data.Events)

	return data.Events
}
