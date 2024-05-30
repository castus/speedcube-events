package db

import (
	"github.com/castus/speedcube-events/logger"
)

type Types struct {
	Unknown  string
	WCA      string
	Cube4Fun string
	PPO      string
}

var CompetitionType = Types{
	Unknown:  "Unknown",
	WCA:      "WCA",
	Cube4Fun: "Cube4Fun",
	PPO:      "PPO",
}

var PageTypes = struct {
	Info        string
	Competitors string
}{
	Info:        "info",
	Competitors: "competitors",
}

type Competition struct {
	Type                              string
	TypeSpecificId                    string
	WCAId                             string // Legacy, use Type for that
	Id                                string
	Header, Name, URL, Place, LogoURL string
	ContactName, ContactURL           string
	Date                              string
	Distance                          string
	Duration                          string
	HasWCA                            bool     // Does the score will save in WCA
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
