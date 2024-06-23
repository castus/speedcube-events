package diff

import (
	"github.com/castus/speedcube-events/db"
	"github.com/castus/speedcube-events/logger"
)

type Differences struct {
	Added   []string
	Passed  []string
	Changed []string
}

var log = logger.Default()

func (d Differences) IsEmpty() bool {
	return !d.HasAdded() && !d.HasPassed() && !d.HasChanged()
}
func (d Differences) HasAdded() bool {
	return len(d.Added) > 0
}
func (d Differences) HasPassed() bool {
	return len(d.Passed) > 0
}
func (d Differences) HasChanged() bool {
	return len(d.Changed) > 0
}

func Diff(localItemsDatabase *db.Database, database *db.Database) Differences {
	// Local (scraped) items are the source of truth
	var diff = Differences{
		Added:   []string{},
		Passed:  []string{},
		Changed: []string{},
	}

	for _, item := range localItemsDatabase.GetAll() {
		diffItem := database.Get(item.Id)
		if diffItem == nil {
			diff.Added = append(diff.Added, item.Id)
			continue
		}

		if !item.IsEqualTo(*diffItem) {
			diff.Changed = append(diff.Changed, item.Id)
		}
	}

	for _, item := range database.GetAll() {
		diffItem := localItemsDatabase.Get(item.Id)
		if diffItem == nil {
			if !item.HasPassed {
				diff.Passed = append(diff.Passed, item.Id)
				continue
			}
		}
	}

	return diff
}

func (d *Differences) PrintDifferencesInfo() {
	if d.HasChanged() {
		log.Info("Items to change", "length", len(d.Changed))
	}
	if d.HasAdded() {
		log.Info("Items to add", "length", len(d.Added))
	}
	if d.HasPassed() {
		log.Info("Items to mark as passed", "length", len(d.Passed))
	}
}
