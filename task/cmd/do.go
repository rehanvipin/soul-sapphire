package cmd

import (
	"fmt"
	"log"

	"github.com/boltdb/bolt"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(doCmd)
}

var doCmd = &cobra.Command{
	Use:   "do [task-no]",
	Short: "Check off a task from the list",
	Run: func(cmd *cobra.Command, args []string) {
		for _, task := range args {
			complete(task)
			// Remove required task from the db
			fmt.Println("Completed task", task)
		}
	},
}

func complete(taskID string) {
	db, err := bolt.Open(dbName, 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucketName))
		key := []byte(taskID)
		b.Delete(key)
		return nil
	})
}
