package utils

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/rotisserie/eris"
	"io"
)

func ListS3Objects(bucket string, prefix string) (*s3.ListObjectsV2Output, error) {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-east-1")},
	)
	if err != nil {
		return nil, eris.Wrapf(err, "failed to create aws session")
	}

	// Create S3 service client
	svc := s3.New(sess)

	// Get the list of items
	resp, err := svc.ListObjectsV2(&s3.ListObjectsV2Input{
		Bucket: aws.String(bucket),
		Prefix: aws.String(prefix),
	})
	if err != nil {
		return nil, eris.Wrapf(err, "Unable to list items in bucket %q, %v", bucket, err)
	}

	return resp, nil

}

func GetObjectInput(bucket string, key string) (string, error) {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-east-1")},
	)
	if err != nil {
		return "", eris.Wrapf(err, "failed to create aws session")
	}

	// Create S3 service client
	svc := s3.New(sess)

	// Create the request to get the object
	input := &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	}
	// Call S3 to retrieve the object
	result, err := svc.GetObject(input)
	if err != nil {
		fmt.Println("Error retrieving object from S3: ", err)
		return "", eris.Wrapf(err, "failed to fetch object")
	}

	// Read the body of the S3 object into memory
	buffer, err := io.ReadAll(result.Body)
	if err != nil {
		return "", eris.Wrapf(err, "failed to read object")
	}

	// Close the body of the object to prevent memory leaks
	result.Body.Close()

	// Do something with the object content, `buffer` contains the file data
	// For example, you could print the size of the file:
	return string(buffer), nil

}
