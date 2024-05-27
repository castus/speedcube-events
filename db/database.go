package db

import (
	"context"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

func tableName() string {
	return os.Getenv("TABLE_NAME")
}

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
		TableName: aws.String(tableName()),
		Item:      item,
	})
	if err != nil {
		log.Error("Couldn't add item to table. Here's why", err)
	}
	return err
}

func AddItemBatch(c *dynamodb.Client, item Competition) (int, error) {
	return AddItemsBatch(c, []Competition{item})
}

func AddItemsBatch(c *dynamodb.Client, items []Competition) (int, error) {
	var err error
	var item map[string]types.AttributeValue
	maxItems := 100
	written := 0
	batchSize := 25 // DynamoDB allows a maximum batch size of 25 items.
	start := 0
	end := start + batchSize
	for start < maxItems && start < len(items) {
		var writeReqs []types.WriteRequest
		if end > len(items) {
			end = len(items)
		}
		for _, competition := range items[start:end] {
			item, err = attributevalue.MarshalMap(competition)
			if err != nil {
				log.Error("Couldn't marshal competition", competition.Name, "Here's why: ", err)
			} else {
				writeReqs = append(
					writeReqs,
					types.WriteRequest{PutRequest: &types.PutRequest{Item: item}},
				)
			}
		}
		_, err = c.BatchWriteItem(context.TODO(), &dynamodb.BatchWriteItemInput{
			RequestItems: map[string][]types.WriteRequest{tableName(): writeReqs}})
		if err != nil {
			log.Error("Couldn't add a batch of items to table", tableName(), "Here's why", err)
		} else {
			written += len(writeReqs)
		}
		start = end
		end += batchSize
	}

	return written, err
}

func DeleteItem(c *dynamodb.Client, ID string, name string) error {
	id, err := attributevalue.Marshal(ID)
	if err != nil {
		panic(err)
	}
	nameValue, err := attributevalue.Marshal(name)
	if err != nil {
		panic(err)
	}
	_, err = c.DeleteItem(context.TODO(), &dynamodb.DeleteItemInput{
		TableName: aws.String(tableName()),
		Key:       map[string]types.AttributeValue{"Id": id, "Name": nameValue},
	})
	if err != nil {
		log.Error("Couldn't delete item from the table.", ID, err)
	}
	return err
}

func DeleteItems(c *dynamodb.Client, items Competitions) error {
	for _, item := range items {
		err := DeleteItem(c, item.Id, item.Name)
		if err != nil {
			return err
		}
	}

	return nil
}

func AllItems(c *dynamodb.Client) (Competitions, error) {
	var competitions Competitions
	var err error
	var response *dynamodb.ScanOutput
	response, err = c.Scan(context.TODO(), &dynamodb.ScanInput{
		TableName: aws.String(tableName()),
	})
	if err != nil {
		log.Error("Could fetch all items", "error", err)
	} else {
		err = attributevalue.UnmarshalListOfMaps(response.Items, &competitions)
		if err != nil {
			log.Error("Couldn't unmarshal query response. Here's why:", "error", err)
		}
	}
	return competitions, err
}
