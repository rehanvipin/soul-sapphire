package cmd

import (
	"errors"
	"fmt"

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
	db, err := bolt.Open(dbName, 0600, nil)
	if err != nil {
		fmt.Println("Could not find the database file")
	}
	defer db.Close()

	// Create bucket if it does not exist
	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucketName))
		if b == nil {
			_, makerr := tx.CreateBucket([]byte(bucketName))
			if makerr != nil {
				fmt.Println("Could not make bucket")
				return makerr
			}
		}
		return nil
	})

	if err != nil {
		fmt.Println("Cannot use db")
	}

	err = db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucketName))
		if b == nil {
			fmt.Println("Could get such a bucket")
			return errors.New("Bucket not found")
		}

		b.ForEach(func(k, v []byte) error {
			fmt.Printf("%s - %s\n", k, v)
			return nil
		})
		return nil
	})
}
