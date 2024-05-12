package externalFetcher

import (
	"fmt"

	"github.com/castus/speedcube-events/db"
)

func FetchConfigCube4Fun(competitions db.Competitions) []ExternalFetchConfig {
	var cube4FunItems = []ExternalFetchConfig{}
	for _, competition := range competitions {
		if competition.Type == db.CompetitionType.Cube4Fun {
			cube4FunItems = append(cube4FunItems, cube4FunConfig(db.Cube4FunPages.Info, competition))
			cube4FunItems = append(cube4FunItems, cube4FunConfig(db.Cube4FunPages.Competitors, competition))
		}
	}

	return cube4FunItems
}

func cube4FunConfig(page string, competition db.Competition) ExternalFetchConfig {
	return ExternalFetchConfig{
		Type:         competition.Type,
		Id:           competition.TypeSpecificId,
		URL:          fmt.Sprintf("%s/%s", competition.URL, page),
		S3BucketPath: fmt.Sprintf("%s/%s/%s.html", competition.Type, competition.TypeSpecificId, page),
	}
}
