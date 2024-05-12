package externalParser

import (
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/castus/speedcube-events/exporter"
	"github.com/castus/speedcube-events/logger"
)

var log = logger.Default()

func Run() {
	getAllObjects()
	log.Info("Run")
}

func getAllObjects() {
	c, err := exporter.S3Client()
	if err != nil {
		log.Error("Couldn't get database client", err)
		panic(err)
	}

	bucketName := os.Getenv("S3_WEB_DATA_BUCKET_NAME")
	out, err := c.ListObjectsV2(context.TODO(), &s3.ListObjectsV2Input{
		Bucket: aws.String(bucketName),
	})

	for _, item := range out.Contents {
		fmt.Println(*item.Key)
	}
	// out, err = c.PutObject(context.TODO(), &s3.PutObjectInput{
	// 	Bucket: aws.String(bucketName),
	// 	Key:    aws.String(fileName),
	// 	Body:   file,
	// })
	if err != nil {
		log.Error("Couldn't get contents of bucket.", "error", err)
	}
	log.Info(bucketName)
}
