package diff

import (
	"github.com/castus/speedcube-events/db"
)

type Differences struct {
	Added   []db.Competition
	Removed []db.Competition
	Changed []db.Competition
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

func Diff(local []db.Competition, database []db.Competition) Differences {
	// Local items are source of truth
	var diff = Differences{
		Added:   []db.Competition{},
		Removed: []db.Competition{},
		Changed: []db.Competition{},
	}

	for _, item := range local {
		diffItem := findById(item.Id, database)
		if diffItem == nil {
			diff.Added = append(diff.Added, item)
			continue
		}

		if !item.IsEqualTo(*diffItem) {
			diff.Changed = append(diff.Changed, item)
		}
	}

	for _, item := range database {
		diffItem := findById(item.Id, local)
		if diffItem == nil {
			diff.Removed = append(diff.Removed, item)
			continue
		}
	}

	return diff
}

func findById(id string, competitions []db.Competition) *db.Competition {
	for _, item := range competitions {
		if item.Id == id {
			return &item
		}
	}

	return nil
}
