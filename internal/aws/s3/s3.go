package s3

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
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

// GetObjectVersion retrieves the metadata of an object and returns its version ID
func (s *S3Client) GetObjectVersion(ctx context.Context, key string) (string, error) {
	input := &s3.HeadObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
	}

	result, err := s.client.HeadObject(ctx, input)
	if err != nil {
		return "", fmt.Errorf("failed to get object metadata: %w", err)
	}

	versionID := aws.ToString(result.VersionId)
	//fmt.Printf("File %s in bucket %s has version ID: %s\n", key, s.bucket, versionID)
	return versionID, nil
}

// ObjectInfo holds the key (filename) and version ID of an object
type ObjectInfo struct {
	Key       string
	VersionID string
}

// GetAllObjectVersions retrieves the filename and version ID for all objects in the S3 bucket
func (s *S3Client) GetAllObjectVersions(ctx context.Context) ([]ObjectInfo, error) {
	input := &s3.ListObjectsV2Input{
		Bucket: aws.String(s.bucket),
	}

	result, err := s.client.ListObjectsV2(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to list objects: %w", err)
	}

	var objects []ObjectInfo
	for _, item := range result.Contents {
		// Attempt to retrieve the version ID for each object
		versionID, err := s.GetObjectVersion(ctx, *item.Key)
		if err != nil {
			fmt.Printf("Failed to get version for object %s: %v\n", *item.Key, err)
			continue
		}

		// Append object with key and (optional) version ID
		objects = append(objects, ObjectInfo{
			Key:       *item.Key,
			VersionID: versionID, // May be empty if versioning is not enabled
		})
	}

	return objects, nil
}
