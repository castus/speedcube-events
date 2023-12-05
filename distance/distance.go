package distance

import (
	"fmt"
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
