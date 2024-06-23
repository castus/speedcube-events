package distance

import (
	"fmt"
	"math"
	"time"

	"github.com/castus/speedcube-events/logger"
)

type TravelInfo struct {
	Distance string
	Duration string
}

type TravelInfoDTO struct {
	Place      string
	DatabaseId string
}

const (
	host = "https://api.mapbox.com"
)

var log = logger.Default()

func GetTravelData(ids []TravelInfoDTO) map[string]TravelInfo {
	var events = make(map[string]TravelInfo)
	for _, item := range ids {
		if item.Place != "Bydgoszcz" {
			continue
		}
		travelInfo := fetchTravelInfo(item.Place, item.DatabaseId)
		events[item.DatabaseId] = TravelInfo{
			Distance: travelInfo.Distance,
			Duration: travelInfo.Duration,
		}
	}

	return events
}

func fetchTravelInfo(place string, id string) TravelInfo {
	log.Info("Trying to fetch fetch travel info.", "ID", id)

	travelInfo, err := distance(place)
	if err != nil {
		log.Error("Fail to fetch travel info", "ID", id, "error", err)
	}

	log.Info("Found fetch travel info.", "ID", id, "Distance", travelInfo.Distance, "Duration", travelInfo.Duration)

	return travelInfo
}

//func IncludeTravelInfo(competitions db.Competitions, databaseItems db.Competitions) db.Competitions {
//	var newArray db.Competitions
//	for _, item := range competitions {
//		if item.Place == "zawody online" {
//			log.Info("No need to fetch travel info, detected ONLINE event", "eventID", item.Id)
//			newArray = append(newArray, item)
//			continue
//		}
//		if item.HasPassed {
//			log.Info("No need to fetch travel info, detected PASSED event", "eventID", item.Id)
//			newArray = append(newArray, item)
//			continue
//		}
//
//		dbItem := databaseItems.FindByID(item.Id)
//		if dbItem != nil && (dbItem.Distance != "" || dbItem.Duration != "" || dbItem.HasPassed) {
//			log.Info("No need to fetch travel info, event has already have it or has been passed", "eventID", item.Id)
//			item.Distance = dbItem.Distance
//			item.Duration = dbItem.Duration
//			newArray = append(newArray, item)
//
//			continue
//		}
//
//		log.Info("Trying to fetch travel info", "eventID", item.Id)
//		travelInfo, err := distance(item.Place)
//		if err == nil {
//			log.Info("Success to fetch travel info", "eventID", item.Id, "distance", travelInfo.Distance, "duration", travelInfo.Duration)
//			item.Distance = travelInfo.Distance
//			item.Duration = travelInfo.Duration
//		} else {
//			log.Error("Fail to fetch travel info", "eventID", item.Id, "error", err)
//		}
//		newArray = append(newArray, item)
//	}
//
//	return newArray
//}

func distance(city string) (TravelInfo, error) {
	coords, err := coordinates(city)
	if err != nil {
		return TravelInfo{}, err
	}

	direction, err := Direction(coords.Latitude, coords.Longitude)
	if err != nil {
		return TravelInfo{}, err
	}

	return TravelInfo{
		Distance: distanceMessage(direction.Distance),
		Duration: durationMessage(direction.Duration),
	}, nil
}

func distanceMessage(distance float64) string {
	rounded := math.Round(distance) / 1000 // to KM

	return fmt.Sprintf("%.0fkm", rounded)
}

func durationMessage(duration float64) string {
	seconds := math.Round(duration)
	dur := time.Duration(seconds) * time.Second
	hours := int(dur.Hours())
	minutes := int(dur.Minutes()) % 60

	return fmt.Sprintf("%d godzin, %d minut", hours, minutes)
}
