package sqs

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

// struct to hold sqs client info
type SQSClient struct {
	client   *sqs.Client
	queueURL string
}

// NewSQSClient initializes a new SQS client for the specified queue URL
func NewSQSClient(ctx context.Context, queueURL string) (*SQSClient, error) {
	// Load the configuration with the AWS profile
	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithRegion("us-east-1"),
		config.WithSharedConfigProfile("aws"), // Specify the AWS profile
	)
	if err != nil {
		return nil, fmt.Errorf("unable to load SDK config: %v", err)
	}

	// Print a message indicating successful SQS client creation
	fmt.Println("Successfully created SQS client and connected to queue:", queueURL)

	return &SQSClient{
		client:   sqs.NewFromConfig(cfg),
		queueURL: queueURL,
	}, nil
}
