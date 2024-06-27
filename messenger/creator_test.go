package messenger

import (
	"testing"
)

func TestCreating(t *testing.T) {
	mockTeamplates()
	m := PrepareMessage([]MessengerDTO{}, []MessengerDTO{}, []MessengerDTO{})

	if m == "\n\n" {
		t.Fatalf(`Expect empty string, git got %s`, m)
	}
}

func TestLogoImageCreation(t *testing.T) {
	out := logoHTML("http://google.com/img.jpg")
	if out != `<img src="http://google.com/img.jpg" width="100" />` {
		t.Fatalf(`Expect image string, git got %s`, out)
	}
}

func TestPlaceBuilding(t *testing.T) {
	out := placeHTML("Warszawa", "0km", "1h")
	if out != `<p>Warszawa <small>0km, 1h jazdy autem</small></p>` {
		t.Fatalf(`Expect place string, git got %s`, out)
	}
}

func TestPlaceOnline(t *testing.T) {
	out := placeHTML("zawody online", "0km", "1h")
	if out != `<p>zawody online</p>` {
		t.Fatalf(`Expect place string, git got %s`, out)
	}
}

func TestEvents(t *testing.T) {
	out := eventsHTML([]string{"333", "222"})
	if out != `<p>Konkurencje: <img src="./333.svg" alt="333"/> <img src="./222.svg" alt="222"/> </p>` {
		t.Fatalf(`Expect events string, git got %s`, out)
	}
}

func mockTeamplates() {
	headerTemplate = ""
	footerTemplate = ""
	logoTemplate = `<img src="{{.}}" width="100" />`
	placeTemplate = `<p>{{.Place}}{{if ne .Place "zawody online"}} <small>{{.Distance}}, {{.Duration}} jazdy autem</small>{{end}}</p>`
	eventsTemplate = `{{if ne (len .) 0}}<p>Konkurencje: {{range .}}<img src="./{{.}}.svg" alt="{{.}}"/> {{end}}{{end}}</p>`
}
