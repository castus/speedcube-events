package diff

import (
	"testing"

	"github.com/castus/speedcube-events/db"
)

func TestTheSame(t *testing.T) {
	diff := runDiff([]db.Competition{itemOne}, []db.Competition{itemOne})
	if !diff.IsEmpty() {
		t.Error(`Compared items are not the same`)
	}
}
func TestAdded(t *testing.T) {
	diff := runDiff([]db.Competition{itemOne}, []db.Competition{})
	if !diff.HasAdded() {
		t.Error(`There are no added items during difference check`)
	}
}
func TestRemoved(t *testing.T) {
	diff := runDiff([]db.Competition{}, []db.Competition{itemOne})
	if !diff.HasPassed() {
		t.Error(`There are no removed items during difference check`)
	}
}
func TestChangesType(t *testing.T) {
	itemTwo := itemOne
	itemTwo.Type = "New Type"
	diff := runDiff([]db.Competition{itemOne}, []db.Competition{itemTwo})
	if !diff.HasChanged() {
		t.Error(`There are changed items during difference check`)
	}
}
func TestChangesTypeSpecificId(t *testing.T) {
	itemTwo := itemOne
	itemTwo.TypeSpecificId = "New TypeSpecificId"
	diff := runDiff([]db.Competition{itemOne}, []db.Competition{itemTwo})
	if !diff.HasChanged() {
		t.Error(`There are changed items during difference check`)
	}
}
func TestNoChangesWCAId(t *testing.T) {
	itemTwo := itemOne
	itemTwo.WCAId = "New WCAId"
	diff := runDiff([]db.Competition{itemOne}, []db.Competition{itemTwo})
	if diff.HasChanged() {
		t.Error(`There are changed items during difference check`)
	}
}
func TestNoChangesDistance(t *testing.T) {
	itemTwo := itemOne
	itemTwo.Distance = "New Distance"
	diff := runDiff([]db.Competition{itemOne}, []db.Competition{itemTwo})
	if diff.HasChanged() {
		t.Error(`There are changed items during difference check`)
	}
}
func TestNoChangesDuration(t *testing.T) {
	itemTwo := itemOne
	itemTwo.Duration = "New Duration"
	diff := runDiff([]db.Competition{itemOne}, []db.Competition{itemTwo})
	if diff.HasChanged() {
		t.Error(`There are changed items during difference check`)
	}
}
func TestNoChangesEvents(t *testing.T) {
	itemTwo := itemOne
	itemTwo.Events = []string{"New Events"}
	diff := runDiff([]db.Competition{itemOne}, []db.Competition{itemTwo})
	if diff.HasChanged() {
		t.Error(`There are changed items during difference check`)
	}
}
func TestNoChangesMainEvent(t *testing.T) {
	itemTwo := itemOne
	itemTwo.MainEvent = "New MainEvent"
	diff := runDiff([]db.Competition{itemOne}, []db.Competition{itemTwo})
	if diff.HasChanged() {
		t.Error(`There are changed items during difference check`)
	}
}
func TestNoChangesCompetitorLimit(t *testing.T) {
	itemTwo := itemOne
	itemTwo.CompetitorLimit = 123
	diff := runDiff([]db.Competition{itemOne}, []db.Competition{itemTwo})
	if diff.HasChanged() {
		t.Error(`There are changed items during difference check`)
	}
}
func TestNoChangesRegistered(t *testing.T) {
	itemTwo := itemOne
	itemTwo.Registered = 123
	diff := runDiff([]db.Competition{itemOne}, []db.Competition{itemTwo})
	if diff.HasChanged() {
		t.Error(`There are changed items during difference check`)
	}
}
func TestChangesName(t *testing.T) {
	itemTwo := itemOne
	itemTwo.Name = "New Name"
	diff := runDiff([]db.Competition{itemOne}, []db.Competition{itemTwo})
	if !diff.HasChanged() {
		t.Error(`There are changed items during difference check`)
	}
}
func TestChangesURL(t *testing.T) {
	itemTwo := itemOne
	itemTwo.URL = "New URL"
	diff := runDiff([]db.Competition{itemOne}, []db.Competition{itemTwo})
	if !diff.HasChanged() {
		t.Error(`There are changed items during difference check`)
	}
}
func TestChangesPlace(t *testing.T) {
	itemTwo := itemOne
	itemTwo.Place = "New Place"
	diff := runDiff([]db.Competition{itemOne}, []db.Competition{itemTwo})
	if !diff.HasChanged() {
		t.Error(`There are changed items during difference check`)
	}
}
func TestChangesLogoURL(t *testing.T) {
	itemTwo := itemOne
	itemTwo.LogoURL = "New LogoURL"
	diff := runDiff([]db.Competition{itemOne}, []db.Competition{itemTwo})
	if !diff.HasChanged() {
		t.Error(`There are changed items during difference check`)
	}
}
func TestChangesContactName(t *testing.T) {
	itemTwo := itemOne
	itemTwo.ContactName = "New ContactName"
	diff := runDiff([]db.Competition{itemOne}, []db.Competition{itemTwo})
	if !diff.HasChanged() {
		t.Error(`There are changed items during difference check`)
	}
}
func TestChangesContactURL(t *testing.T) {
	itemTwo := itemOne
	itemTwo.ContactURL = "New ContactURL"
	diff := runDiff([]db.Competition{itemOne}, []db.Competition{itemTwo})
	if !diff.HasChanged() {
		t.Error(`There are changed items during difference check`)
	}
}
func TestChangesHasWCA(t *testing.T) {
	itemTwo := itemOne
	itemTwo.HasWCA = false
	diff := runDiff([]db.Competition{itemOne}, []db.Competition{itemTwo})
	if !diff.HasChanged() {
		t.Error(`There are changed items during difference check`)
	}
}
func TestChangesHasPassed(t *testing.T) {
	itemTwo := itemOne
	itemTwo.HasPassed = true
	diff := runDiff([]db.Competition{itemOne}, []db.Competition{itemTwo})
	if !diff.HasChanged() {
		t.Error(`There are changed items during difference check`)
	}
}
func TestChangesDate(t *testing.T) {
	itemTwo := itemOne
	itemTwo.Date = "New Date"
	diff := runDiff([]db.Competition{itemOne}, []db.Competition{itemTwo})
	if !diff.HasChanged() {
		t.Error(`There are changed items during difference check`)
	}
}
func TestChangesHeader(t *testing.T) {
	itemTwo := itemOne
	itemTwo.Header = "New Header"
	diff := runDiff([]db.Competition{itemOne}, []db.Competition{itemTwo})
	if !diff.HasChanged() {
		t.Error(`There are changed items during difference check`)
	}
}
func TestChanges(t *testing.T) {
	itemTwo := itemOne
	itemTwo.Name = "Competition 2 Name"
	diff := runDiff([]db.Competition{itemOne}, []db.Competition{itemTwo})
	if !diff.HasChanged() {
		t.Error(`There are no changed items during difference check`)
	}
}

func runDiff(local []db.Competition, database []db.Competition) Differences {
	items1Database := db.InitializeWith(local)
	items2Database := db.InitializeWith(database)

	return Diff(&items1Database, &items2Database)
}

var itemOne = db.Competition{
	Type:            "Competition 1 Type",
	TypeSpecificId:  "Competition 1 TypeSpecificId",
	WCAId:           "Competition 1 WCAId",
	Id:              "competition-1-id",
	Header:          "Competition 1 Header",
	Name:            "Competition 1 Name",
	URL:             "Competition 1 URL",
	Place:           "Competition 1 Place",
	LogoURL:         "Competition 1 LogoURL",
	ContactName:     "Competition 1 ContactName",
	ContactURL:      "Competition 1 ContactURL",
	Date:            "Competition 1 Date",
	Distance:        "Competition 1 Distance",
	Duration:        "Competition 1 Duration",
	MainEvent:       "Competition 1 MainEvent",
	HasWCA:          true,
	HasPassed:       false,
	Events:          []string{"333", "444"},
	CompetitorLimit: 150,
	Registered:      40,
}
