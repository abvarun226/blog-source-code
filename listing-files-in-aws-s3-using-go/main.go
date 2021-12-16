package main

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
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

	s3Keys := make([]string, 0)

	// list files under `blog` directory in `work-with-s3` bucket.
	if err := s3client.ListObjectsPagesWithContext(ctx, &s3.ListObjectsInput{
		Bucket: aws.String(bucket),
		Prefix: aws.String("blog/"), // list files in the directory.
	}, func(o *s3.ListObjectsOutput, b bool) bool { // callback func to enable paging.
		for _, o := range o.Contents {
			s3Keys = append(s3Keys, *o.Key)
		}
		return true
	}); err != nil {
		log.Fatalf("failed to list items in s3 directory: %v", err)
	}

	log.Printf("number of files under `blog` directory: %d", len(s3Keys))
	for _, k := range s3Keys {
		log.Printf("file: %s", k)
	}
}
