package handler

import (
	"context"
	"fmt"

	"github.com/dennis-yeom/robin/internal/aws/sqs"
)

type Handler struct {
	// handler is struct with 3 pointers to clients.
	// this is how we bring together all the data.
	sqs *sqs.SQSClient
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
