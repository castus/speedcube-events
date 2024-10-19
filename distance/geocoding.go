package distance

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
)

type Coordinate struct {
	Longitude float64 `json:"longitude"`
	Latitude  float64 `json:"latitude"`
}

type Features struct {
	Features []Feature `json:"features"`
}

type Feature struct {
	Properties Property `json:"properties"`
}

type Property struct {
	FeatureType   string     `json:"feature_type"`
	Name          string     `json:"name"`
	NamePreferred string     `json:"name_preferred"`
	Coordinates   Coordinate `json:"coordinates"`
}

func geocodingURL(city string) string {
	encodedCity := url.QueryEscape(city)

	return fmt.Sprintf("%s/search/geocode/v6/forward?q=%s&country=pl&proximity=ip&types=place&language=pl&access_token=%s", host, encodedCity, os.Getenv("MAPS_TOKEN"))
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

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return Coordinate{}, err
	}
	var data Features
	err = json.Unmarshal(body, &data)
	if err != nil {
		log.Error("Couldn't unmarshal JSON", "error", err)
		return Coordinate{}, err
	}

	for _, feature := range data.Features {
		if feature.Properties.Name == place {
			return Coordinate{
				Longitude: feature.Properties.Coordinates.Longitude,
				Latitude:  feature.Properties.Coordinates.Latitude,
			}, nil
		}
	}

	log.Debug("Couldn't find place to get coordinates", "place", place)
	return Coordinate{}, nil
}
