package exporter

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/castus/speedcube-events/dataFetch"
	"github.com/castus/speedcube-events/db"
	"github.com/castus/speedcube-events/logger"
	"github.com/castus/speedcube-events/s3"
)

var log = logger.Default()

const (
	fileName = "data.json"
)

func ExportForFrontend(database db.Database) {
	items := database.GetAll()
	exportToFile(items, fileName)
	exportToStorage(fileName)
	cleanup(fileName)
}

func ExportLocal(database db.Database) {
	items := database.GetAll()
	exportToFile(items, fileName)
}

func PersistDatabase(database db.Database) {
	database.StoreInDynamoDB()
}

func SaveWebpageAsFile(name string) {
	URL := fmt.Sprintf("%s/%s", dataFetch.Host, dataFetch.EventsPath)
	r, ok := dataFetch.WebFetcher{}.Fetch(URL)
	if !ok {
		log.Error("Couldn't fetch webpage to save it on disk")
		return
	}

	file, err := os.Create(name)
	if err != nil {
		log.Error("Couldn't create webpage file", "error", err)
		panic(err)
	}

	defer file.Close()

	content, _ := io.ReadAll(r)
	_, err = file.Write(content)
	if err != nil {
		log.Error("Couldn't write to a webpage file", "error", err)
		panic(err)
	}

	log.Info("Webpage file created.")
}

func exportToFile(items db.CompetitionsCollection, fileName string) {
	j, _ := json.MarshalIndent(items, "", "    ")
	file, err := os.Create(fileName)
	if err != nil {
		log.Error("Couldn't create database file", "error", err)
		panic(err)
	}

	defer file.Close()
	_, _ = file.Write(j)

	log.Info("Database file created.")
}

func exportToStorage(fileName string) {
	file, err := os.Open(fileName)
	if err != nil {
		log.Error("Couldn't open file to upload.", "file", fileName, "error", err)
	} else {
		defer file.Close()
		bucketName := os.Getenv("S3_BUCKET_NAME")
		s3.Save(bucketName, fileName, file)
	}
}

func cleanup(fileName string) {
	err := os.Remove(fileName)
	if err != nil {
		log.Error("Couldn't cleanup database file.", "file", fileName, "error", err)
	}
}
