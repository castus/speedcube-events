package externalParser

import (
	"github.com/castus/speedcube-events/externalFetcher"
	"io"

	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/castus/speedcube-events/dataFetch"
	"github.com/castus/speedcube-events/db"
)

func PPOParse(body io.Reader, competitionType string, id string, pageName string, eventsMap dataFetch.EventsMap, dbItem db.Competition, c *dynamodb.Client) {
	log.Info("Found PPO event, parsing ...", "type", competitionType, "id", id, "pageName", pageName)
	pageNameItems := strings.Split(pageName, ".")
	pageKey := pageNameItems[0]
	if pageKey == externalFetcher.PageTypes.Info {
		parsePPOInfo(body, dbItem, c, eventsMap)
	} else if pageKey == externalFetcher.PageTypes.Competitors {
		parsePPOCompetitors(body, dbItem, c)
	}
}

func parsePPOInfo(body io.Reader, dbItem db.Competition, c *dynamodb.Client, eventsMap dataFetch.EventsMap) {
	doc, err := goquery.NewDocumentFromReader(body)
	if err != nil {
		log.Error("Couldn't parse HTML with GoQuery", "error", err)
	}
	doc.Find(".MuiTypography-h5").Each(func(i int, row *goquery.Selection) {
		text := trim(row.Text())
		text = strings.ToLower(text)

		if text == "limit zawodników" {
			limit := trim(row.Next().Text())
			i, err := strconv.Atoi(limit)
			if err != nil {
				log.Error("Error during converting competitor limit from string to int for PPO", "error", err)
			} else {
				dbItem.CompetitorLimit = i
				log.Info("Found Competitor limit for PPO", "competitorLimit", i)
			}
		}

		if text == "konkurencje" {
			var eventsArr = []string{}
			events := row.Next()
			events.Find("img").Each(func(i int, row *goquery.Selection) {
				src, exists := row.Attr("src")
				if exists {
					sources := strings.Split(src, "/")
					last := sources[len(sources)-1]
					names := strings.Split(last, ".")
					competition := names[0]
					_, exists := eventsMap.NameById(competition)
					if exists {
						eventsArr = append(eventsArr, competition)
					}
				}
			})
			log.Info("Found competitions for PPO")
			dbItem.Events = eventsArr
		}
	})

	_, err = db.AddItemBatch(c, dbItem)
	if err != nil {
		log.Error("Couldn't save item to database after update PPO", "error", err)
		panic(err)
	}
	log.Info("Successfully update PPO event info page", "id", dbItem.Id)
}

func parsePPOCompetitors(body io.Reader, dbItem db.Competition, c *dynamodb.Client) {
	doc, err := goquery.NewDocumentFromReader(body)
	if err != nil {
		log.Error("Couldn't parse HTML with GoQuery", "error", err)
	}
	doc.Find(".MuiTable-root").Each(func(i int, sel *goquery.Selection) {
		rows := sel.Find("tbody tr")
		registered := rows.Length()
		dbItem.Registered = registered
		log.Info("Found Registered competitors for PPO", "registered", registered)
	})

	_, err = db.AddItemBatch(c, dbItem)
	if err != nil {
		log.Error("Couldn't save item to database after update PPO", "error", err)
		panic(err)
	}
	log.Info("Successfully update PPO event competitors page", "id", dbItem.Id)
}
