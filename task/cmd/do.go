package cmd

import (
	"errors"
	"fmt"
	"log"
	"strconv"

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
		resetSeq()
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
		if b == nil {
			fmt.Println("Cannot do that, add something to list first")
			return errors.New("Empty database")
		}
		key := []byte(taskID)
		b.Delete(key)

		return nil
	})
}

func resetSeq() {
	db, err := bolt.Open(dbName, 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	err = db.Update(func(tx *bolt.Tx) error {
		var tasks [][]byte
		b := tx.Bucket([]byte(bucketName))
		// Get all remaining tasks
		b.ForEach(func(k, v []byte) error {
			tasks = append(tasks, v)
			b.Delete(k)
			return nil
		})
		// Reset start sequence to 0
		b.SetSequence(0)
		for _, v := range tasks {
			k, _ := b.NextSequence()
			b.Put([]byte(strconv.Itoa(int(k))), v)
		}
		return nil
	})
}
