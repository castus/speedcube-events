package distance

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"

	"github.com/tidwall/gjson"
)

type DirectionMeasures struct {
	Distance float64
	Duration float64
}

func direction(lat float64, long float64) (DirectionMeasures, error) {
	response, err := http.Get(directionsURL(lat, long))
	if err != nil {
		return DirectionMeasures{}, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return DirectionMeasures{}, fmt.Errorf("error: HTTP request status code is %d instead of 200", response.StatusCode)
	}

	bodyBytes, err := io.ReadAll(response.Body)
	if err != nil {
		return DirectionMeasures{}, err
	}
	responseJSON := string(bodyBytes)

	routes := gjson.Get(responseJSON, "routes")
	if !routes.Exists() {
		return DirectionMeasures{}, nil
	}

	dist := routes.Get("0.distance")
	duration := routes.Get("0.duration")

	if !dist.Exists() || !duration.Exists() {
		return DirectionMeasures{}, nil
	}

	return DirectionMeasures{
		Distance: dist.Float(),
		Duration: duration.Float(),
	}, nil
}

func directionsURL(lat float64, long float64) string {
	coords := fmt.Sprintf("%s,%s;%f,%f", os.Getenv("HOME_LON"), os.Getenv("HOME_LAT"), long, lat)
	decodedCoordinates := url.QueryEscape(coords)

	return fmt.Sprintf("%s/directions/v5/mapbox/driving/%s?alternatives=false&geometries=geojson&language=en&overview=full&steps=true&alley_bias=-1&radiuses=unlimited;unlimited&access_token=%s", host, decodedCoordinates, os.Getenv("MAPS_TOKEN"))
}
