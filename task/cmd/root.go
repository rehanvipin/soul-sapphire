package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

const bucketName = "bezos"

var rootCmd = &cobra.Command{
	Use:   "task",
	Short: "Task is a CLI To-do list",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Run task help to see usage")
	},
}

// Execute is the interface to run the main app
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
