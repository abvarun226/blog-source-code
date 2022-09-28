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
	GetQueueUrl(ctx context.Context,
		params *sqs.GetQueueUrlInput,
		optFns ...func(*sqs.Options)) (*sqs.GetQueueUrlOutput, error)

	DeleteQueue(ctx context.Context,
		params *sqs.DeleteQueueInput,
		optFns ...func(*sqs.Options)) (*sqs.DeleteQueueOutput, error)
}

func GetQueueURL(c context.Context, api SQSQueueAPI, input *sqs.GetQueueUrlInput) (*sqs.GetQueueUrlOutput, error) {
	return api.GetQueueUrl(c, input)
}

func DeleteQueue(c context.Context, api SQSQueueAPI, input *sqs.DeleteQueueInput) (*sqs.DeleteQueueOutput, error) {
	return api.DeleteQueue(c, input)
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

// delete a queue with the given name.
func deleteQueue(ctx context.Context, c *sqs.Client, queue *string) {
	// Get the URL for the queue
	input := &sqs.GetQueueUrlInput{
		QueueName: queue,
	}
	resultGet, err := GetQueueURL(ctx, c, input)
	if err != nil {
		log.Printf("error getting the queue URL: %v", err)
		return
	}
	queueURL := resultGet.QueueUrl

	// delete the queue using the queue URL
	dqInput := &sqs.DeleteQueueInput{
		QueueUrl: queueURL,
	}
	if _, err := DeleteQueue(ctx, c, dqInput); err != nil {
		log.Printf("error deleting the queue: %v", err)
		return
	}

	log.Printf("deleted queue with URL: %s", *queueURL)
}

func main() {
	ctx := context.TODO()

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
	c := client(ctx, awsURL, awsRegion)

	// delete a queue with the given name.
	deleteQueue(ctx, c, queue)
}
