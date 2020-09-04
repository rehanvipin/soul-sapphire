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
	fmt.Println("Dealer has", dealerHand)
	fmt.Println("Player has", playerHand)

	// Points table
	points := map[rune]int{
		'J': 10,
		'Q': 10,
		'K': 10,
	}
	for i := 2; i <= 10; i++ {
		points[rune('0'+i)] = i
	}

	// Calculate points
	var playerScore, dealerScore = 0, 0
	for i := range playerHand {
		playerScore += points[playerHand[i].Number]
	}
	for i := range dealerHand {
		dealerScore += points[dealerHand[i].Number]
	}

	fmt.Printf("Player: %2d\n", playerScore)
	fmt.Printf("Dealer: %2d\n", dealerScore)
}
