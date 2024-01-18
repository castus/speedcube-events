package distance

import (
	"fmt"
	"github.com/castus/speedcube-events/db"
	"github.com/castus/speedcube-events/logger"
	"math"
	"time"
)

type TravelInfo struct {
	Distance string
	Duration string
}

const (
	host = "https://api.mapbox.com"
)

var log = logger.Default()

func IncludeTravelInfo(competitions db.Competitions, databaseItems db.Competitions) db.Competitions {
	var newArray db.Competitions
	for _, item := range competitions {
		dbItem := databaseItems.FindByID(item.Id)
		if item.Place == "zawody online" || (dbItem != nil && (dbItem.Distance != "" || dbItem.Duration != "")) {
			log.Info("No need to fetch travel info", "eventID", item.Id)
			newArray = append(newArray, *dbItem)
			continue
		}

		log.Info("Trying to fetch travel info", "eventID", item.Id)
		travelInfo, err := Distance(item.Place)
		if err == nil {
			log.Info("Success to fetch travel info", "eventID", item.Id, "distance", travelInfo.Distance, "duration", travelInfo.Duration)
			item.Distance = travelInfo.Distance
			item.Duration = travelInfo.Duration
		} else {
			log.Error("Fail to fetch travel info", "eventID", item.Id, "error", err)
		}
		newArray = append(newArray, item)
	}

	return newArray
}

func Distance(city string) (TravelInfo, error) {
	coordinates, err := Coordinates(city)
	if err != nil {
		return TravelInfo{}, err
	}

	direction, err := Direction(coordinates.Latitude, coordinates.Longitude)
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
