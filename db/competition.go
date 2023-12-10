package db

import (
	"fmt"
	"strings"

	"github.com/castus/speedcube-events/logger"
)

type Competition struct {
	Id                                string
	Header, Name, URL, Place, LogoURL string
	ContactName, ContactURL           string
	HasWCA                            bool
	Date                              string
	Distance                          string
	Duration                          string
}

var log = logger.Default()

func (c Competition) IsEqualTo(competition Competition) bool {
	if c.Id != competition.Id &&
		c.Header != competition.Header &&
		c.Name != competition.Name &&
		c.URL != competition.URL &&
		c.Place != competition.Place &&
		c.LogoURL != competition.LogoURL &&
		c.ContactName != competition.ContactName &&
		c.ContactURL != competition.ContactURL &&
		c.HasWCA != competition.HasWCA &&
		c.Date != competition.Date {
		log.Debug("Item changed", "from", c, "to", competition)
		return false
	}

	return true
}

func (c Competition) PrintHTMLContent() string {
	var message = []string{}

	message = append(message, "<table border=\"1\" cellpadding=\"10px\" style=\"margin: 0; border-collapse: collapse; width: 100%;\">")

	header := fmt.Sprintf("%s <small>(%s)</small>", c.Name, c.Header)
	if len(c.URL) > 0 {
		header = fmt.Sprintf("<a href=\"%s\">%s</a>", c.URL, header)
	}
	if c.HasWCA {
		header = fmt.Sprintf("%s <img src=\"https://www.speedcubing.pl/images/wca_small_logo.png\" width=\"30\" height=\"30\" />", header)
	}
	message = append(message, fmt.Sprintf("<tr><td colspan=\"2\"><h2 style=\"margin: 0; font-weight: normal\">%s</h2></td></tr>", header))

	message = append(message, fmt.Sprintf("<tr><td style=\"width: 120px; text-align: center\"><img src=\"%s\" width=\"100\" height=\"100\" /></td>", c.LogoURL))
	message = append(message, fmt.Sprintf("<td valign=\"top\">"))
	message = append(message, fmt.Sprintf("<p style=\"margin: 0 0 3px\">%s</p>", c.Date))

	placeMessage := c.Place
	if c.Place != "zawody online" {
		placeMessage = fmt.Sprintf("%s, <small>%s, %s jazdy autem</small>", placeMessage, c.Distance, c.Duration)
	}
	message = append(message, fmt.Sprintf("<p style=\"margin: 0 0 3px\">%s</p>", placeMessage))

	message = append(message, fmt.Sprintf("<p style=\"margin: 0 0 3px\"><a href=\"mailto:%s\">%s (%s)</a></p>", c.ContactURL, c.ContactName, c.ContactURL))
	message = append(message, fmt.Sprintf("</td>"))
	message = append(message, fmt.Sprintf("</tr>"))
	message = append(message, "</table>")

	return strings.Join(message, "\n")
}
