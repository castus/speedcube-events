package db

import (
	"github.com/castus/speedcube-events/logger"
)

var log = logger.Default()

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

type Competition struct {
	Id              string
	Header          string
	Name            string
	URL             string
	Place           string
	LogoURL         string
	ContactName     string
	ContactURL      string
	HasWCA          bool // Will the score save in WCA
	HasPassed       bool // Event moved to Past tab
	Date            string
	Type            string
	TypeSpecificId  string
	WCAId           string // Legacy, use Type for that
	Distance        string
	Duration        string
	Events          []string // WCA API scrap
	MainEvent       string   // WCA GeneralInfo scrap
	CompetitorLimit int      // WCA GeneralInfo scrap
	Registered      int      // WCA Registrations scrap
}

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

func (c Competition) ExtractWCAId() string {
	if c.Type != CompetitionType.WCA {
		panic("This is not a WCA Type, please use WCA to get WCA ids")
	}

	var id string
	if c.WCAId != "" {
		id = c.WCAId
	} else {
		id = c.TypeSpecificId
	}

	return id
}
