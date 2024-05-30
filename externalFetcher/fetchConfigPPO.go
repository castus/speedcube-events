package externalFetcher

import (
	"fmt"

	"github.com/castus/speedcube-events/db"
)

func FetchConfigPPO(competitions db.Competitions) []ExternalFetchConfig {
	var items = []ExternalFetchConfig{}
	for _, competition := range competitions {
		if competition.Type == db.CompetitionType.PPO {
			items = append(items, ppoConfig(db.PageTypes.Info, competition))
			items = append(items, ppoConfig(db.PageTypes.Competitors, competition))
		}
	}

	return items
}

func ppoConfig(page string, competition db.Competition) ExternalFetchConfig {
	return ExternalFetchConfig{
		Type:         competition.Type,
		Id:           competition.Id,
		URL:          competition.URL,
		S3BucketPath: fmt.Sprintf("%s/%s/%s.html", competition.Type, competition.Id, page),
		PageType:     page,
	}
}
