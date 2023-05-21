package main

import (
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

func query(svc *dynamodb.DynamoDB) {
	input := &dynamodb.QueryInput{
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":val": {
				N: aws.String("1"),
			},
		},
		KeyConditionExpression: aws.String("ID = :val"),
		TableName:              aws.String("Employee"),
	}

	result, err := svc.Query(input)
	if err != nil {
		log.Fatalf("failed to query: %v", err)
		return
	}

	for _, item := range result.Items {
		emp := Employee{}

		err = dynamodbattribute.UnmarshalMap(item, &emp)
		if err != nil {
			log.Fatalf("failed to unmarshal item: %v", err)
			return
		}

		fmt.Printf("%+v\\n", emp)
	}
}
