package externalParser

import (
	"io"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/castus/speedcube-events/dataFetch"
	"github.com/castus/speedcube-events/db"
	"github.com/castus/speedcube-events/externalFetcher"
)

func Cube4FunParse(body io.Reader, competitionType string, id string, pageName string, eventsMap dataFetch.EventsMap, dbItem *db.Competition) {
	log.Info("Found Cube4Fun event, parsing ...", "type", competitionType, "id", id, "pageName", pageName)
	pageNameItems := strings.Split(pageName, ".")
	pageKey := pageNameItems[0]
	if pageKey == externalFetcher.PageTypes.Info {
		parseInfo(body, dbItem, eventsMap)
	}
}

func parseInfo(body io.Reader, dbItem *db.Competition, eventsMap dataFetch.EventsMap) {
	doc, err := goquery.NewDocumentFromReader(body)
	if err != nil {
		log.Error("Couldn't parse HTML with GoQuery", "error", err)
	}

	doc.Find("p").Each(func(i int, row *goquery.Selection) {
		text := trim(row.Text())
		text = strings.ToLower(text)
		if text == "competitor limit" || text == "limit zawodników" {
			limit := trim(row.Next().Text())
			limits := strings.Split(limit, "/")
			var registered string
			var all string
			if len(limits) < 2 {
				registered = "0"
				all = limits[0]
			} else {
				registered = limits[0]
				all = limits[1]
			}

			i, err := strconv.Atoi(all)
			if err != nil {
				log.Error("Error during converting competitor limit from string to int for Cube4Fun", "error", err)
			} else {
				dbItem.CompetitorLimit = i
				log.Info("Found Competitor limit for Cube4Fun", "competitorLimit", i)
			}

			i, err = strconv.Atoi(registered)
			if err != nil {
				log.Error("Error during converting registered competitors number from string to int for Cube4Fun", "error", err)
			} else {
				dbItem.Registered = i
				log.Info("Found Registered competitors number limit for Cube4Fun", "registered", i)
			}
		}
		if text == "events" || text == "konkurencje" {
			var eventsArr = []string{}
			events := row.Next()
			events.Find("title").Each(func(i int, row *goquery.Selection) {
				ev := trim(row.Text())
				eventId, exists := eventsMap.IdByName(ev)
				if exists {
					eventsArr = append(eventsArr, eventId)
				}
			})
			dbItem.Events = eventsArr
		}
	})
	log.Info("Successfully update Cube4Fun event, ready to export")
}

func trim(text string) string {
	text = strings.TrimSpace(text)
	text = strings.Trim(text, "\n")

	return text
}
