package main

import (
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

func delete(svc *dynamodb.DynamoDB) {
	input := &dynamodb.DeleteItemInput{
		TableName: aws.String("Employee"),
		Key: map[string]*dynamodb.AttributeValue{
			"ID": {
				N: aws.String("1"),
			},
		},
	}

	if _, err := svc.DeleteItem(input); err != nil {
		log.Fatalf("failed to delete item: %v", err)
	}

	fmt.Println("Item deleted successfully")
}
