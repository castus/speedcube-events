package db

import (
	"github.com/castus/speedcube-events/logger"
)

type Competition struct {
	WCAId                             string
	Id                                string
	Header, Name, URL, Place, LogoURL string
	ContactName, ContactURL           string
	HasWCA                            bool
	Date                              string
	Distance                          string
	Duration                          string
	HasPassed                         bool     // Event moved to Past tab
	Events                            []string // WCA API scrap
	MainEvent                         string   // WCA GeneralInfo scrap
	CompetitorLimit                   int      // WCA GeneralInfo scrap
	Registered                        int      // WCA Registrations scrap
}

var log = logger.Default()

func (c Competition) IsEqualTo(competition Competition) bool {
	if c.Id == competition.Id &&
		c.Header == competition.Header &&
		c.Name == competition.Name &&
		c.URL == competition.URL &&
		c.Place == competition.Place &&
		c.LogoURL == competition.LogoURL &&
		c.ContactName == competition.ContactName &&
		c.ContactURL == competition.ContactURL &&
		c.HasWCA == competition.HasWCA &&
		c.HasPassed == competition.HasPassed &&
		c.Date == competition.Date {
		return true
	}

	log.Debug("Item changed", "from", c, "to", competition)
	return false
}

func (c Competition) HasWCAPage() bool {
	return c.WCAId != ""
}
