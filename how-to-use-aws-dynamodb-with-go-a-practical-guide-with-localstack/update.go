package main

import (
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

func update(svc *dynamodb.DynamoDB) {
	input := &dynamodb.UpdateItemInput{
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":n": {
				S: aws.String("Jane Doe"),
			},
		},
		TableName: aws.String("Employee"),
		Key: map[string]*dynamodb.AttributeValue{
			"ID": {
				N: aws.String("1"),
			},
		},
		UpdateExpression: aws.String("set #name = :n"),
		ExpressionAttributeNames: map[string]*string{
			"#name": aws.String("Name"),
		},
	}

	if _, err := svc.UpdateItem(input); err != nil {
		log.Fatalf("failed to update item: %v", err)
	}

	fmt.Println("Item updated successfully")
}
