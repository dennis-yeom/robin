package s3

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type S3Client struct {
	client *s3.Client
	bucket string
}

// NewS3Client initializes a new S3 client for the specified bucket and endpoint
func NewS3Client(ctx context.Context, bucket, endpoint string) (*S3Client, error) {
	// Load the configuration with the Linode profile
	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithRegion("us-east-1"),
		config.WithSharedConfigProfile("linode"), // Specify the Linode profile
	)
	if err != nil {
		return nil, fmt.Errorf("unable to load SDK config: %w", err)
	}

	// Configure the S3 client with a custom endpoint for Linode
	s3Client := s3.New(s3.Options{
		Region:           "us-east-1",
		EndpointResolver: s3.EndpointResolverFromURL(endpoint),
		Credentials:      cfg.Credentials,
		UsePathStyle:     true, // Enable path-style addressing
	})

	fmt.Println("s3 client initialized")

	return &S3Client{
		client: s3Client,
		bucket: bucket,
	}, nil
}
