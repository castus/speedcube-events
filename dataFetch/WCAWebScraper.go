package dataFetch

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/castus/speedcube-events/db"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const (
	registrationsPath = "registrations"
	WCAHost           = "worldcubeassociation.org"
)

func IncludeRegistrations(competitions db.Competitions) db.Competitions {
	var newCompetitions db.Competitions
	for _, competition := range competitions {
		if competition.HasWCAPage() {
			competition.Registered = registrations(competition.URL)
			time.Sleep(time.Millisecond * 500)
		}
		newCompetitions = append(newCompetitions, competition)
	}

	return newCompetitions
}

func registrations(URL string) int {
	webURL := fmt.Sprintf("%s/%s", URL, registrationsPath)
	log.Info("Trying to fetch registrations page", "webURL", webURL)
	doc := fetchWebPageBody(webURL)

	numberOfRegistrations := doc.Find("#competition-data table tbody tr").Size()
	log.Info("Found registrations.", "numberOfRegistrations", numberOfRegistrations, "webURL", webURL)

	return numberOfRegistrations
}

type GeneralInfo struct {
	MainEvent       string
	CompetitorLimit int
}

func IncludeGeneralInfo(competitions db.Competitions) db.Competitions {
	var newCompetitions db.Competitions
	for _, competition := range competitions {
		if competition.HasWCAPage() {
			gi := generalInfo(competition.URL)
			competition.MainEvent = gi.MainEvent
			competition.CompetitorLimit = gi.CompetitorLimit
			time.Sleep(time.Millisecond * 500)
		}
		newCompetitions = append(newCompetitions, competition)
	}

	return newCompetitions
}

func generalInfo(URL string) GeneralInfo {
	gi := GeneralInfo{}

	webURL := fmt.Sprintf("%s", URL)
	log.Info("Trying to fetch general info page", "webURL", webURL)
	doc := fetchWebPageBody(webURL)

	log.Info("Trying to find main event and competitor limit", "webURL", webURL)
	doc.Find("#general-info .dl-horizontal dt").Each(func(i int, s *goquery.Selection) {
		if s.Text() == "Main event" {
			cls, _ := s.Next().Find("i").Attr("class")
			clsItems := strings.Split(cls, " ")
			for _, item := range clsItems {
				if strings.Contains(item, "event-") {
					eventSplit := strings.Split(item, "-")
					gi.MainEvent = eventSplit[len(eventSplit)-1]
					log.Info("Found main event", "mainEvent", gi.MainEvent)
					break
				}
			}
		}
		if s.Text() == "Competitor limit" {
			text := s.Next().Text()
			reg := regexp.MustCompile("[^0-9]+")
			text = reg.ReplaceAllString(text, "")
			i, err := strconv.Atoi(text)
			if err != nil {
				log.Error("Error during converting competitor limit from string to int", "webURL", webURL, "string", text, "error", err)
			} else {
				gi.CompetitorLimit = i
				log.Info("Found Competitor limit", "competitorLimit", gi.CompetitorLimit)
			}
		}
	})

	return gi
}

func fetchWebPageBody(URL string) *goquery.Document {
	res, err := http.Get(URL)
	if err != nil {
		log.Error("Couldn't fetch page to scrap", "error", err, "url", URL)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Error("Status code error", "status code", res.StatusCode, "status", res.Status)
	}
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Error("Couldn't load HTML", "error", err, "url", URL)
	}

	return doc
}
