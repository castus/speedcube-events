package scraper

import (
	"fmt"
	"os"

	"github.com/castus/speedcube-events/dataFetch"
	"github.com/castus/speedcube-events/db"
	"github.com/castus/speedcube-events/diff"
	"github.com/castus/speedcube-events/distance"
	"github.com/castus/speedcube-events/exporter"
	"github.com/castus/speedcube-events/externalFetcher"
	"github.com/castus/speedcube-events/logger"
	"github.com/castus/speedcube-events/messenger"
	"github.com/spf13/cobra"
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
			updateWCAItem(event.Id, mergedDatabase, wcaAPIData[event.Id], event.WCAId, true)
		}

		// C4F events are WCA compliant, but their Id is different
		// so we need to find the id using WCA search
		// and fill the rest of the details to not use browser scrap
		onlyC4FEvents := mergedDatabase.FilterScrapCube4FunEligible()
		for _, event := range onlyC4FEvents {
			WCAId, err := dataFetch.SearchWCAApiData(makeSearchWCAApiDTO(event))
			if err != nil {
				log.Error("Error when finding Cube4Fun event on WCA", "error", err)
				continue
			}

			dto := getWCAApiDTO(event.Id, WCAId)
			out := []dataFetch.WCAApiDTO{dto}
			wcaAPIData = dataFetch.GetWCAApiData(out)

			// C4F is weird, because the registered number of users is not available on WCA
			// We have to scrap their website for that number
			updateWCAItem(event.Id, mergedDatabase, wcaAPIData[event.Id], WCAId, false)
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

		onlyScrapedItems := []db.Competition{}
		for _, item := range localItemsIds {
			dbItem := mergedDatabase.Get(item)
			onlyScrapedItems = append(onlyScrapedItems, *dbItem)
		}

		onlyScrapDatabase := db.InitializeWith(onlyScrapedItems)
		k8SCube4FunDTO := makeK8SCube4FunDTO(onlyScrapDatabase.FilterScrapCube4FunEligible())
		k8SPPODTO := makeK8SPPODTO(onlyScrapDatabase.FilterScrapPPOEligible())

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

		exporter.ExportForFrontend(*mergedDatabase)
		exporter.PersistDatabase(*mergedDatabase)
	},
}

func updateWCAItem(id string, database *db.Database, wcaAPIData dataFetch.WCAApiResponse, wcaID string, includeRegistered bool) {
	dbItem := database.Get(id)
	dbItem.WCAId = wcaID
	dbItem.Events = wcaAPIData.Events
	dbItem.MainEvent = wcaAPIData.MainEvent
	dbItem.CompetitorLimit = wcaAPIData.CompetitorLimit
	if includeRegistered {
		dbItem.Registered = wcaAPIData.Registered
	}
	dbItem.Longitude = wcaAPIData.Longitude
	dbItem.Latitude = wcaAPIData.Latitude
	database.Update(*dbItem)
}

func makeWCAApiDTO(competitions db.CompetitionsCollection) []dataFetch.WCAApiDTO {
	var items []dataFetch.WCAApiDTO
	for _, competition := range competitions {
		items = append(items, getWCAApiDTO(competition.Id, competition.ExtractWCAId()))
	}

	return items
}

func getWCAApiDTO(databaseId string, otherId string) dataFetch.WCAApiDTO {
	return dataFetch.WCAApiDTO{
		DatabaseId: databaseId,
		OtherId:    otherId,
	}
}

func makeSearchWCAApiDTO(competition *db.Competition) dataFetch.SearchWCAApiDTO {
	item := dataFetch.SearchWCAApiDTO{
		Name: competition.Name,
	}

	return item
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
			Latitude:   competition.Latitude,
			Longitude:  competition.Longitude,
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
