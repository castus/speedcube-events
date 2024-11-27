package messenger

import (
	"fmt"
	"strings"
	"testing"
)

var commonItem = MessengerDTO{
	LogoURL:         "http://logo.url",
	URL:             "http://event.url",
	Name:            "Event Name",
	HasWCA:          true,
	Date:            "12-13 August",
	Place:           "Warsaw",
	Distance:        "123km",
	Duration:        "1h 10min",
	Events:          []string{"333", "222"},
	MainEvent:       "333",
	CompetitorLimit: 120,
	Registered:      10,
	ContactURL:      "contact@email.com",
	ContactName:     "Contact Name",
}

func TestSection_added(t *testing.T) {
	out := PrepareMessage([]MessengerDTO{commonItem}, []MessengerDTO{}, []MessengerDTO{})
	expectedSectionHeaders := 1
	sectionHeaders := strings.Count(out, "section-header")
	if sectionHeaders != expectedSectionHeaders {
		t.Error(fmt.Sprintf("Expected %d section headers, got", expectedSectionHeaders), sectionHeaders)
	}
}

func TestSection_passed(t *testing.T) {
	out := PrepareMessage([]MessengerDTO{}, []MessengerDTO{commonItem}, []MessengerDTO{})
	expectedSectionHeaders := 1
	sectionHeaders := strings.Count(out, "section-header")
	if sectionHeaders != expectedSectionHeaders {
		t.Error(fmt.Sprintf("Expected %d section headers, got", expectedSectionHeaders), sectionHeaders)
	}
}

func TestSection_changed(t *testing.T) {
	out := PrepareMessage([]MessengerDTO{}, []MessengerDTO{}, []MessengerDTO{commonItem})
	expectedSectionHeaders := 1
	sectionHeaders := strings.Count(out, "section-header")
	if sectionHeaders != expectedSectionHeaders {
		t.Error(fmt.Sprintf("Expected %d section headers, got", expectedSectionHeaders), sectionHeaders)
	}
}

func TestSection_everyType(t *testing.T) {
	out := PrepareMessage([]MessengerDTO{commonItem}, []MessengerDTO{commonItem, commonItem}, []MessengerDTO{commonItem})
	expectedSectionHeaders := 3
	sectionHeaders := strings.Count(out, "section-header")
	if sectionHeaders != expectedSectionHeaders {
		t.Error(fmt.Sprintf("Expected %d section headers, got", expectedSectionHeaders), sectionHeaders)
	}
}

func TestCompetitionItems_forEverySection(t *testing.T) {
	out := PrepareMessage([]MessengerDTO{commonItem}, []MessengerDTO{commonItem, commonItem}, []MessengerDTO{commonItem})
	expectedCompetitionItems := 4
	competitionItems := strings.Count(out, "competition-item")
	if competitionItems != expectedCompetitionItems {
		t.Error(fmt.Sprintf("Expected %d competition items, got", expectedCompetitionItems), competitionItems)
	}
}

func TestCompetitionItems_added(t *testing.T) {
	out := PrepareMessage([]MessengerDTO{commonItem}, []MessengerDTO{}, []MessengerDTO{})
	expectedCompetitionItems := 1
	competitionItems := strings.Count(out, "competition-item")
	if competitionItems != expectedCompetitionItems {
		t.Error(fmt.Sprintf("Expected %d competition items, got", expectedCompetitionItems), competitionItems)
	}
}

func TestHeader_noURL(t *testing.T) {
	item := commonItem
	item.URL = ""
	out := PrepareMessage([]MessengerDTO{item}, []MessengerDTO{}, []MessengerDTO{})

	if strings.Contains(out, "header-with-link") {
		t.Error(`Expected lack of header-with-link but got one`)
	}
}

func TestHeader_withURL(t *testing.T) {
	item := commonItem
	item.URL = "http://google.com"
	out := PrepareMessage([]MessengerDTO{item}, []MessengerDTO{}, []MessengerDTO{})

	if !strings.Contains(out, "header-with-link") {
		t.Error(`Expected header-with-link but got one`)
	}
}

func TestHeader_noWCA(t *testing.T) {
	item := commonItem
	item.HasWCA = false
	out := PrepareMessage([]MessengerDTO{item}, []MessengerDTO{}, []MessengerDTO{})

	if strings.Contains(out, "header-wca") {
		t.Error(`Expected lack of header-wca but got one`)
	}
}

func TestHeader_withWCA(t *testing.T) {
	item := commonItem
	item.HasWCA = true
	out := PrepareMessage([]MessengerDTO{item}, []MessengerDTO{}, []MessengerDTO{})

	if !strings.Contains(out, "header-wca") {
		t.Error(`Expected header-wca but got one`)
	}
}

func TestMainEvent_empty(t *testing.T) {
	item := commonItem
	item.MainEvent = ""
	out := PrepareMessage([]MessengerDTO{item}, []MessengerDTO{}, []MessengerDTO{})

	if strings.Contains(out, "main-event") {
		t.Error(`Expected lack of main-event but got one`)
	}
}

func TestMainEvent_present(t *testing.T) {
	item := commonItem
	item.MainEvent = "333"
	out := PrepareMessage([]MessengerDTO{item}, []MessengerDTO{}, []MessengerDTO{})

	if !strings.Contains(out, "main-event") {
		t.Error(`Expected main-event but got one`)
	}
}

func TestRegistered_empty(t *testing.T) {
	item := commonItem
	item.CompetitorLimit = 0
	out := PrepareMessage([]MessengerDTO{item}, []MessengerDTO{}, []MessengerDTO{})

	if strings.Contains(out, "registered") {
		t.Error(`Expected lack of registered but got one`)
	}
}

func TestRegistered_present(t *testing.T) {
	item := commonItem
	item.CompetitorLimit = 100
	out := PrepareMessage([]MessengerDTO{item}, []MessengerDTO{}, []MessengerDTO{})

	if !strings.Contains(out, "registered") {
		t.Error(`Expected registered but got one`)
	}
}

func TestRegistered_limitReached(t *testing.T) {
	item := commonItem
	item.CompetitorLimit = 100
	item.Registered = 100
	out := PrepareMessage([]MessengerDTO{item}, []MessengerDTO{}, []MessengerDTO{})

	if !strings.Contains(out, "registered") && !strings.Contains(out, "registered-limit") {
		t.Error(`Expected registered and registered-limit but got one`)
	}
}

func TestEvents_empty(t *testing.T) {
	item := commonItem
	item.Events = []string{}
	out := PrepareMessage([]MessengerDTO{item}, []MessengerDTO{}, []MessengerDTO{})

	if strings.Contains(out, "Konkurencje") {
		t.Error(`Expected lack of events but got one`)
	}
}

func TestEvents_present(t *testing.T) {
	item := commonItem
	item.Events = []string{"333", "222"}
	out := PrepareMessage([]MessengerDTO{item}, []MessengerDTO{}, []MessengerDTO{})

	if !strings.Contains(out, "events") && !strings.Contains(out, "event-image-333") && !strings.Contains(out, "event-image-222") {
		t.Error(`Expected events, event-image-222, event-image-333 but got one`)
	}
}

func TestDistanceInfo_empty(t *testing.T) {
	item := commonItem
	item.Distance = "1km"
	item.Duration = "1h"
	item.Place = "zawody online"
	out := PrepareMessage([]MessengerDTO{item}, []MessengerDTO{}, []MessengerDTO{})

	if strings.Contains(out, "distance-info") {
		t.Error(`Expected lack of distance-info but got one`)
	}
}

func TestDistanceInfo_present(t *testing.T) {
	item := commonItem
	item.Distance = "1km"
	item.Duration = "1h"
	item.Place = "Warsaw"
	out := PrepareMessage([]MessengerDTO{item}, []MessengerDTO{}, []MessengerDTO{})

	if !strings.Contains(out, "distance-info") {
		t.Error(`Expected distance-info but got one`)
	}
}
