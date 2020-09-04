package main

import (
	"fmt"
	"log"
	"soul-sapphire/deck"
)

func main() {
	pack, initerr := deck.New(
		deck.Duplicate(3),
		deck.Shuffle(),
	)
	if initerr != nil {
		log.Fatalln("Could not create such a pack of cards")
	}

	// Create hands
	var dealerHand, playerHand deck.Deck

	// Deal
	var card deck.Card
	for i := 0; i < 2; i++ {
		card = pack[0]
		playerHand = append(playerHand, card)
		card = pack[1]
		dealerHand = append(dealerHand, card)
		pack = pack[2:]
	}
	fmt.Println(dealerHand)
	fmt.Println(playerHand)
}
