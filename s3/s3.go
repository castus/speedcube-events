package s3

import (
	"context"
	"io"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/castus/speedcube-events/logger"
)

var log = logger.Default()

func s3Client() (*s3.Client, error) {
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

func AllKeys(bucketName string) []string {
	var keys = []string{}
	c, err := s3Client()
	if err != nil {
		log.Error("Couldn't get S3Client client", "error", err)
		panic(err)
	}

	out, err := c.ListObjectsV2(context.TODO(), &s3.ListObjectsV2Input{
		Bucket: aws.String(bucketName),
	})
	if err != nil {
		log.Error("Couldn't get contents of bucket.", "error", err)
		return keys
	}

	for _, item := range out.Contents {
		keys = append(keys, *item.Key)
	}

	return keys
}

func Contents(bucketName string, key string) string {
	c, err := s3Client()
	log.Info("Trying to get S3 object.", "bucketName", bucketName, "objectKey", key)
	if err != nil {
		log.Error("Couldn't get S3Client client", "error", err)
		panic(err)
	}

	out, err := c.GetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(key),
	})
	if err != nil {
		log.Error("Couldn't get object from bucket.", "bucketName", bucketName, "objectKey", key, "error", err)
	}
	s, err := io.ReadAll(out.Body)
	if err != nil {
		panic(err)
	}

	return string(s)
}

func Save(bucketName string, key string, content io.Reader) {
	c, err := s3Client()
	if err != nil {
		log.Error("Couldn't get S3Client client", "error", err)
		panic(err)
	}

	log.Info("Saving to S3 bucket", "name", bucketName)
	_, err = c.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(key),
		Body:   content,
	})
	if err != nil {
		log.Error("Couldn't save file to bucket.", "bucketName", bucketName, "objectKey", key, "error", err)
	}
}
