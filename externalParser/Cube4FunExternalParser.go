package externalParser

import (
	"fmt"
	"io"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/castus/speedcube-events/dataFetch"
	"github.com/castus/speedcube-events/db"
)

var eventNames = make(map[string]string)

func Cube4FunParse(body io.Reader, competitionType string, id string, pageName string, eventsMap []dataFetch.EventMap) {
	log.Info("Parsing Cube4Fun", "type", competitionType, "id", id, "pageName", pageName)
	pageNameItems := strings.Split(pageName, ".")
	pageKey := pageNameItems[0]
	prepareEventsMap(eventsMap)
	if pageKey == db.Cube4FunPages.Info {
		parseInfo(body)
	}
}

func prepareEventsMap(eventsMap []dataFetch.EventMap) {
	for _, event := range eventsMap {
		eventNames[event.Name] = event.Id
	}
}

func parseInfo(body io.Reader) {
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
			registered := limits[0]
			all := limits[1]

			fmt.Println(registered)
			fmt.Println(all)
		}
		if text == "events" || text == "konkurencje" {
			var eventsArr = []string{}
			events := row.Next()
			events.Find("title").Each(func(i int, row *goquery.Selection) {
				ev := trim(row.Text())
				eventsArr = append(eventsArr, eventNames[ev])
			})
			fmt.Println(eventsArr)
		}
	})
}

func trim(text string) string {
	text = strings.TrimSpace(text)
	text = strings.Trim(text, "\n")

	return text
}
