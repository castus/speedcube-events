package main

import (
	"fmt"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/castus/speedcube-events/dataFetch"
	"github.com/castus/speedcube-events/db"
	"github.com/castus/speedcube-events/diff"
	"github.com/castus/speedcube-events/distance"
	"github.com/castus/speedcube-events/logger"
	"github.com/castus/speedcube-events/messenger"
)

var log = logger.Default()

func main() {
	scrappedCompetitions := dataFetch.ScrapCompetitions()
	//printer.PrettyPrint(scrappedCompetitions)

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

	diffIDs := diff.Diff(scrappedCompetitions, dbCompetitions)
	displayLogMessage(diffIDs)

	fullDataCompetitions := updateDatabase(scrappedCompetitions, dbCompetitions, c, diffIDs.Removed)

	if diffIDs.IsEmpty() {
		log.Info("No changes in the events, skipping sending email.")
	} else {
		message := ""
		message = fmt.Sprintf("%s\n%s\n", message, messenger.PrepareMessageForAdded(diffIDs, fullDataCompetitions))
		message = fmt.Sprintf("%s\n%s\n", message, messenger.PrepareMessageForChanged(diffIDs, fullDataCompetitions))
		message = fmt.Sprintf("%s\n%s\n", message, messenger.PrepareMessageForRemoved(diffIDs, dbCompetitions))
		messenger.Send(message)
	}
}

func updateDatabase(scrappedCompetitions db.Competitions, dbCompetitions db.Competitions, client *dynamodb.Client, itemsToRemove []string) db.Competitions {
	log.Info("Trying to update database.")
	scrappedCompetitions = dataFetch.IncludeEvents(scrappedCompetitions)
	scrappedCompetitions = dataFetch.IncludeRegistrations(scrappedCompetitions)
	scrappedCompetitions = dataFetch.IncludeGeneralInfo(scrappedCompetitions)
	scrappedCompetitions = distance.IncludeTravelInfo(scrappedCompetitions, dbCompetitions)
	writes, err := db.AddItemsBatch(client, scrappedCompetitions)
	if err != nil {
		log.Error("Couldn't save batch of items to database", "error", err, "savedItems", writes, "allItems", len(scrappedCompetitions))
		panic(err)
	}
	log.Info("Saved batch of items to database", "savedItems", writes, "allItems", len(scrappedCompetitions))

	if len(itemsToRemove) > 0 {
		var competitionsToRemove db.Competitions
		for _, id := range itemsToRemove {
			dbItem := dbCompetitions.FindByID(id)
			if dbItem != nil {
				competitionsToRemove = append(competitionsToRemove, *dbItem)
			}
		}
		log.Info("Some items have been removed, removing from database")
		err = db.DeleteItems(client, competitionsToRemove)
		if err != nil {
			log.Error("Couldn't delete item", "error", err)
			panic(err)
		}
	}

	return scrappedCompetitions
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
