package diff

import (
	"github.com/castus/speedcube-events/db"
)

type Differences struct {
	Added   []string
	Removed []string
	Changed []string
}

func (d Differences) IsEmpty() bool {
	return !d.HasAdded() && !d.HasRemoved() && !d.HasChanged()
}
func (d Differences) HasAdded() bool {
	return len(d.Added) > 0
}
func (d Differences) HasRemoved() bool {
	return len(d.Removed) > 0
}
func (d Differences) HasChanged() bool {
	return len(d.Changed) > 0
}

func Diff(local db.Competitions, database db.Competitions) Differences {
	// Local (scraped) items are the source of truth
	var diff = Differences{
		Added:   []string{},
		Removed: []string{},
		Changed: []string{},
	}

	for _, item := range local {
		diffItem := database.FindByID(item.Id)
		if diffItem == nil {
			diff.Added = append(diff.Added, item.Id)
			continue
		}

		if !item.IsEqualTo(*diffItem) {
			diff.Changed = append(diff.Changed, item.Id)
		}
	}

	for _, item := range database {
		diffItem := local.FindByID(item.Id)
		if diffItem == nil {
			diff.Removed = append(diff.Removed, item.Id)
			continue
		}
	}

	return diff
}
