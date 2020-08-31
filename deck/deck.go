package deck

import (
	"sort"
)

// Card is the primitive type of the deck
type Card struct {
	number rune
	suit   string
}

// Deck is the type that is visible to the user
type Deck []Card

// Options is the base used by all functional options
type Options func(deck *Deck) (*Deck, error)

// Functional options, to be passed to New

// Sort sorts the deck according to some user defined function
func Sort(less func(Deck, int, int) bool) Options {
	return func(deck *Deck) (*Deck, error) {
		sort.SliceStable(*deck, func(i, j int) bool {
			return less(*deck, i, j)
		})
		return deck, nil
	}
}

// DefaultCompare is the default comparison function for sort
func DefaultCompare(deck Deck, first, second int) bool {
	var numberpoints = map[rune]int{'A': 1, '2': 2, '3': 3, '4': 4, '5': 5, '6': 6,
		'7': 7, '8': 8, '9': 9, 'J': 10, 'Q': 11, 'K': 12}
	var suitpoints = map[string]int{"spade": 1, "diamond": 2, "club": 3, "heart": 4}
	return suitpoints[deck[first].suit] > suitpoints[deck[second].suit] &&
		numberpoints[deck[first].number] > numberpoints[deck[second].number]
}

// New creates a new deck of cards and returns them
func New(opts ...Options) (Deck, error) {
	var deck = make(Deck, 0)
	numbers := []rune{'A', '2', '3', '4', '5', '6', '7', '8', '9', 'J', 'Q', 'K'}
	suits := []string{"spade", "diamond", "club", "heart"}
	for _, suit := range suits {
		for _, number := range numbers {
			deck = append(deck, Card{
				number: number,
				suit:   suit,
			})
		}
	}
	// Apply the options
	var prepDeck = &deck
	var err error
	for _, opt := range opts {
		prepDeck, err = opt(prepDeck)
		if err != nil {
			return nil, err
		}
	}
	deck = *prepDeck
	return deck, nil
}
