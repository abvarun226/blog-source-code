package main

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

func main() {
	ctx := context.Background()
	bucket := "work-with-s3"

	conf := aws.NewConfig().
		WithRegion("us-west-2").
		WithEndpoint("http://127.0.0.1:4566").
		WithS3ForcePathStyle(true)

	sess, err := session.NewSession(conf)
	if err != nil {
		log.Fatalf("failed to create a new aws session: %v", sess)
	}

	s3client := s3.New(sess)

	// create a new s3 bucket.
	if _, err := s3client.CreateBucket(&s3.CreateBucketInput{Bucket: aws.String(bucket)}); err != nil {
		if aerr, ok := err.(awserr.Error); ok && aerr.Code() != "Conflict" && aerr.Code() != "BucketAlreadyOwnedByYou" {
			log.Fatalf("failed to create a new s3 bucket: %v", err)
		}
	}

	// iterator to delete all files under `blog` directory.
	iter := s3manager.NewDeleteListIterator(s3client, &s3.ListObjectsInput{
		Bucket: aws.String(bucket),
		Prefix: aws.String("blog/"),
	})

	// use the iterator to delete the files.
	if err := s3manager.NewBatchDeleteWithClient(s3client).Delete(ctx, iter); err != nil {
		log.Fatalf("failed to delete files under given directory: %v", err)
	}
}
