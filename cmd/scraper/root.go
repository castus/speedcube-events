package scraper

import (
	"github.com/castus/speedcube-events/distance"
	"github.com/castus/speedcube-events/printer"

	"github.com/castus/speedcube-events/dataFetch"
	"github.com/castus/speedcube-events/db"
	"github.com/castus/speedcube-events/diff"
	"github.com/castus/speedcube-events/logger"
	"github.com/spf13/cobra"
)

var log = logger.Default()
var useMock bool
var printK8SConfig bool

func Setup() *cobra.Command {
	Cmd.Flags().BoolVarP(&useMock, "mock", "m", false, "Use mock file to scrap data")
	Cmd.Flags().BoolVarP(&printK8SConfig, "k8s", "k", false, "Print Kubernetes config")

	return Cmd
}

var Cmd = &cobra.Command{
	Use:   "scrape",
	Short: "Scrape single source of truth, parse it and export to DynamoDB",
	Run: func(cmd *cobra.Command, args []string) {
		database := db.Database{}
		database.Initialize()
		var fetcher dataFetch.DataFetcher = dataFetch.WebFetcher{}

		if useMock {
			fetcher = dataFetch.FileFetcher{}
		}

		scrapedCompetitions := dataFetch.ScrapCompetitions(fetcher)
		if len(scrapedCompetitions) == 0 {
			log.Info("No scraped competitions, finishing.")
			return
		}
		localItemsDatabase := db.InitializeWith(scrapedCompetitions)

		diffIDs := diff.Diff(&localItemsDatabase, &database)
		diffIDs.PrintDifferencesInfo()

		merger := db.NewMerger(diffIDs.Added, diffIDs.Passed, diffIDs.Changed)
		mergedDatabase := merger.Merge(localItemsDatabase, database)

		onlyWCAEvents := mergedDatabase.FilterWCAApiEligible()
		wcaAPIData := dataFetch.GetWCAApiData(makeWCAApiDTO(onlyWCAEvents))
		for _, event := range onlyWCAEvents {
			dbItem := mergedDatabase.Get(event.Id)
			dbItem.Events = wcaAPIData[event.Id].Events
			dbItem.MainEvent = wcaAPIData[event.Id].MainEvent
			dbItem.CompetitorLimit = wcaAPIData[event.Id].CompetitorLimit
			dbItem.Registered = wcaAPIData[event.Id].Registered
			mergedDatabase.Update(*dbItem)
		}

		onlyTravelEligible := mergedDatabase.FilterTravelInfoEligible()
		travelData := distance.GetTravelData(makeTravelInfoDTO(onlyTravelEligible))
		for _, event := range onlyTravelEligible {
			dbItem := mergedDatabase.Get(event.Id)
			dbItem.Distance = travelData[event.Id].Distance
			dbItem.Duration = travelData[event.Id].Duration
			mergedDatabase.Update(*dbItem)
		}

		for _, event := range mergedDatabase.GetAll() {
			printer.PrettyPrint(event)
		}
		return

		// if diffIDs.IsEmpty() {
		// 	log.Info("No changes in the events, skipping sending email.")
		// } else {
		// 	message := messenger.PrepareHeader()
		// 	message = fmt.Sprintf("%s\n%s\n", message, messenger.PrepareMessageForAdded(diffIDs, fullDataCompetitions))
		// 	message = fmt.Sprintf("%s\n%s\n", message, messenger.PrepareMessageForChanged(diffIDs, fullDataCompetitions))
		// 	message = fmt.Sprintf("%s\n%s\n", message, messenger.PrepareMessageForRemoved(diffIDs, dbCompetitions))
		// 	message = fmt.Sprintf("%s\n%s\n", message, messenger.PrepareFooter())

		// 	messenger.Send(message)
		// }

		// if printK8SConfig {
		// 	fmt.Println(externalFetcher.GetK8sJobsConfig(fullDataCompetitions))
		// 	return
		// }

		// if os.Getenv("KUBERNETES_SERVICE_HOST") != "" && os.Getenv("KUBERNETES_SERVICE_PORT") != "" {
		// 	externalFetcher.SpinK8sJobsToFetchExternalData(fullDataCompetitions)
		// 	log.Info("Running k8s job.")
		// } else {
		// 	log.Info("Detected local environment, skipping spinning K8s jobs to fetch external resources.")
		// }
	},
}

func makeWCAApiDTO(competitions db.CompetitionsCollection) []dataFetch.WCAApiDTO {
	var items []dataFetch.WCAApiDTO
	for _, competition := range competitions {
		items = append(items, dataFetch.WCAApiDTO{
			DatabaseId: competition.Id,
			OtherId:    competition.ExtractWCAId(),
		})
	}

	return items
}

func makeTravelInfoDTO(competitions db.CompetitionsCollection) []distance.TravelInfoDTO {
	var items []distance.TravelInfoDTO
	for _, competition := range competitions {
		items = append(items, distance.TravelInfoDTO{
			DatabaseId: competition.Id,
			Place:      competition.Place,
		})
	}

	return items
}

// func updateDatabase(scrappedCompetitions db.Competitions, dbCompetitions db.Competitions, client *dynamodb.Client, itemsToRemove []string) db.Competitions {
// 	log.Info("Trying to update database.")
// 	scrappedCompetitions = dataFetch.IncludeEvents(scrappedCompetitions)
// 	scrappedCompetitions = dataFetch.IncludeRegistrations(scrappedCompetitions, dataFetch.WebFetcher{})
// 	scrappedCompetitions = dataFetch.IncludeGeneralInfo(scrappedCompetitions, dataFetch.WebFetcher{})
// 	scrappedCompetitions = distance.IncludeTravelInfo(scrappedCompetitions, dbCompetitions)
// 	// printer.PrettyPrint(scrappedCompetitions)
// 	// os.Exit(1)
// 	writes, err := db.AddItemsBatch(client, scrappedCompetitions)
// 	if err != nil {
// 		log.Error("Couldn't save batch of items to database", "error", err, "savedItems", writes, "allItems", len(scrappedCompetitions))
// 		panic(err)
// 	}
// 	log.Info("Saved batch of items to database", "savedItems", writes, "allItems", len(scrappedCompetitions))

// 	if len(itemsToRemove) > 0 {
// 		var competitionsToRemove db.Competitions
// 		for _, id := range itemsToRemove {
// 			dbItem := dbCompetitions.FindByID(id)
// 			if dbItem != nil {
// 				dbItem.HasPassed = true
// 				competitionsToRemove = append(competitionsToRemove, *dbItem)
// 			}
// 		}
// 		writes, err := db.AddItemsBatch(client, competitionsToRemoveasd)
// 		if err != nil {
// 			log.Error("Couldn't save batch of items to database", "error", err, "savedItems", writes, "allItems", len(scrappedCompetitions))
// 			panic(err)
// 		}
// 		log.Info("Some items have been removed, marking them as passed events", "savedItems", writes, "allItems", len(scrappedCompetitions))
// 	}3

// 	return scrappedCompetitions
// }
