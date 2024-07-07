package externalParser

import (
	"bytes"
	"os"
	"strings"

	"github.com/castus/speedcube-events/dataFetch"
	"github.com/castus/speedcube-events/db"
	"github.com/castus/speedcube-events/exporter"
	"github.com/castus/speedcube-events/logger"
	"github.com/castus/speedcube-events/s3"
)

var log = logger.Default()

func Run() {
	database := db.Database{}
	database.Initialize()

	log.Info("Trying to get S3 objects for external parsing")
	bucketName := os.Getenv("S3_WEB_DATA_BUCKET_NAME")
	allKeys := s3.AllKeys(bucketName)
	eventsMap := dataFetch.InitializeEventsMap()
	for _, key := range allKeys {
		items := strings.Split(key, "/")
		externalType := items[0]
		id := items[1]
		externalPageName := items[2]
		log.Info("Trying to parse object", "key", key)
		dbItem := database.Get(id)
		if dbItem.HasPassed {
			log.Info("Event passed, nothing to parse", "key", key)
			continue
		}

		s3Content := s3.Contents(bucketName, key)
		if externalType == db.CompetitionType.Cube4Fun {
			Cube4FunParse(bytes.NewReader([]byte(s3Content)), externalType, id, externalPageName, eventsMap, dbItem)
		} else if externalType == db.CompetitionType.PPO {
			PPOParse(bytes.NewReader([]byte(s3Content)), externalType, id, externalPageName, eventsMap, dbItem)
		}
	}
	exporter.Export(database)
}
