package main

import (
	"fmt"
	"soul-sapphire/deck"
)

func main() {
	deck, err := deck.New(
		deck.Sort(deck.DefaultCompare),
		deck.Shuffle(),
		deck.JokerAdd(2),
		deck.Remove([]rune{'2', '3', 0}, []string{"diamond", "club"}, nil),
		deck.Duplicate(2),
	)
	if err != nil {
		panic(err)
	}
	fmt.Println(deck)
}
