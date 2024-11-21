package handler

import (
	"context"
	"fmt"
	"time"

	"github.com/dennis-yeom/robin/internal/aws/s3"
	"github.com/dennis-yeom/robin/internal/aws/sqs"
	"github.com/dennis-yeom/robin/internal/mongo"
	"github.com/dennis-yeom/robin/internal/redis"
	"github.com/spf13/viper"
)

type Handler struct {
	// handler is struct with 3 pointers to clients.
	// this is how we bring together all the data.
	sqs   *sqs.SQSClient
	s3    *s3.S3Client
	mongo *mongo.MongoClient
	redis *redis.RedisClient
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

// WithMongoDB sets up the MongoDB client for the handler using configuration from config.yaml
func WithMongoDB() HandlerOptions {
	return func(h *Handler) error {
		// Retrieve MongoDB configuration from Viper
		uri := viper.GetString("mongo.uri")
		dbName := viper.GetString("mongo.dbName")

		// Validate configuration
		if uri == "" || dbName == "" {
			return fmt.Errorf("MongoDB configuration is missing in config.yaml")
		}

		// Initialize the MongoDB client
		mongoClient, err := mongo.NewMongoClient(context.Background(), uri, dbName)
		if err != nil {
			return fmt.Errorf("failed to initialize MongoDB client: %w", err)
		}

		fmt.Println("MongoDB client successfully initialized")
		h.mongo = mongoClient
		return nil
	}
}

// WithRedis sets up the Redis client for the handler without testing the connection
func WithRedis() HandlerOptions {
	return func(h *Handler) error {
		// Retrieve the Redis port from Viper configuration
		port := viper.GetInt("redis.port")
		if port == 0 {
			return fmt.Errorf("redis port is not set or invalid in config.yaml")
		}

		// Initialize the Redis client
		redisClient := redis.New(port)

		// Assign the Redis client to the handler
		h.redis = redisClient

		fmt.Printf("Redis client created on port %d\n", port)
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

// ListObjectVersions prints all objects in the S3 bucket with their version IDs
func (h *Handler) ListObjectVersions() error {
	// Retrieve all objects with their version IDs
	objects, err := h.s3.GetAllObjectVersions(context.Background())
	if err != nil {
		return fmt.Errorf("failed to list object versions: %w", err)
	}

	// Print each object's key and version ID
	fmt.Println("Objects in bucket with their version IDs:")
	for _, obj := range objects {
		fmt.Printf(" - Key: %s, Version ID: %s\n", obj.Key, obj.VersionID)
	}

	return nil
}

// RedisPing tests the connection to Redis using the Ping function
func (h *Handler) RedisPing(ctx context.Context) error {
	if h.redis == nil {
		return fmt.Errorf("Redis client is not initialized")
	}

	// Call the Ping function from the Redis client
	if err := h.redis.Ping(ctx); err != nil {
		return fmt.Errorf("failed to connect to Redis: %w", err)
	}

	fmt.Println("Successfully connected to Redis!")
	return nil
}

func (h *Handler) Watch(t int) error {
	// Create a context for processing
	ctx := context.Background()

	// Set up a ticker to run every `t` seconds
	ticker := time.NewTicker(time.Duration(t) * time.Second)
	// Ensure the ticker is stopped when Watch exits
	defer ticker.Stop()

	fmt.Println("Starting periodic check on queue...")

	for range ticker.C {
		fmt.Println("Ticker triggered... checking queue for messages...")

		// Call ReceiveMessage with appropriate parameters
		success, err := h.ReceiveMessage(ctx, 30, 10, 1) // 30s visibility timeout, 10s wait time, 1 message
		if err != nil {
			fmt.Printf("Error while checking messages: %v\n", err)
			continue
		}

		if success {
			fmt.Println("A new file was added and processed from the queue!")
		} else {
			fmt.Println("No new files found in the queue.")
		}
	}
	return nil
}
