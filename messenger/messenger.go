package messenger

import (
	"log/slog"
	"os"
	"strings"

	"github.com/castus/speedcube-events/diff"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

var log = slog.New(slog.NewTextHandler(os.Stdout, nil))

func Send(competitions diff.Differences) {
	m := mail.NewV3Mail()
	from := mail.NewEmail(os.Getenv("MAIL_FROM_NAME"), os.Getenv("MAIL_FROM_EMAIL"))
	m.SetFrom(from)

	m.Subject = "Kalendarz imprez się zmienił"

	msg := prepareMessage(competitions)
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

func prepareMessage(diff diff.Differences) string {
	var message = []string{}

	message = append(message, "Kalendarz imprez się zmienił.")
	if diff.HasAdded() {
		message = append(message, "<h1 style=\"margin: 40px 0 20px; font-weight: normal\">Imprezy dodane</h1>")
		for _, item := range diff.Added {
			message = append(message, item.PrintHTMLContent())
		}
	}
	if diff.HasChanged() {
		message = append(message, "<h1 style=\"margin: 40px 0 20px; font-weight: normal\">Zmiany w istniejących imprezach</h1>")
		for _, item := range diff.Changed {
			message = append(message, item.PrintHTMLContent())
		}
	}
	if diff.HasRemoved() {
		message = append(message, "<h1 style=\"margin: 40px 0 20px; font-weight: normal\">Imprezy usunięte</h1>")
		for _, item := range diff.Removed {
			message = append(message, item.PrintHTMLContent())
		}
	}

	return strings.Join(message, "\n")
}
