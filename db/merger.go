package db

type Merger struct {
	added   []string
	passed  []string
	changed []string
}

func NewMerger(added []string, passed []string, changed []string) *Merger {
	return &Merger{
		added:   added,
		passed:  passed,
		changed: changed,
	}
}

func (m *Merger) Merge(localDatabase Database, remoteDatabase Database) *Database {
	mergedDatabase := &remoteDatabase
	log.Info("merging events")
	if len(m.added) != 0 {
		addNew(m.added, localDatabase, mergedDatabase)
	}
	if len(m.changed) != 0 {
		updateChanged(m.changed, localDatabase, mergedDatabase)
	}
	if len(m.passed) != 0 {
		markAsPassed(m.passed, mergedDatabase)
	}

	return mergedDatabase
}

func addNew(added []string, from Database, to *Database) *Database {
	for _, addedItem := range added {
		to.Add(*from.Get(addedItem))
	}

	return to
}

func markAsPassed(passed []string, to *Database) *Database {
	for _, item := range passed {
		dbItem := to.Get(item)
		if dbItem.Id == item {
			dbItem.HasPassed = true
			to.Update(*dbItem)
		}
	}

	return to
}

func updateChanged(changed []string, from Database, to *Database) *Database {
	for _, id := range changed {
		dbItem := to.Get(id)
		localItem := from.Get(id)
		if dbItem.Id == id {
			if dbItem.Id != localItem.Id {
				dbItem.Id = localItem.Id
			}
			if dbItem.Header != localItem.Header {
				dbItem.Header = localItem.Header
			}
			if dbItem.Name != localItem.Name {
				dbItem.Name = localItem.Name
			}
			if dbItem.URL != localItem.URL {
				dbItem.URL = localItem.URL
			}
			if dbItem.Place != localItem.Place {
				dbItem.Place = localItem.Place
			}
			if dbItem.LogoURL != localItem.LogoURL {
				dbItem.LogoURL = localItem.LogoURL
			}
			if dbItem.ContactName != localItem.ContactName {
				dbItem.ContactName = localItem.ContactName
			}
			if dbItem.ContactURL != localItem.ContactURL {
				dbItem.ContactURL = localItem.ContactURL
			}
			if dbItem.HasWCA != localItem.HasWCA {
				dbItem.HasWCA = localItem.HasWCA
			}
			if dbItem.HasPassed != localItem.HasPassed {
				dbItem.HasPassed = localItem.HasPassed
			}
			if dbItem.Date != localItem.Date {
				dbItem.Date = localItem.Date
			}

			to.Update(*dbItem)
		}
	}

	return to
}
