package main

import (
	"github.com/castus/speedcube-events/db"
	"github.com/castus/speedcube-events/diff"
	"github.com/castus/speedcube-events/distance"
	"github.com/castus/speedcube-events/logger"
	"github.com/castus/speedcube-events/messenger"
	"github.com/castus/speedcube-events/scrapper"
)

var log = logger.Default()

func main() {
	scrappedCompetitions := scrapper.Scrap()
	c, err := db.GetClient()
	if err != nil {
		log.Error("Couldn't get database client", err)
		panic(err)
	}

	dbCompetitions, err := db.AllItems(c)
	if err != nil {
		log.Error("Couldn't fetch items from database", err)
		panic(err)
	}

	itemsToNotify := diff.Diff(scrappedCompetitions, dbCompetitions)
	displayLogMessage(itemsToNotify)

	if itemsToNotify.IsEmpty() {
		log.Info("No changes in the events, skipping sending email.")
	} else {
		itemsToNotify.Added = addTravelInfoToItems(itemsToNotify.Added)
		itemsToNotify.Changed = addTravelInfoToItems(itemsToNotify.Changed)

		itemsToSaveToDatabase := itemsToNotify.Added
		itemsToSaveToDatabase = append(itemsToSaveToDatabase, itemsToNotify.Changed...)
		writes, err := db.AddItemsBatch(c, itemsToSaveToDatabase)
		if err != nil {
			log.Error("Couldn't save batch of items to database", err, "Saved items", writes, "All items", len(itemsToSaveToDatabase))
			panic(err)
		}
		messenger.Send(itemsToNotify)
	}

	if itemsToNotify.HasRemoved() {
		log.Info("Some items have been removed, removing from database")
		err = db.DeleteItems(c, itemsToNotify.Removed)
		if err != nil {
			log.Error("Couldn't delete item", err)
			panic(err)
		}
	}
}

func addTravelInfoToItems(competitions []db.Competition) []db.Competition {
	var newArray []db.Competition
	for _, item := range competitions {
		if item.Place == "zawody online" || item.Distance != "" || item.Duration != "" {
			newArray = append(newArray, item)
			continue
		}

		travelInfo, err := distance.Distance(item.Place)
		if err == nil {
			item.Distance = travelInfo.Distance
			item.Duration = travelInfo.Duration
		}
		newArray = append(newArray, item)
	}

	return newArray
}

func displayLogMessage(diffs diff.Differences) {
	if diffs.HasChanged() {
		log.Info("Items to change", "length", len(diffs.Changed))
	}
	if diffs.HasAdded() {
		log.Info("Items to add", "length", len(diffs.Added))
	}
	if diffs.HasRemoved() {
		log.Info("Items to remove", "length", len(diffs.Removed))
	}
}
