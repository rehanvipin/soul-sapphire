package main

import (
	"fmt"
	"soul-sapphire/deck"
)

func main() {
	deck, err := deck.New(
		deck.Sort(deck.DefaultCompare),
	)
	if err != nil {
		panic(err)
	}
	fmt.Println(deck)
}
