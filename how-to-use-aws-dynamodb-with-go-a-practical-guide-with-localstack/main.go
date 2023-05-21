package main

import (
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

type Employee struct {
	ID   int
	Name string
	Age  int
}

func main() {
	sess, err := session.NewSession(&aws.Config{
		Region:   aws.String("us-east-1"),
		Endpoint: aws.String("http://localhost:4566"),
	})
	if err != nil {
		log.Fatalf("failed to start new session: %v", err)
	}

	svc := dynamodb.New(sess)

	deleteTable(svc)
}
