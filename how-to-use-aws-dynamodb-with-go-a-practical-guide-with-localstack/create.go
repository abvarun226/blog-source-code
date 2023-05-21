package main

import (
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

func createTable(svc *dynamodb.DynamoDB) {
	input := &dynamodb.CreateTableInput{
		AttributeDefinitions: []*dynamodb.AttributeDefinition{
			{AttributeName: aws.String("ID"), AttributeType: aws.String("N")},
		},
		KeySchema: []*dynamodb.KeySchemaElement{
			{AttributeName: aws.String("ID"), KeyType: aws.String("HASH")},
		},
		ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
			ReadCapacityUnits:  aws.Int64(5),
			WriteCapacityUnits: aws.Int64(5),
		},
		TableName: aws.String("Employee"),
	}

	if _, err := svc.CreateTable(input); err != nil {
		log.Fatalf("failed to create table: %v", err)
	}

	fmt.Println("Table created successfully")
}
