package main

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

type SQSQueueAPI interface {
	ListQueues(ctx context.Context,
		params *sqs.ListQueuesInput,
		optFns ...func(*sqs.Options)) (*sqs.ListQueuesOutput, error)
}

func ListQueues(c context.Context, api SQSQueueAPI, input *sqs.ListQueuesInput) (*sqs.ListQueuesOutput, error) {
	return api.ListQueues(c, input)
}

// creates an sqs client.
func client(ctx context.Context, awsURL, region string) *sqs.Client {
	// customResolver is required here since we use localstack and need to point the aws url to localhost.
	customResolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
		return aws.Endpoint{
			PartitionID:   "aws",
			URL:           awsURL,
			SigningRegion: region,
		}, nil

	})

	// load the default aws config along with custom resolver.
	cfg, err := config.LoadDefaultConfig(ctx, config.WithEndpointResolverWithOptions(customResolver))
	if err != nil {
		log.Fatalf("configuration error: %v", err)
	}

	return sqs.NewFromConfig(cfg)
}

// list all queues in AWS SQS account.
func listQueues(ctx context.Context, c *sqs.Client) {
	inputList := &sqs.ListQueuesInput{}

	resultList, err := ListQueues(ctx, c, inputList)
	if err != nil {
		log.Printf("error retrieving queue URLs: %v", err)
		return
	}

	for i, url := range resultList.QueueUrls {
		log.Printf("%d: %s", i+1, url)
	}
}

func main() {
	ctx := context.TODO()

	awsURL := "http://127.0.0.1:4566"
	awsRegion := "us-west-2"

	// create aws client.
	c := client(ctx, awsURL, awsRegion)

	// list all queues in the SQS account.
	listQueues(ctx, c)
}
