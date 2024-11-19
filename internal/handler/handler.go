package handler

import (
	"context"
	"fmt"

	"github.com/dennis-yeom/robin/internal/aws/s3"
	"github.com/dennis-yeom/robin/internal/aws/sqs"
)

type Handler struct {
	// handler is struct with 3 pointers to clients.
	// this is how we bring together all the data.
	sqs *sqs.SQSClient
	s3  *s3.S3Client
}

type HandlerOptions func(*Handler) error

// initialize new handler instance
func New(opts ...HandlerOptions) (*Handler, error) {
	h := &Handler{}

	// error checking for all options
	for _, opt := range opts {
		if err := opt(h); err != nil {
			return nil, err
		}
	}

	println("Successfully created handler instance!")

	return h, nil
}

// WithSQS is an option to initialize the SQS client in Demo
func WithSQS(sqsUrl string) HandlerOptions {
	return func(h *Handler) error {
		// Use NewSQSClient to initialize SQSClient with the specified queue URL
		sqsClient, err := sqs.NewSQSClient(context.Background(), sqsUrl)
		if err != nil {
			return fmt.Errorf("failed to initialize SQS client: %w", err)
		}
		fmt.Println("SQS client successfully initialized and assigned.")
		h.sqs = sqsClient
		return nil
	}
}

// WithS3 sets up the S3 client for the Demo struct
func WithS3(bucket string, endpoint string) HandlerOptions {
	return func(h *Handler) error {
		// Retrieve the endpoint from configuration
		if endpoint == "" {
			return fmt.Errorf("endpoint must be set in the config file")
		}

		// Initialize the S3 client with the bucket and endpoint
		s3Client, err := s3.NewS3Client(context.TODO(), bucket, endpoint)
		if err != nil {
			return fmt.Errorf("failed to initialize S3 client: %v", err)
		}

		h.s3 = s3Client
		return nil
	}
}

// ReceiveMessage retrieves and processes messages from SQS through the handler
func (h *Handler) ReceiveMessage(ctx context.Context, visibilityTimeout int32, waitTimeSeconds int32, maxMessages int32) (bool, error) {
	if h.sqs == nil {
		return false, fmt.Errorf("SQS client is not initialized")
	}

	// Call the SQS client's ReceiveMessage function
	success, err := h.sqs.ReceiveMessage(ctx, visibilityTimeout, waitTimeSeconds, maxMessages)

	if err != nil {
		return false, fmt.Errorf("handler failed to receive message: %w", err)
	}

	if success {
		fmt.Println("Message successfully received and processed via handler.")
	} else {
		fmt.Println("No messages received via handler.")
	}

	return success, nil
}
