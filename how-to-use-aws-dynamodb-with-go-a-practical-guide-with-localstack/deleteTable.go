package main

import (
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

func deleteTable(svc *dynamodb.DynamoDB) {
	input := &dynamodb.DeleteTableInput{
		TableName: aws.String("Employee2"),
	}

	if _, err := svc.DeleteTable(input); err != nil {
		log.Fatalf("failed to delete item: %v", err)
	}

	fmt.Println("Item deleted successfully")
}
