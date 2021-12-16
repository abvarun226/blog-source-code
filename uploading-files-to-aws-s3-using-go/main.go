package main

import (
	"log"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

func main() {
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

	// upload couple of files to s3 bucket we just created.
	if err := uploadFile(sess, bucket, "blog/file1.txt", "this is file 1"); err != nil {
		log.Fatalf("failed to upload file: %v", err)
	}
	if err := uploadFile(sess, bucket, "blog/file2.txt", "this is file 2"); err != nil {
		log.Fatalf("failed to upload file: %v", err)
	}
}

func uploadFile(sess *session.Session, bucket, s3path, content string) error {
	uploader := s3manager.NewUploader(sess)

	_, err := uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(s3path),
		Body:   strings.NewReader(content),
	})

	return err
}
