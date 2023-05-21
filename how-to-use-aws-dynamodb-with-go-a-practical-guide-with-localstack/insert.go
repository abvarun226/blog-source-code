package main

import (
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

func insert(svc *dynamodb.DynamoDB) {
	item := Employee{
		ID:   1,
		Name: "John Doe",
		Age:  30,
	}

	av, err := dynamodbattribute.MarshalMap(item)
	if err != nil {
		log.Fatalf("failed to marshal item: %v", err)
	}

	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String("Employee"),
	}

	if _, err := svc.PutItem(input); err != nil {
		log.Fatalf("failed to put item: %v", err)
	}

	fmt.Println("Item inserted successfully")
}
