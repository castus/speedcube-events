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

type Competitions []Competition

func AddItemBatch(c *dynamodb.Client, item Competition) (int, error) {
	panic("Don't use it")
}

type Database struct {
	items     map[string]Competition
	client    *dynamodb.Client
	tableName string
}

func (d *Database) Initialize() {
	d.tableName = os.Getenv("TABLE_NAME")
	c, err := d.getClient()
	if err != nil {
		log.Error("Couldn't get database client", "error", err)
		panic(err)
	}
	d.client = c
	d.items = make(map[string]Competition)

	competitions, err := d.fetchAllItems()
	if err != nil {
		log.Error("Couldn't fetch items from database", "error", err)
		panic(err)
	}

	for _, item := range competitions {
		d.items[item.Id] = item
	}
}

func InitializeWith(competitions []Competition) Database {
	d := Database{}
	for _, item := range competitions {
		d.Add(item)
	}

	return d
}

func (d *Database) Add(item Competition) {
	if len(d.items) == 0 {
		d.items = make(map[string]Competition)
	}
	_, thereIsAnItem := d.items[item.Id]
	if thereIsAnItem {
		msg := "You try to add item that's already in the database"
		log.Error(msg)
		panic(msg)
	}

	d.items[item.Id] = item
}

func (d *Database) Update(item Competition) {
	_, thereIsAnItem := d.items[item.Id]
	if !thereIsAnItem {
		msg := "You try to update item that's not in the database"
		log.Error(msg)
		panic(msg)
	}

	d.items[item.Id] = item
}

func (d *Database) Get(id string) *Competition {
	item, ok := d.items[id]
	if !ok {
		return nil
	}

	return &item
}

func (d *Database) GetAll() CompetitionsCollection {
	var items = CompetitionsCollection{}
	for _, v := range d.items {
		items = append(items, &v)
	}

	return items
}

func (d *Database) GetIds() []string {
	var items = []string{}
	for _, v := range d.items {
		items = append(items, v.Id)
	}

	return items
}

func (d *Database) FilterWCAApiEligible() CompetitionsCollection {
	var items = d.GetAll()
	items = items.FilterWCAEvents()

	return items
}

func (d *Database) FilterScrapCube4FunEligible() CompetitionsCollection {
	var items = d.GetAll()
	out := items.FilterCube4Fun()

	return out.FilterHasURL()
}

func (d *Database) FilterScrapPPOEligible() CompetitionsCollection {
	var items = d.GetAll()
	out := items.FilterPPO()

	return out.FilterHasURL()
}

func (d *Database) FilterTravelInfoEligible() CompetitionsCollection {
	var items = d.GetAll()
	items = items.FilterNotPassed()
	items = items.FilterNotOnline()
	items = items.FilterEmptyDistanceOrDuration()

	return items
}

func (d *Database) StoreInDynamoDB() {
	var items = []Competition{}
	for _, v := range d.items {
		items = append(items, v)
	}

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
				log.Error("Couldn't marshal competition", competition.Name, "Here's why: ", "error", err)
			} else {
				writeReqs = append(
					writeReqs,
					types.WriteRequest{PutRequest: &types.PutRequest{Item: item}},
				)
			}
		}
		_, err = d.client.BatchWriteItem(context.TODO(), &dynamodb.BatchWriteItemInput{
			RequestItems: map[string][]types.WriteRequest{d.tableName: writeReqs}})
		if err != nil {
			log.Error("Couldn't add a batch of items to table", d.tableName, "Here's why", "error", err)
		} else {
			written += len(writeReqs)
		}
		start = end
		end += batchSize
	}

	if err != nil {
		log.Error("Couldn't save batch of items to database", "error", err, "savedItems", written, "allItems", len(items))
		panic(err)
	}
	log.Info("Saved batch of items to database", "savedItems", written, "allItems", len(items))
}

func (d *Database) getClient() (*dynamodb.Client, error) {
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

func (d *Database) fetchAllItems() ([]Competition, error) {
	var competitions []Competition
	var err error
	response, err := d.client.Scan(context.TODO(), &dynamodb.ScanInput{
		TableName: aws.String(d.tableName),
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
