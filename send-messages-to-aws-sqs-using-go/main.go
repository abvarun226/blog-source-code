package main

import (
	"context"
	"flag"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
)

type SQSQueueAPI interface {
	GetQueueUrl(ctx context.Context,
		params *sqs.GetQueueUrlInput,
		optFns ...func(*sqs.Options)) (*sqs.GetQueueUrlOutput, error)

	SendMessage(ctx context.Context,
		params *sqs.SendMessageInput,
		optFns ...func(*sqs.Options)) (*sqs.SendMessageOutput, error)
}

func GetQueueURL(c context.Context, api SQSQueueAPI, input *sqs.GetQueueUrlInput) (*sqs.GetQueueUrlOutput, error) {
	return api.GetQueueUrl(c, input)
}

func SendMessage(c context.Context, api SQSQueueAPI, input *sqs.SendMessageInput) (*sqs.SendMessageOutput, error) {
	return api.SendMessage(c, input)
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

// send a message to a queue.
func sendMessage(ctx context.Context, c *sqs.Client, queue *string) {
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

	// Send a message with attributes to the given queue
	messageInput := &sqs.SendMessageInput{
		DelaySeconds: 10,
		MessageAttributes: map[string]types.MessageAttributeValue{
			"Blog": {
				DataType:    aws.String("String"),
				StringValue: aws.String("The Code Library"),
			},
			"Article": {
				DataType:    aws.String("Number"),
				StringValue: aws.String("10"),
			},
		},
		MessageBody: aws.String("article about sending a message to AWS SQS using Go"),
		QueueUrl:    queueURL,
	}

	resp, err := SendMessage(ctx, c, messageInput)
	if err != nil {
		log.Printf("error sending the message: %v", err)
		return
	}

	log.Printf("Message ID: %s", *resp.MessageId)
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

	// send a message to the given queue
	sendMessage(ctx, c, queue)
}
