package cmd

import (
	"fmt"
	"log"

	"github.com/dennis-yeom/robin/internal/handler"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
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

}

// Execute runs the RootCmd and handles any errors
func Execute() {
	err := RootCmd.Execute()
	if err != nil {
		log.Fatal(err)
	}
}
