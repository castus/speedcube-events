package messenger

import (
	"fmt"
	"log/slog"
	"net/smtp"
	"os"
	"strings"
)

var log = slog.New(slog.NewTextHandler(os.Stdout, nil))

type Mail struct {
	Sender  string
	Subject string
	Body    string
}

func buildMessage(mail Mail, to string) string {
	msg := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\r\n"
	msg += fmt.Sprintf("From: Zawody Speedcuberskie w Polsce <%s>\r\n", mail.Sender)
	msg += fmt.Sprintf("To: %s\r\n", to)
	msg += fmt.Sprintf("Subject: %s\r\n", mail.Subject)
	msg += fmt.Sprintf("\r\n%s\r\n", mail.Body)

	return msg
}

func Send(message string) {	
	sender := "kalendarz@krzysztofromanowski.pl"
	user := os.Getenv("MAIL_SMTP_USER")
	password := os.Getenv("MAIL_SMTP_PASSWORD")
	subject := "Zmiany w zawodach Speedcuberskich"
	body := message
	addr := "smtp-relay.brevo.com:587"
	host := "smtp-relay.brevo.com"

	request := Mail{
		Sender:  sender,
		Subject: subject,
		Body:    body,
	}

	tos := strings.Split(os.Getenv("NOTIFY_EMAILS"), ",")

	for _, email := range tos {
		msg := buildMessage(request, email)
		auth := smtp.PlainAuth("", user, password, host)
		err := smtp.SendMail(addr, auth, sender, []string{email}, []byte(msg))

		if err != nil {
			log.Error("Error sending email", "recipient", email, "trace", err)
		} else {
			log.Info("Email sent successfully", "recipient", email)
		}
	}
}
