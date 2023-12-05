package main

import (
	"github.com/castus/speedcube-events/db"
	"github.com/castus/speedcube-events/diff"
	"github.com/castus/speedcube-events/messanger"
	"github.com/castus/speedcube-events/scrapper"
)

func main() {
	scrappedCompetitions := scrapper.Scrap()
	c, _ := db.GetClient()
	dbCompetitions, err := db.AllItems(c)
	if err != nil {
		panic(err)
	}

	itemsToNotify := diff.Diff(scrappedCompetitions, dbCompetitions)
	if !itemsToNotify.IsEmpty() {
		messanger.Send(itemsToNotify)
	}

	if itemsToNotify.HasRemoved() {
		for _, item := range itemsToNotify.Removed {
			err := db.DeleteItem(c, item)
			if err != nil {
				panic(err)
			}
		}
	}

	// loop over scrapper and add to database
}
