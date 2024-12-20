package cmd

import (
	"context"
	"fmt"
	"log"

	"github.com/dennis-yeom/robin/internal/handler"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	t int

	// root command
	RootCmd = &cobra.Command{
		Use:   "robin",
		Short: "runs robin",
		Long:  "main command for robin",
		Run: func(cmd *cobra.Command, arg []string) {
			fmt.Println("running robin...\n for options: go run main.go --help")
		},
	}

	// handler command
	HandlerCmd = &cobra.Command{
		Use:   "handler",
		Short: "creates instance of handler",
		RunE: func(cmd *cobra.Command, args []string) error {
			_, err := handler.New()
			if err != nil {
				return err
			}
			return nil
		},
	}

	SQSClientCmd = &cobra.Command{
		Use:   "sqs",
		Short: "instantiates sqs client",
		RunE: func(cmd *cobra.Command, args []string) error {
			_, err := handler.New(
				handler.WithSQS(viper.GetString("sqs.url")),
			)

			if err != nil {
				return err
			}

			return nil

		},
	}

	// s3 command
	S3Cmd = &cobra.Command{
		Use:   "s3",
		Short: "instantiates an S3 client",
		RunE: func(cmd *cobra.Command, args []string) error {
			// Get the bucket and endpoint from the configuration file
			bucket := viper.GetString("s3.bucket")
			endpoint := viper.GetString("s3.endpoint")

			// Create a handler with an S3 client
			_, err := handler.New(
				handler.WithS3(bucket, endpoint),
			)
			if err != nil {
				return fmt.Errorf("failed to initialize S3 client: %w", err)
			}

			//fmt.Println("S3 client successfully created!")
			return nil
		},
	}

	// MongoDB command
	MongoCmd = &cobra.Command{
		Use:   "mongo",
		Short: "initializes the Mongo client",
		RunE: func(cmd *cobra.Command, args []string) error {
			// Create the handler with MongoDB client
			_, err := handler.New(
				handler.WithMongoDB(),
			)
			if err != nil {
				return fmt.Errorf("failed to initialize MongoDB client: %w", err)
			}

			fmt.Println("MongoDB client successfully initialized!")
			return nil
		},
	}

	// Redis command
	RedisCmd = &cobra.Command{
		Use:   "redis",
		Short: "initializes the Redis client and checks connection",
		RunE: func(cmd *cobra.Command, args []string) error {
			// Initialize the handler with Redis client
			h, err := handler.New(
				handler.WithRedis(),
			)
			if err != nil {
				return fmt.Errorf("failed to initialize Redis client: %w", err)
			}

			// Test the Redis connection using RedisPing
			if err := h.RedisPing(context.Background()); err != nil {
				return fmt.Errorf("failed to connect to Redis: %w", err)
			}

			// Print success message if Ping is successful
			fmt.Println("Redis server is up and connection is successful!")
			return nil
		},
	}

	// GetMsg command
	GetMsgCmd = &cobra.Command{
		Use:   "getmsg",
		Short: "receives a message from the SQS queue",
		RunE: func(cmd *cobra.Command, args []string) error {
			// Create a handler instance with sqs
			h, err := handler.New(
				handler.WithSQS(viper.GetString("sqs.url")),
			)
			if err != nil {
				return fmt.Errorf("failed to initialize handler: %w", err)
			}

			// call  ReceiveMessage function
			success, err := h.ReceiveMessage(context.Background(), 30, 10, 1) // 30s visibility, 10s wait time, 1 message
			if err != nil {
				return fmt.Errorf("error while receiving message: %w", err)
			}

			if success {
				fmt.Println("Message successfully received and processed.")
			} else {
				fmt.Println("No messages found in the queue.")
			}

			return nil
		},
	}

	// WatchCmd calls the Watch function to periodically check the SQS queue
	WatchCmd = &cobra.Command{
		Use:   "watch",
		Short: "Watches the SQS queue periodically and processes new files",
		RunE: func(cmd *cobra.Command, args []string) error {
			// Initialize the handler with necessary clients
			h, err := handler.New(
				handler.WithSQS(viper.GetString("sqs.url")),
				//handler.WithS3(viper.GetString("s3.bucket"), viper.GetString("s3.endpoint")),
				//handler.WithMongoDB(),
				//handler.WithRedis(),
			)

			if err != nil {
				return fmt.Errorf("failed to initialize handler: %w", err)
			}

			// Call the Watch function with the interval from the flag
			if err := h.Watch(t); err != nil {
				return fmt.Errorf("error in Watch: %w", err)
			}

			return nil
		},
	}

	// lists all files and their versions
	ListCmd = &cobra.Command{
		Use:   "list",
		Short: "lists contents and versions in buckets",
		RunE: func(cmd *cobra.Command, args []string) error {
			// Get the bucket name and endpoint from configuration
			bucket := viper.GetString("s3.bucket")
			endpoint := viper.GetString("s3.endpoint")

			// check if config filled
			if bucket == "" {
				return fmt.Errorf("bucket must be set in the config file")
			}
			if endpoint == "" {
				return fmt.Errorf("endpoint must be set in the config file")
			}

			// configure with bucket and endpoint
			h, err := handler.New(
				handler.WithS3(bucket, endpoint),
			)
			if err != nil {
				return fmt.Errorf("failed to configure handler with S3 client: %v", err)
			}

			// list and err check
			if err := h.ListObjectVersions(); err != nil {
				return fmt.Errorf("failed to list object versions: %v", err)
			}

			return nil
		},
	}
)

func init() {
	viper.SetConfigName(".config") // name of config file (without extension)
	viper.SetConfigType("yaml")    // required since we're using .yml
	viper.AddConfigPath(".")       // look for the config fil

	// check to see if config file exists
	if err := viper.ReadInConfig(); err != nil {
		log.Printf("No configuration file found; using defaults or command-line args: %v", err)
	}

	// add commands to root cmd
	RootCmd.AddCommand(HandlerCmd)
	RootCmd.AddCommand(SQSClientCmd)
	RootCmd.AddCommand(GetMsgCmd)
	RootCmd.AddCommand(S3Cmd)
	RootCmd.AddCommand(ListCmd)
	RootCmd.AddCommand(MongoCmd)
	RootCmd.AddCommand(RedisCmd)
	RootCmd.AddCommand(WatchCmd)

	// Flags for TestCmd
	WatchCmd.PersistentFlags().IntVarP(&t, "time", "t", 5, "number of seconds to wait")

}

// Execute runs the RootCmd and handles any errors
func Execute() {
	err := RootCmd.Execute()
	if err != nil {
		log.Fatal(err)
	}
}
