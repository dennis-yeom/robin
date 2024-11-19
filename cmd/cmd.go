package cmd

import (
	"fmt"
	"log"

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
)

func init() {
	viper.SetConfigName(".config") // name of config file (without extension)
	viper.SetConfigType("yaml")    // required since we're using .yml
	viper.AddConfigPath(".")       // look for the config fil

	// check to see if config file exists
	if err := viper.ReadInConfig(); err != nil {
		log.Printf("No configuration file found; using defaults or command-line args: %v", err)
	}

}

// Execute runs the RootCmd and handles any errors
func Execute() {
	err := RootCmd.Execute()
	if err != nil {
		log.Fatal(err)
	}
}
