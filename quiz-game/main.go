package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"sync"
	"time"
)

var locked sync.Mutex

func main() {
	csvFile := flag.String("csv", "problems.csv", "Location of quiz Q&A")
	timeLimit := flag.Int("time", 30, "Timer for each question")
	shuffle := flag.Bool("random", false, "Suffle the order of questions")
	flag.Parse()

	f, err := os.Open(*csvFile)
	check(err)
	r := csv.NewReader(f)
	records, err := r.ReadAll()
	check(err)
	maxPoints := len(records)
	var points = 0

	// Shuffle if asked
	rand.Seed(time.Now().UnixNano())
	if *shuffle {
		rand.Shuffle(len(records), func(i, j int) {
			records[i], records[j] = records[j], records[i]
		})
	}

	// Press tab to start quiz
	var start rune
	fmt.Println("Press enter to start")
	for {
		fmt.Scanf("%c", &start)
		if start == rune('\n') {
			break
		}
	}

	var fiber = make(chan int)
	var timeout = false
	for _, x := range records {
		question, answer := x[0], x[1]
		go QuizQuestion(question, answer, fiber)
		// Wait until it finishes
		waiter := time.After(time.Duration(*timeLimit) * time.Second)
		var cont = false
		for !cont {
			select {
			case <-waiter:
				fmt.Println("Time up")
				timeout = true
				cont = true
			case i := <-fiber:
				points += i
				cont = true
			}
		}
		if timeout {
			break
		}
	}

	fmt.Printf("You got %d out of %d\n", points, maxPoints)
	f.Close()
}

// QuizQuestion quizzes on stdin for input and passes the score to the main thread
func QuizQuestion(question string, answer string, fiber chan int) {
	var guess string
	locked.Lock()
	fmt.Printf("What is %s ?: ", question)
	fmt.Scanln(&guess)
	locked.Unlock()
	if cleaner(guess) == cleaner(answer) {
		fiber <- 1
	} else {
		fiber <- 0
	}
}

func cleaner(input string) string {
	x := strings.Title(input)
	x = strings.ToLower(x)
	return x
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}
