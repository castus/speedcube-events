package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/castus/speedcube-events/dataFetch"
	"github.com/castus/speedcube-events/db"
	"github.com/castus/speedcube-events/diff"
	"github.com/castus/speedcube-events/distance"
	"github.com/castus/speedcube-events/exporter"
	"github.com/castus/speedcube-events/logger"
	"github.com/castus/speedcube-events/messenger"
	"github.com/castus/speedcube-events/printer"
)

var log = logger.Default()

func main() {
	args := os.Args

	if len(args) > 1 && strings.Contains(args[1], "saveWebpage") {
		exporter.SaveWebpageAsFile("kalendarz-imprez.html")
		log.Info("Webpage saved on disk")
		return
	}

	var fetcher dataFetch.DataFetcher = dataFetch.WebFetcher{}

	if len(args) > 1 && strings.Contains(args[1], "mock") {
		fetcher = dataFetch.FileFetcher{}
	}

	scrappedCompetitions := dataFetch.ScrapCompetitions(fetcher)
	if len(scrappedCompetitions) == 0 {
		log.Info("No scraped competitions, finishing")
		return
	}
	// printer.PrettyPrint(scrappedCompetitions)

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

	if len(args) > 1 && strings.Contains(args[1], "exportDatabase") {
		exporter.Export(dbCompetitions)
		return
	}

	diffIDs := diff.Diff(scrappedCompetitions, dbCompetitions)
	displayLogMessage(diffIDs)
	fullDataCompetitions := updateDatabase(scrappedCompetitions, dbCompetitions, c, diffIDs.Removed)

	if diffIDs.IsEmpty() {
		log.Info("No changes in the events, skipping sending email.")
	} else {
		message := messenger.PrepareHeader()
		message = fmt.Sprintf("%s\n%s\n", message, messenger.PrepareMessageForAdded(diffIDs, fullDataCompetitions))
		message = fmt.Sprintf("%s\n%s\n", message, messenger.PrepareMessageForChanged(diffIDs, fullDataCompetitions))
		message = fmt.Sprintf("%s\n%s\n", message, messenger.PrepareMessageForRemoved(diffIDs, dbCompetitions))
		message = fmt.Sprintf("%s\n%s\n", message, messenger.PrepareFooter())

		messenger.Send(message)
	}
}

func updateDatabase(scrappedCompetitions db.Competitions, dbCompetitions db.Competitions, client *dynamodb.Client, itemsToRemove []string) db.Competitions {
	log.Info("Trying to update database.")
	scrappedCompetitions = dataFetch.IncludeEvents(scrappedCompetitions)
	scrappedCompetitions = dataFetch.IncludeRegistrations(scrappedCompetitions, dataFetch.WebFetcher{})
	scrappedCompetitions = dataFetch.IncludeGeneralInfo(scrappedCompetitions, dataFetch.WebFetcher{})
	scrappedCompetitions = distance.IncludeTravelInfo(scrappedCompetitions, dbCompetitions)
	printer.PrettyPrint(scrappedCompetitions)
	os.Exit(1)
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
				dbItem.HasPassed = true
				competitionsToRemove = append(competitionsToRemove, *dbItem)
			}
		}
		writes, err := db.AddItemsBatch(client, competitionsToRemove)
		if err != nil {
			log.Error("Couldn't save batch of items to database", "error", err, "savedItems", writes, "allItems", len(scrappedCompetitions))
			panic(err)
		}
		log.Info("Some items have been removed, marking them as passed events", "savedItems", writes, "allItems", len(scrappedCompetitions))
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
		log.Info("Items to mark as passed", "length", len(diffs.Removed))
	}
}
