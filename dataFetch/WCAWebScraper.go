package dataFetch

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/castus/speedcube-events/db"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const (
	registrationsPath = "registrations"
	WCAHost           = "worldcubeassociation.org"
)

func IncludeRegistrations(competitions db.Competitions, fetcher DataFetcher) db.Competitions {
	var newCompetitions db.Competitions
	for _, competition := range competitions {
		if competition.HasWCAPage() {
			competition.Registered = registrations(competition.URL, fetcher)
			time.Sleep(time.Millisecond * 500)
		}
		newCompetitions = append(newCompetitions, competition)
	}

	return newCompetitions
}

func registrations(URL string, fetcher DataFetcher) int {
	webURL := fmt.Sprintf("%s/%s", URL, registrationsPath)
	log.Info("Trying to fetch registrations page", "webURL", webURL)
	r, ok := fetcher.Fetch(webURL)
	if !ok {
		return 0
	}
	doc, err := goquery.NewDocumentFromReader(r)

	if err != nil {
		return 0
	}

	numberOfRegistrations := doc.Find("#competition-data table tbody tr").Size()
	log.Info("Found registrations.", "numberOfRegistrations", numberOfRegistrations, "webURL", webURL)

	return numberOfRegistrations
}

type GeneralInfo struct {
	MainEvent       string
	CompetitorLimit int
}

func IncludeGeneralInfo(competitions db.Competitions, fetcher DataFetcher) db.Competitions {
	var newCompetitions db.Competitions
	for _, competition := range competitions {
		if competition.HasWCAPage() {
			gi := generalInfo(competition.URL, fetcher)
			competition.MainEvent = gi.MainEvent
			competition.CompetitorLimit = gi.CompetitorLimit
			time.Sleep(time.Millisecond * 500)
		}
		newCompetitions = append(newCompetitions, competition)
	}

	return newCompetitions
}

func generalInfo(URL string, fetcher DataFetcher) GeneralInfo {
	gi := GeneralInfo{}

	webURL := fmt.Sprintf("%s", URL)
	log.Info("Trying to fetch general info page", "webURL", webURL)
	r, ok := fetcher.Fetch(webURL)
	if !ok {
		return gi
	}
	doc, err := goquery.NewDocumentFromReader(r)

	if err != nil {
		return gi
	}

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
