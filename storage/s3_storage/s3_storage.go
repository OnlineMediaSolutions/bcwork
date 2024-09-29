package s3storage

import (
	"io"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/rotisserie/eris"
)

type S3 interface {
	ListS3Objects(bucket string, prefix string) (*s3.ListObjectsV2Output, error)
	GetObjectInput(bucket string, key string) ([]byte, error)
}

type S3Storage struct {
	client *s3.S3
}

var _ S3 = (*S3Storage)(nil)

func New() (*S3Storage, error) {
	const region = "us-east-1"

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(region)},
	)
	if err != nil {
		return nil, eris.Wrapf(err, "failed to create aws session")
	}

	return &S3Storage{client: s3.New(sess)}, nil
}

// ListS3Objects Get the list of items
func (s *S3Storage) ListS3Objects(bucket string, prefix string) (*s3.ListObjectsV2Output, error) {
	resp, err := s.client.ListObjectsV2(&s3.ListObjectsV2Input{
		Bucket: aws.String(bucket),
		Prefix: aws.String(prefix),
	})
	if err != nil {
		return nil, eris.Wrapf(err, "unable to list items in bucket %q, %v", bucket, err)
	}

	return resp, nil

}

// GetObjectInput retrieve the object
func (s *S3Storage) GetObjectInput(bucket string, key string) ([]byte, error) {
	// Create the request to get the object
	input := &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	}
	// Call S3 to retrieve the object
	result, err := s.client.GetObject(input)
	if err != nil {
		return nil, eris.Wrapf(err, "error retrieving object from S3")
	}

	// Read the body of the S3 object into memory
	buffer, err := io.ReadAll(result.Body)
	if err != nil {
		return nil, eris.Wrapf(err, "failed to read object")
	}
	defer result.Body.Close()

	// Do something with the object content, `buffer` contains the file data
	// For example, you could print the size of the file:
	return buffer, nil
}
