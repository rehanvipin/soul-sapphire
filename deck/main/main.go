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
		deck.Remove('2', '3', "diamond", "club", "joker"),
		deck.Duplicate(2),
	)
	if err != nil {
		panic(err)
	}
	fmt.Println(deck)
}
