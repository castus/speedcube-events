package db

import (
	"context"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

const (
	tableName = "SpeedcubeEvents"
)

func GetClient() (*dynamodb.Client, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion("eu-central-1"),
		config.WithCredentialsProvider(credentials.StaticCredentialsProvider{
			Value: aws.Credentials{
				AccessKeyID:     os.Getenv("AWS_API_KEY"),
				SecretAccessKey: os.Getenv("AWS_API_SECRET"),
				SessionToken:    "",
				Source:          "Speedcube Events app",
			},
		}),
	)
	if err != nil {
		return nil, err
	}

	c := dynamodb.NewFromConfig(cfg)
	return c, nil
}

func PutItem(c *dynamodb.Client, competition Competition) (err error) {
	item, err := attributevalue.MarshalMap(competition)
	if err != nil {
		panic(err)
	}

	_, err = c.PutItem(context.TODO(), &dynamodb.PutItemInput{
		TableName: aws.String(tableName),
		Item:      item,
	})
	if err != nil {
		log.Printf("Couldn't add item to table. Here's why: %v\n", err)
	}
	return err
}

func DeleteItem(c *dynamodb.Client, competition Competition) error {
	id, err := attributevalue.Marshal(competition.Id)
	if err != nil {
		panic(err)
	}
	name, err := attributevalue.Marshal(competition.Name)
	if err != nil {
		panic(err)
	}
	_, err = c.DeleteItem(context.TODO(), &dynamodb.DeleteItemInput{
		TableName: aws.String(tableName),
		Key:       map[string]types.AttributeValue{"Id": id, "Name": name},
	})
	if err != nil {
		log.Printf("Couldn't delete %v from the table. Here's why: %v\n", competition.Id, err)
	}
	return err
}

func AllItems(c *dynamodb.Client) ([]Competition, error) {
	var competitions []Competition
	var err error
	var response *dynamodb.ScanOutput
	response, err = c.Scan(context.TODO(), &dynamodb.ScanInput{
		TableName: aws.String(tableName),
	})
	if err != nil {
		log.Println(err)
	} else {
		err = attributevalue.UnmarshalListOfMaps(response.Items, &competitions)
		if err != nil {
			log.Printf("Couldn't unmarshal query response. Here's why: %v\n", err)
		}
	}
	return competitions, err
}
