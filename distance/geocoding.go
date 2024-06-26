package distance

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"

	"github.com/tidwall/gjson"
)

type Coordinate struct {
	Longitude float64
	Latitude  float64
}

func geocodingURL(city string) string {
	encodedCity := url.QueryEscape(city)

	return fmt.Sprintf("%s/search/geocode/v6/forward?q=%s&country=pl&limit=1&proximity=ip&types=place&language=pl&access_token=%s", host, encodedCity, os.Getenv("MAPS_TOKEN"))
}

func coordinates(place string) (Coordinate, error) {
	response, err := http.Get(geocodingURL(place))
	if err != nil {
		return Coordinate{}, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return Coordinate{}, errors.New(fmt.Sprintf("Error: HTTP request status code is %d instead of 200", response.StatusCode))
	}

	bodyBytes, err := io.ReadAll(response.Body)
	if err != nil {
		return Coordinate{}, err
	}
	responseJSON := string(bodyBytes)
	features := gjson.Get(responseJSON, "features")
	if !features.Exists() {
		return Coordinate{}, nil
	}

	bestGuessCity := features.Get("0")
	cityCenter := bestGuessCity.Get("properties")
	coords := cityCenter.Get("coordinates")
	longitude := coords.Get("longitude")
	latitude := coords.Get("latitude")

	fmt.Println(longitude)
	fmt.Println(latitude)
	if !longitude.Exists() || !latitude.Exists() {
		return Coordinate{}, nil
	}

	return Coordinate{
		Latitude:  latitude.Float(),
		Longitude: longitude.Float(),
	}, nil
}
