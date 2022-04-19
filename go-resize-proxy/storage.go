package main

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

// Storage is an interface for handling any upload of images.
type Storage interface {
	Upload(key string, b []byte, contentType string)
}

//---------------

// S3Storage will be implementing the interface Storage
type S3Storage struct {
	S3Client *s3.S3
}

// Setup a connection and get a new instance of S3Storage
func NewS3Storage() *S3Storage {
	const accessKey = "minioadmin" //your_access_key
	const secretKey = "minioadmin" //your_secret_key

	// Configure to use MinIO Server
	s3Config := &aws.Config{
		Credentials:      credentials.NewStaticCredentials(accessKey, secretKey, ""),
		Endpoint:         aws.String(S3_ADDRESS),
		Region:           aws.String("us-east-1"),
		DisableSSL:       aws.Bool(true),
		S3ForcePathStyle: aws.Bool(true),
	}
	newSession := session.New(s3Config)

	s3Storage := &S3Storage{}
	s3Storage.S3Client = s3.New(newSession)

	return s3Storage
}

// Implementing Storage.Upload
func (s *S3Storage) Upload(s3key string, b []byte, contentType string) {

	fmt.Println("****uploading to s3", s3key)

	bucket := aws.String("images")
	key := aws.String(strings.TrimPrefix(s3key, "/images/")) // since we set bucketname separately we need only the s3 key.

	_, err := s.S3Client.PutObject(&s3.PutObjectInput{
		Body:        bytes.NewReader(b),
		Bucket:      bucket,
		Key:         key,
		ContentType: aws.String(contentType)})
	if err != nil {
		fmt.Printf("Failed to upload data to %s/%s, %s\n", *bucket, *key, err.Error())
		return
	}
	fmt.Printf("Successfully uploaded data with key %s into bucket %s\n", *key, *bucket)
}
