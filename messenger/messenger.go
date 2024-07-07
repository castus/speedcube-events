package messenger

import (
	"log/slog"
	"os"
	"strings"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

var log = slog.New(slog.NewTextHandler(os.Stdout, nil))

func Send(message string) {
	m := mail.NewV3Mail()
	from := mail.NewEmail(os.Getenv("MAIL_FROM_NAME"), os.Getenv("MAIL_FROM_EMAIL"))
	m.SetFrom(from)

	m.Subject = "Kalendarz imprez się zmienił"

	c := mail.NewContent("text/html", message)
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
		log.Error("Couldn't send email", "error", err)
	} else {
		log.Info("Email send")
	}
}
