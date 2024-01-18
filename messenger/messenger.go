package messenger

import (
	"fmt"
	"github.com/castus/speedcube-events/db"
	"log/slog"
	"os"
	"strings"

	"github.com/castus/speedcube-events/diff"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

var log = slog.New(slog.NewTextHandler(os.Stdout, nil))

func Send(IDs diff.Differences, competitions db.Competitions) {
	m := mail.NewV3Mail()
	from := mail.NewEmail(os.Getenv("MAIL_FROM_NAME"), os.Getenv("MAIL_FROM_EMAIL"))
	m.SetFrom(from)

	m.Subject = "Kalendarz imprez się zmienił"

	msg := prepareMessage(IDs, competitions)
	c := mail.NewContent("text/html", msg)
	m.AddContent(c)

	p := mail.NewPersonalization()
	var tos = []*mail.Email{}
	splitEmails := strings.Split(os.Getenv("NOTIFY_EMAILS"), ",")
	for _, item := range splitEmails {
		tos = append(tos, mail.NewEmail(item, item))
	}
	p.AddTos(tos...)
	m.AddPersonalizations(p)

	client := sendgrid.NewSendClient(os.Getenv("MAIL_TOKEN"))
	_, err := client.Send(m)
	if err != nil {
		log.Error("Couldn't send email", err)
	} else {
		log.Info("Email send")
	}
}

func prepareMessage(IDs diff.Differences, competitions db.Competitions) string {
	var message = []string{}

	message = append(message, "Kalendarz imprez się zmienił.")
	if IDs.HasAdded() {
		message = append(message, "<h1 style=\"margin: 40px 0 20px; font-weight: normal\">Imprezy dodane</h1>")
		for _, item := range IDs.Added {
			message = append(message, FormattedItemAsHTML(*competitions.FindByID(item)))
		}
	}
	if IDs.HasChanged() {
		message = append(message, "<h1 style=\"margin: 40px 0 20px; font-weight: normal\">Zmiany w istniejących imprezach</h1>")
		for _, item := range IDs.Changed {
			message = append(message, FormattedItemAsHTML(*competitions.FindByID(item)))
		}
	}
	if IDs.HasRemoved() {
		message = append(message, "<h1 style=\"margin: 40px 0 20px; font-weight: normal\">Imprezy usunięte</h1>")
		for _, item := range IDs.Removed {
			message = append(message, FormattedItemAsHTML(*competitions.FindByID(item)))
		}
	}

	return strings.Join(message, "\n")
}

func FormattedItemAsHTML(c db.Competition) string {
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

	if len(c.Events) > 0 {
		message = append(message, fmt.Sprintf("<p style=\"margin: 0 0 3px\">Konkurencje: "))
		for _, event := range c.Events {
			message = append(message, cubeImageForTag(event))
		}
		message = append(message, fmt.Sprintf("</p>"))
	}
	if c.MainEvent != "" {
		message = append(message, fmt.Sprintf("<p style=\"margin: 0 0 3px\">Konkurencja główna: %s</p>", cubeImageForTag(c.MainEvent)))
	}
	if c.CompetitorLimit > 0 {
		message = append(message, fmt.Sprintf("<p style=\"margin: 0 0 3px\">Zarejestrowanych: %d/%d</p>", c.Registered, c.CompetitorLimit))
	}

	message = append(message, fmt.Sprintf("<p style=\"margin: 0 0 3px\"><a href=\"mailto:%s\">%s (%s)</a></p>", c.ContactURL, c.ContactName, c.ContactURL))
	message = append(message, fmt.Sprintf("</td>"))
	message = append(message, fmt.Sprintf("</tr>"))
	message = append(message, "</table>")

	return strings.Join(message, "\n")
}

func cubeImageForTag(tag string) string {
	return fmt.Sprintf("<img src=\"https://raw.githubusercontent.com/cubing/icons/main/svgs/event/%s.svg\" width=\"20\" height=\"20\" alt=\"%s\" />", tag, tag)
}
