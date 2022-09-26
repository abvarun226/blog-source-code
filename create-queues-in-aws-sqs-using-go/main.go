package main

import (
	"context"
	"flag"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

type SQSQueueAPI interface {
	CreateQueue(ctx context.Context,
		params *sqs.CreateQueueInput,
		optFns ...func(*sqs.Options)) (*sqs.CreateQueueOutput, error)
}

func CreateQueue(c context.Context, api SQSQueueAPI, input *sqs.CreateQueueInput) (*sqs.CreateQueueOutput, error) {
	return api.CreateQueue(c, input)
}

// creates an sqs client.
func client(awsURL, region string) *sqs.Client {
	// customResolver is required here since we use localstack and need to point the aws url to localhost.
	customResolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
		return aws.Endpoint{
			PartitionID:   "aws",
			URL:           awsURL,
			SigningRegion: region,
		}, nil

	})

	// load the default aws config along with custom resolver.
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithEndpointResolverWithOptions(customResolver))
	if err != nil {
		log.Fatalf("configuration error: %v", err)
	}

	return sqs.NewFromConfig(cfg)
}

// create a queue with the given name and attribute.
func createQueue(c *sqs.Client, queue *string, attr map[string]string) {
	input := &sqs.CreateQueueInput{
		QueueName:  queue,
		Attributes: attr,
	}

	result, err := CreateQueue(context.TODO(), c, input)
	if err != nil {
		log.Printf("error creating the queue: %v", err)
		return
	}

	log.Printf("Queue URL: %s", *result.QueueUrl)
}

func main() {
	// name of the queue as a command line option.
	queue := flag.String("q", "", "name of the queue")
	flag.Parse()

	// queue cannot be empty string.
	if *queue == "" {
		log.Println("-q argument is required. Specify a name for the queue")
		return
	}

	awsURL := "http://127.0.0.1:4566"
	awsRegion := "us-west-2"

	// create aws client.
	c := client(awsURL, awsRegion)

	// queue attributes.
	queueAttributes := map[string]string{
		"DelaySeconds":           "60",
		"MessageRetentionPeriod": "86400",
	}

	// create a queue with the given name.
	createQueue(c, queue, queueAttributes)
}
