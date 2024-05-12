package externalParser

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/castus/speedcube-events/dataFetch"
	"github.com/castus/speedcube-events/db"
	"github.com/castus/speedcube-events/logger"
	"github.com/castus/speedcube-events/s3"
)

var log = logger.Default()

func Run() {
	log.Info("Running external parsing, trying to get S3 objects")
	bucketName := os.Getenv("S3_WEB_DATA_BUCKET_NAME")
	allKeys := s3.AllKeys(bucketName)
	for _, key := range allKeys {
		items := strings.Split(key, "/")
		externalType := items[0]
		externalId := items[1]
		externalPageName := items[2]
		log.Info("Parsing object", "key", key)
		if externalType == db.CompetitionType.Cube4Fun {
			file, err := os.Open(fmt.Sprintf("exported/%s-%s", externalId, externalPageName))
			if err != nil {
				log.Error("Couldn't open file", err)
				panic(err)
			}
			defer file.Close()

			data, err := io.ReadAll(file)
			if err != nil {
				log.Error("Unable to read file: ", err)
			}
			Cube4FunParse(bytes.NewReader(data), externalType, externalId, externalPageName, dataFetch.EventsMap())
			// Cube4FunParse(s3.Contents(bucketName, key), externalType, externalId, externalPageName)
		}
	}
}
