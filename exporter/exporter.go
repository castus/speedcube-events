package exporter

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/castus/speedcube-events/dataFetch"
	"github.com/castus/speedcube-events/db"
	"github.com/castus/speedcube-events/logger"
	"io"
	"os"
)

var log = logger.Default()

const (
	fileName = "data.json"
)

func Export(items db.Competitions) {
	exportToFile(items, fileName)
	exportToStorage(fileName)
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
		log.Error("Couldn't create webpage file", err)
		panic(err)
	}

	defer file.Close()

	content, err := io.ReadAll(r)
	_, err = file.Write(content)
	if err != nil {
		log.Error("Couldn't write to a webpage file", err)
		panic(err)
	}

	log.Info("Webpage file created.")
}

func exportToFile(items db.Competitions, fileName string) {
	j, _ := json.MarshalIndent(items, "", "    ")
	file, err := os.Create(fileName)
	if err != nil {
		log.Error("Couldn't create database file", err)
		panic(err)
	}

	defer file.Close()
	_, err = file.Write(j)

	log.Info("Database file created.")
}

func exportToStorage(fileName string) {
	c, err := storageClient()
	if err != nil {
		log.Error("Couldn't get database client", err)
		panic(err)
	}

	file, err := os.Open(fileName)
	if err != nil {
		log.Error("Couldn't open file to upload.", "file", fileName, "error", err)
	} else {
		defer file.Close()
		log.Info("Saving to S3 bucket", "name", os.Getenv("S3_BUCKET_NAME"))
		bucketName := os.Getenv("S3_BUCKET_NAME")
		_, err = c.PutObject(context.TODO(), &s3.PutObjectInput{
			Bucket: aws.String(bucketName),
			Key:    aws.String(fileName),
			Body:   file,
		})
		if err != nil {
			log.Error("Couldn't save file to bucket.", "file", fileName, "bucketName", bucketName, "objectKey", fileName, "error", err)
		}
	}
}

func storageClient() (*s3.Client, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion("eu-central-1"),
		config.WithCredentialsProvider(credentials.StaticCredentialsProvider{
			Value: aws.Credentials{
				AccessKeyID:     os.Getenv("AWS_S3_API_KEY"),
				SecretAccessKey: os.Getenv("AWS_S3_API_SECRET"),
				SessionToken:    "",
				Source:          "Speedcube Events app",
			},
		}),
	)
	if err != nil {
		return nil, err
	}

	c := s3.NewFromConfig(cfg)

	return c, nil
}
