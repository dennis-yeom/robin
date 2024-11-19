package sqs

import (
	"context"
	"fmt"
	"log"

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

// check the queue to see if there are any messages
func (s *SQSClient) ReceiveMessage(ctx context.Context, visibilityTimeout int32, waitTimeSeconds int32, maxMessages int32) (bool, error) {
	// Call the ReceiveMessage API to fetch messages
	output, err := s.client.ReceiveMessage(ctx, &sqs.ReceiveMessageInput{
		QueueUrl:            &s.queueURL,
		MaxNumberOfMessages: maxMessages,
		VisibilityTimeout:   visibilityTimeout, // Time message stays hidden for others
		WaitTimeSeconds:     waitTimeSeconds,   // Long polling duration
	})
	if err != nil {
		return false, fmt.Errorf("failed to receive messages: %v", err)
	}

	// Check if any messages were received
	if len(output.Messages) == 0 {
		return false, nil // No messages received
	}

	// Process the received messages
	for _, message := range output.Messages {
		// Handle the message
		fmt.Printf("Received message: %s\n", *message.Body)

		// Delete the message after processing
		if err := s.DeleteMessage(ctx, message.ReceiptHandle); err != nil {
			log.Printf("Failed to delete message: %v", err)
		} else {
			fmt.Printf("Successfully deleted message: %s\n", *message.Body)
		}
	}

	// Return true to indicate that at least one message was processed
	return true, nil
}

// DeleteMessage deletes a message from the SQS queue
func (s *SQSClient) DeleteMessage(ctx context.Context, receiptHandle *string) error {
	_, err := s.client.DeleteMessage(ctx, &sqs.DeleteMessageInput{
		QueueUrl:      &s.queueURL,
		ReceiptHandle: receiptHandle,
	})
	if err != nil {
		return fmt.Errorf("failed to delete message: %v", err)
	}
	return nil
}
