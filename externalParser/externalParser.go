package externalParser

import (
	"bytes"
	"os"
	"strings"

	"github.com/castus/speedcube-events/dataFetch"
	"github.com/castus/speedcube-events/db"
	"github.com/castus/speedcube-events/logger"
	"github.com/castus/speedcube-events/s3"
)

var log = logger.Default()

func Run() {
	c, err := db.GetClient()
	if err != nil {
		log.Error("Couldn't get database client", err)
		panic(err)
	}

	dbCompetitions, err := db.AllItems(c)
	if err != nil {
		log.Error("Couldn't fetch items from database", err)
		panic(err)
	}

	log.Info("Running external parsing, trying to get S3 objects")
	bucketName := os.Getenv("S3_WEB_DATA_BUCKET_NAME")
	allKeys := s3.AllKeys(bucketName)
	eventsMap := dataFetch.EventsMap()
	for _, key := range allKeys {
		items := strings.Split(key, "/")
		externalType := items[0]
		id := items[1]
		externalPageName := items[2]
		log.Info("Trying to parse object", "key", key)
		if externalType == db.CompetitionType.Cube4Fun {
			dbItem := dbCompetitions.FindByID(id)
			if dbItem.HasPassed {
				continue
			}

			s3Content := s3.Contents(bucketName, key)
			Cube4FunParse(bytes.NewReader([]byte(s3Content)), externalType, id, externalPageName, eventsMap, dbItem, c)
		}
	}
}