package cmd

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/boltdb/bolt"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(addCmd)
}

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add task to the list",
	Run: func(cmd *cobra.Command, args []string) {
		// Connect to db and add task to list
		task := strings.Join(args, " ")
		addTask(task)
		// Confirmation
		fmt.Printf("Added task - %v - to the list\n", task)
	},
}

func addTask(task string) {
	db, err := bolt.Open(dbName, 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	value := []byte(task)
	var keyb []byte

	err = db.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte(bucketName))
		if err != nil {
			return err
		}
		key, kerr := bucket.NextSequence()
		if kerr != nil {
			return kerr
		}

		keyb = []byte(strconv.Itoa(int(key)))

		err = bucket.Put(keyb, value)
		if err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		log.Fatal(err)
	}
}

// func check() {
// 	// Retrieve data
// 	err = db.View(func(tx *bolt.Tx) error {
// 		bucket := tx.Bucket([]byte(bucketName))
// 		if bucket == nil {
// 			return fmt.Errorf("could not find bucket %v", bucket)
// 		}

// 		val := bucket.Get(keyb)
// 		fmt.Println(string(val))

// 		return nil
// 	})
// }
