package db

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/castus/speedcube-events/distance"
)

type Competition struct {
	Id                                string
	Header, Name, URL, Place, LogoURL string
	ContactName, ContactURL           string
	HasWCA                            bool
	Date                              string
}

func (c Competition) IsEqualTo(competition Competition) bool {
	if reflect.DeepEqual(c, competition) {
		return true
	} else {
		fmt.Println("BYŁO:")
		fmt.Println(c)
		fmt.Println("JEST:")
		fmt.Println(competition)
		return false
	}
}

func (c Competition) PrintHTMLContent() string {
	var message = []string{}

	message = append(message, "<table border=\"1\" cellpadding=\"10px\" style=\"margin: 0; border-collapse: collapse; width: 100%;\">")
	if len(c.URL) > 0 {
		message = append(message, fmt.Sprintf("<tr><td colspan=\"2\"><h2 style=\"margin: 0; font-weight: normal\"><a href=\"%s\">%s <small>(%s)</small></a></h2></td></tr>", c.URL, c.Name, c.Header))
	} else {
		message = append(message, fmt.Sprintf("<tr><td colspan=\"2\"><h2 style=\"margin: 0; font-weight: normal\">%s <small>(%s)</small></h2></td></tr>", c.Name, c.Header))
	}

	message = append(message, fmt.Sprintf("<tr><td style=\"width: 120px; text-align: center\"><img src=\"%s\" width=\"100\" height=\"100\" /></td>", c.LogoURL))
	message = append(message, fmt.Sprintf("<td valign=\"top\">"))
	message = append(message, fmt.Sprintf("<p style=\"margin: 0 0 3px\">%s</p>", c.Date))

	placeMessage := c.Place
	if c.Place != "zawody online" {
		travelInfo, err := distance.Distance(c.Place)
		if err == nil {
			placeMessage = fmt.Sprintf("%s, <small>%s, %s jazdy autem</small>", placeMessage, travelInfo.Distance, travelInfo.Duration)
		}
	}
	message = append(message, fmt.Sprintf("<p style=\"margin: 0 0 3px\">%s</p>", placeMessage))

	message = append(message, fmt.Sprintf("<p style=\"margin: 0 0 3px\"><a href=\"mailto:%s\">%s (%s)</a></p>", c.ContactURL, c.ContactName, c.ContactURL))
	message = append(message, fmt.Sprintf("</td>"))
	message = append(message, fmt.Sprintf("</tr>"))
	message = append(message, "</table>")

	return strings.Join(message, "\n")
}
