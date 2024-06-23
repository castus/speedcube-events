package diff

import (
	"testing"

	"github.com/castus/speedcube-events/db"
)

func TestTheSame(t *testing.T) {
	diff := runDiff([]db.Competition{itemOne}, []db.Competition{itemOne})
	if !diff.IsEmpty() {
		t.Fatalf(`Compared items are not the same`)
	}
}
func TestAdded(t *testing.T) {
	diff := runDiff([]db.Competition{itemOne}, []db.Competition{})
	if !diff.HasAdded() {
		t.Fatalf(`There are no added items during difference check`)
	}
}
func TestRemoved(t *testing.T) {
	diff := runDiff([]db.Competition{}, []db.Competition{itemOne})
	if !diff.HasPassed() {
		t.Fatalf(`There are no removed items during difference check`)
	}
}
func TestNoChanges(t *testing.T) {
	itemTwo := itemOne
	itemTwo.Registered = 1
	diff := runDiff([]db.Competition{itemOne}, []db.Competition{itemTwo})
	if diff.HasChanged() {
		t.Fatalf(`There are changed items during difference check`)
	}
}
func TestChanges(t *testing.T) {
	itemTwo := itemOne
	itemTwo.Name = "Competition 2 Name"
	diff := runDiff([]db.Competition{itemOne}, []db.Competition{itemTwo})
	if !diff.HasChanged() {
		t.Fatalf(`There are no changed items during difference check`)
	}
}

func runDiff(local []db.Competition, database []db.Competition) Differences {
	items1Database := db.Database{}
	items1Database.InitializeWith(local)
	items2Database := db.Database{}
	items2Database.InitializeWith(database)

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
