package dataFetch

import (
	"fmt"
	"github.com/castus/speedcube-events/logger"
	"net/http"
	"regexp"
	"strings"
	"unicode"

	"github.com/PuerkitoBio/goquery"
	"github.com/castus/speedcube-events/db"
)

const (
	host       = "https://www.speedcubing.pl"
	eventsPath = "kalendarz-imprez"
)

var log = logger.Default()

func ScrapCompetitions() db.Competitions {
	var competitions db.Competitions

	doc, ok := fetchWebPageBody(fmt.Sprintf("%s/%s", host, eventsPath))
	if !ok {
		return competitions
	}

	doc.Find("#competitions-following .row").Each(func(i int, row *goquery.Selection) {
		header := row.Prev().Text()
		header = strings.TrimSpace(header)
		header = strings.Trim(header, "\n")
		row.Find("li.col-sm-6").Each(func(i int, s *goquery.Selection) {
			competition := db.Competition{}

			competition.Header = header

			title := s.Find(".ulr-title").Text()
			title = strings.TrimSpace(title)
			title = strings.Trim(title, "\n")
			competition.Name = title

			competition.Id = normalizeString(title)

			url, _ := s.Find(".ulr-title a").Attr("href")
			competition.URL = url
			if strings.Contains(url, WCAHost) {
				elements := strings.Split(url, "/")
				competition.WCAId = elements[len(elements)-1]
				log.Info("Found WCA URL, extracting ID.", "WCAId", competition.WCAId, "URL", competition.URL)
			}

			logo, _ := s.Find(".ulr-image").Attr("style")
			logoURL := logoURL(logo)
			competition.LogoURL = fmt.Sprintf("%s/%s", host, logoURL)

			dateFrom := s.Find(".ulr-title").Next().Text()
			dateFrom = strings.TrimSpace(dateFrom)
			competition.Date = dateFrom

			contact := s.Find(".ulr-contact a").Text()
			contact = strings.TrimSpace(contact)
			contact = strings.Trim(contact, "\n")
			contact = strings.ReplaceAll(contact, "@ ", "")
			competition.ContactName = contact

			contactURL, _ := s.Find(".ulr-contact a").Attr("href")
			contactURL = strings.ReplaceAll(contactURL, "mailto:", "")
			competition.ContactURL = contactURL

			competition.HasWCA = s.Has(".ulr-title img").Size() > 0

			place := s.Find(".ulr-contact").Prev().Text()
			place = strings.TrimSpace(place)
			place = strings.Trim(place, "\n")
			competition.Place = place

			competitions = append(competitions, competition)
		})
	})

	return competitions
}

func fetchWebPageBody(URL string) (*goquery.Document, bool) {
	res, err := http.Get(URL)
	if err != nil {
		log.Error("Couldn't fetch page to scrap", "error", err, "url", URL)
		return nil, false
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Error("Status code error", "status code", res.StatusCode, "status", res.Status)
		return nil, false
	}
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Error("Couldn't load HTML", "error", err, "url", URL)
		return nil, false
	}

	return doc, true
}

func normalizeString(s string) string {
	s = strings.ToLower(s)

	var polishReplacements = map[rune]rune{
		'ą': 'a',
		'ć': 'c',
		'ę': 'e',
		'ł': 'l',
		'ń': 'n',
		'ó': 'o',
		'ś': 's',
		'ż': 'z',
		'ź': 'z',
	}
	var b strings.Builder
	for _, r := range s {
		if repl, ok := polishReplacements[r]; ok {
			b.WriteRune(repl)
		} else {
			b.WriteRune(r)
		}
	}
	s = b.String()

	s = strings.ReplaceAll(s, " ", "-")

	s = strings.Map(func(r rune) rune {
		if unicode.IsLetter(r) || unicode.IsNumber(r) || r == '-' {
			return r
		}
		return -1
	}, s)

	return s
}

func logoURL(input string) string {
	pattern := regexp.MustCompile(`\('([^']*)'\)`)
	matches := pattern.FindStringSubmatch(input)

	if len(matches) < 2 {
		return ""
	}

	return matches[1]
}
