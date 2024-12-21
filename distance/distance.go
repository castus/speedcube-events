package distance

import (
	"fmt"
	"math"
	"time"

	"github.com/castus/speedcube-events/logger"
)

type TravelInfo struct {
	Distance  string
	Duration  string
	Latitude  float32
	Longitude float32
}

type TravelInfoDTO struct {
	Latitude   float32
	Longitude  float32
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
		travelInfo := fetchTravelInfo(item.Place, item.DatabaseId, item.Latitude, item.Longitude)
		events[item.DatabaseId] = TravelInfo{
			Distance: travelInfo.Distance,
			Duration: travelInfo.Duration,
		}
	}

	return events
}

func fetchTravelInfo(place string, id string, latitude float32, longitude float32) TravelInfo {
	log.Info("Trying to fetch fetch travel info.", "ID", id)

	coordinate := Coordinate{}

	if latitude == 0 && longitude == 0 {
		coor, err := geoCoding(place)
		if err != nil {
			log.Error("Fail to fetch coordinates for a place", "ID", id, "Place", place, "error", err)
		}
		coordinate.Latitude = coor.Latitude
		coordinate.Longitude = coor.Longitude
	}

	coordinate.Latitude = float64(latitude)
	coordinate.Longitude = float64(longitude)

	travelInfo, err := getTravelInfo(coordinate)
	if err != nil {
		log.Error("Fail to fetch travel info", "ID", id, "lat", latitude, "lon", longitude, "error", err)
	}

	log.Info("Found fetch travel info.", "ID", id, "Distance", travelInfo.Distance, "Duration", travelInfo.Duration)

	return travelInfo
}

func geoCoding(city string) (Coordinate, error) {
	coords, err := coordinates(city)
	if err != nil {
		return Coordinate{}, err
	}

	return Coordinate{
		Latitude:  coords.Latitude,
		Longitude: coords.Longitude,
	}, nil
}

func getTravelInfo(coordinate Coordinate) (TravelInfo, error) {
	direction, err := direction(coordinate.Latitude, coordinate.Longitude)
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
