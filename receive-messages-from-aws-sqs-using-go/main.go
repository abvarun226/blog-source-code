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

	ReceiveMessage(ctx context.Context,
		params *sqs.ReceiveMessageInput,
		optFns ...func(*sqs.Options)) (*sqs.ReceiveMessageOutput, error)
}

func GetQueueURL(c context.Context, api SQSQueueAPI, input *sqs.GetQueueUrlInput) (*sqs.GetQueueUrlOutput, error) {
	return api.GetQueueUrl(c, input)
}

func ReceiveMessage(c context.Context, api SQSQueueAPI, input *sqs.ReceiveMessageInput) (*sqs.ReceiveMessageOutput, error) {
	return api.ReceiveMessage(c, input)
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

// receive a message from a queue.
func recvMessage(ctx context.Context, c *sqs.Client, queue *string) {
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

	// Receive a message with attributes to the given queue
	recvInput := &sqs.ReceiveMessageInput{
		QueueUrl:              queueURL,
		MessageAttributeNames: []string{"All"},
		MaxNumberOfMessages:   1,
		VisibilityTimeout:     int32(10),
	}

	msg, err := ReceiveMessage(ctx, c, recvInput)
	if err != nil {
		log.Printf("error receiving messages: %v", err)
		return
	}

	if msg.Messages == nil {
		log.Printf("No messages found")
		return
	}

	log.Printf("Message ID: %s, Message Body: %s", *msg.Messages[0].MessageId, *msg.Messages[0].Body)
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

	// create aws client
	c := client(ctx, awsURL, awsRegion)

	// receive a message from the given queue
	recvMessage(ctx, c, queue)
}
