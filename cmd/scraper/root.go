package scraper

import (
	"fmt"
	"github.com/castus/speedcube-events/dataFetch"
	"github.com/castus/speedcube-events/db"
	"github.com/castus/speedcube-events/diff"
	"github.com/castus/speedcube-events/distance"
	"github.com/castus/speedcube-events/externalFetcher"
	"github.com/castus/speedcube-events/logger"
	"github.com/castus/speedcube-events/messenger"
	"github.com/spf13/cobra"
	"os"
)

var log = logger.Default()
var useMock bool
var printK8SConfig bool

func Setup() *cobra.Command {
	cmd.Flags().BoolVarP(&useMock, "mock", "m", false, "Use mock file to scrap data")
	cmd.Flags().BoolVarP(&printK8SConfig, "k8s", "k", false, "Print Kubernetes config")

	return cmd
}

var cmd = &cobra.Command{
	Use:   "scrape",
	Short: "Scrape single source of truth, parse it and add all necessary information",
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

		localItemsIds := localItemsDatabase.GetIds()

		// Reset databases for safety, to not use it later
		database = db.Database{}
		localItemsDatabase = db.Database{}

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

		if diffIDs.IsEmpty() {
			log.Info("No changes in the events, skipping sending email.")
		} else {
			message := messenger.PrepareMessage(
				makeMessengerDTO(diffIDs.Added, mergedDatabase),
				makeMessengerDTO(diffIDs.Passed, mergedDatabase),
				makeMessengerDTO(diffIDs.Changed, mergedDatabase))
			messenger.Send(message)
		}

		onlyScrapedItems := db.CompetitionsCollection{}
		for _, item := range localItemsIds {
			dbItem := mergedDatabase.Get(item)
			onlyScrapedItems = append(onlyScrapedItems, dbItem)
		}

		onlyScrapEligible := mergedDatabase.FilterScrapCube4FunEligible()
		k8SCube4FunDTO := makeK8SCube4FunDTO(onlyScrapEligible)
		k8SPPODTO := makeK8SPPODTO(onlyScrapEligible)

		if printK8SConfig {
			fmt.Println(externalFetcher.PrintK8sJobsConfig(k8SCube4FunDTO, k8SPPODTO))
			return
		}

		if os.Getenv("KUBERNETES_SERVICE_HOST") != "" && os.Getenv("KUBERNETES_SERVICE_PORT") != "" {
			externalFetcher.SpinK8sJobsToFetchExternalData(k8SCube4FunDTO, k8SPPODTO)
			log.Info("Running k8s job.")
		} else {
			log.Info("Detected local environment, skipping spinning K8s jobs.")
		}
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

func makeK8SCube4FunDTO(competitions db.CompetitionsCollection) []externalFetcher.K8SConfigCube4FunDTO {
	var items []externalFetcher.K8SConfigCube4FunDTO
	for _, competition := range competitions {
		if competition.Type != db.CompetitionType.Cube4Fun {
			panic(fmt.Sprintf("Expected Cube4Fun item, got %s", competition.Type))
		}
		items = append(items, externalFetcher.K8SConfigCube4FunDTO{
			Type: competition.Type,
			Id:   competition.Id,
			URL:  competition.URL,
		})
	}

	return items
}

func makeK8SPPODTO(competitions db.CompetitionsCollection) []externalFetcher.K8SConfigPPODTO {
	var items []externalFetcher.K8SConfigPPODTO
	for _, competition := range competitions {
		if competition.Type != db.CompetitionType.PPO {
			panic(fmt.Sprintf("Expected PPO item, got %s", competition.Type))
		}
		items = append(items, externalFetcher.K8SConfigPPODTO{
			Type: competition.Type,
			Id:   competition.Id,
			URL:  competition.URL,
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

func makeMessengerDTO(IDs []string, db *db.Database) []messenger.MessengerDTO {
	var items []messenger.MessengerDTO
	for _, id := range IDs {
		item := db.Get(id)
		items = append(items, messenger.MessengerDTO{
			LogoURL:         item.LogoURL,
			Name:            item.Name,
			URL:             item.URL,
			HasWCA:          item.HasWCA,
			Date:            item.Date,
			Distance:        item.Distance,
			Duration:        item.Duration,
			Place:           item.Place,
			Events:          item.Events,
			MainEvent:       item.MainEvent,
			CompetitorLimit: item.CompetitorLimit,
			Registered:      item.Registered,
			ContactURL:      item.ContactURL,
			ContactName:     item.ContactName,
		})
	}

	return items
}
