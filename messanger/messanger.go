package messanger

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/castus/speedcube-events/diff"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

func Send(competitions diff.Differences) {
	msg := prepareMessage(competitions)
	from := mail.NewEmail("Krzysztof Romanowski", "kontakt@krzysztofromanowski.pl")
	subject := "Kalendarz imprez się zmienił"
	to := mail.NewEmail("Krzysztof Romanowski", "castus.pl@gmail.com")
	plainTextContent := ""
	htmlContent := msg
	message := mail.NewSingleEmail(from, subject, to, plainTextContent, htmlContent)
	client := sendgrid.NewSendClient(os.Getenv("MAIL_TOKEN"))
	response, err := client.Send(message)
	if err != nil {
		log.Println(err)
	} else {
		fmt.Println(response.StatusCode)
		fmt.Println(response.Body)
		fmt.Println(response.Headers)
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
