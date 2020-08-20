package cmd

import (
	"fmt"
	"log"

	"github.com/boltdb/bolt"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(listCmd)
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List out all the tasks",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("The tasks are")
		// Get all tasks from the db
		enumerate()
	},
}

func enumerate() {
	db, err := bolt.Open("list.db", 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	err = db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucketName))

		b.ForEach(func(k, v []byte) error {
			fmt.Printf("k:%s & v:%s\n", k, v)
			return nil
		})
		return nil
	})
}
