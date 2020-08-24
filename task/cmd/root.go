package cmd

import (
	"fmt"
	"os"
	"path"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
)

const bucketName = "bezos"

var dbName string

func init() {
	home, herr := homedir.Dir()
	if herr != nil {
		fmt.Println("Could not get home directory")
	}
	dbName = path.Join(home, "list.db")
}

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
