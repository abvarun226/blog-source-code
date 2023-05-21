package main

import (
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

type Employee struct {
	ID      int
	Name    string
	Age     int
	City    string
	Country string
}

func main() {
	sess, err := session.NewSession(&aws.Config{
		Region:   aws.String("us-east-1"),
		Endpoint: aws.String("http://localhost:4566"),
	})
	if err != nil {
		log.Fatalf("failed to create new session: %v", err)
	}

	svc := dynamodb.New(sess)
	input := &dynamodb.CreateTableInput{
		AttributeDefinitions: []*dynamodb.AttributeDefinition{
			{AttributeName: aws.String("ID"), AttributeType: aws.String("N")},
			{AttributeName: aws.String("Name"), AttributeType: aws.String("S")},
			{AttributeName: aws.String("City"), AttributeType: aws.String("S")},
		},
		KeySchema: []*dynamodb.KeySchemaElement{
			{AttributeName: aws.String("ID"), KeyType: aws.String("HASH")},
			{AttributeName: aws.String("Name"), KeyType: aws.String("RANGE")},
		},
		GlobalSecondaryIndexes: []*dynamodb.GlobalSecondaryIndex{
			{
				IndexName: aws.String("CityIndex"),
				KeySchema: []*dynamodb.KeySchemaElement{
					{AttributeName: aws.String("ID"), KeyType: aws.String("HASH")},
					{AttributeName: aws.String("City"), KeyType: aws.String("RANGE")},
				},
				Projection: &dynamodb.Projection{
					ProjectionType: aws.String("ALL"),
				},
				ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
					ReadCapacityUnits:  aws.Int64(5),
					WriteCapacityUnits: aws.Int64(5),
				},
			},
		},
		ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
			ReadCapacityUnits:  aws.Int64(5),
			WriteCapacityUnits: aws.Int64(5),
		},
		TableName: aws.String("Employee2"),
	}

	if _, err := svc.CreateTable(input); err != nil {
		log.Fatalf("failed to create table: %v", err)
	}

	fmt.Println("Table created successfully")

	item := Employee{
		ID:      1,
		Name:    "John Doe",
		Age:     30,
		City:    "New York",
		Country: "USA",
	}

	av, err := dynamodbattribute.MarshalMap(item)
	if err != nil {
		log.Fatalf("failed to marshal item: %v", err)
	}

	input2 := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String("Employee2"),
	}

	if _, err := svc.PutItem(input2); err != nil {
		log.Fatalf("failed to insert item: %v", err)
	}

	fmt.Println("Item inserted successfully")

	input3 := &dynamodb.QueryInput{
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":id":  {N: aws.String("1")},
			":val": {S: aws.String("New York")},
		},
		KeyConditionExpression: aws.String("ID = :id and begins_with(City, :val)"),
		TableName:              aws.String("Employee2"),
		IndexName:              aws.String("CityIndex"),
	}

	result, err := svc.Query(input3)
	if err != nil {
		log.Fatalf("failed to query table: %v", err)
	}

	for _, item := range result.Items {
		emp := Employee{}

		if err := dynamodbattribute.UnmarshalMap(item, &emp); err != nil {
			log.Fatalf("failed to unmarshal item: %v", err)
		}

		fmt.Printf("%+v\\n", emp)
	}
}
