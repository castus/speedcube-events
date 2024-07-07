package dataFetch

import (
	"fmt"
	"io"
	"regexp"
	"strings"
	"unicode"

	"github.com/castus/speedcube-events/logger"

	"github.com/PuerkitoBio/goquery"
	"github.com/castus/speedcube-events/db"
)

const (
	registrationsPath = "registrations"
	WCAHost           = "worldcubeassociation.org"
	Cube4FunHost      = "cube4fun.pl"
	RubiArtHost       = "rubiart.pl"
	Host              = "https://www.speedcubing.pl"
	EventsPath        = "kalendarz-imprez"
)

var log = logger.Default()

type DataFetcher interface {
	Fetch(URL string) (r io.Reader, ok bool)
}

func ScrapCompetitions(fetcher DataFetcher) []db.Competition {
	var competitions []db.Competition

	URL := fmt.Sprintf("%s/%s", Host, EventsPath)
	r, ok := fetcher.Fetch(URL)
	if !ok {
		return competitions
	}

	doc, err := goquery.NewDocumentFromReader(r)
	if err != nil {
		log.Error("Couldn't parse HTML with GoQuery", "error", err, "url", URL)
		return competitions
	}

	log.Info("Parsing kalendarz imprez...")
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
			competitionType, typeSpecificId := getTypeInformation(competition.URL, competition.Name)
			competition.Type = competitionType
			competition.TypeSpecificId = typeSpecificId
			if competitionType != db.CompetitionType.Unknown {
				log.Info("Found known competition type.", "Type", competition.Type, "TypeID", competition.TypeSpecificId)
			} else {
				log.Info("Found unknown competition type.", "Id", competition.Id, "URL", url)
			}

			logo, _ := s.Find(".ulr-image").Attr("style")
			logoURL := logoURL(logo)
			competition.LogoURL = fmt.Sprintf("%s/%s", Host, logoURL)

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

func getTypeInformation(url string, name string) (string, string) {
	if strings.Contains(url, WCAHost) {
		elements := strings.Split(url, "/")
		return db.CompetitionType.WCA, elements[len(elements)-1]
	}
	if strings.Contains(url, Cube4FunHost) {
		elements := strings.Split(url, "/")
		return db.CompetitionType.Cube4Fun, elements[len(elements)-1]
	}
	if strings.Contains(url, RubiArtHost) && strings.Contains(name, "PPO") {
		elements := strings.Split(url, "/")
		return db.CompetitionType.PPO, elements[len(elements)-1]
	}

	return db.CompetitionType.Unknown, ""
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
