package deck

import (
	"math/rand"
	"sort"
	"time"
)

// Card is the primitive type of the deck
// Joker is represented as {0, "joker"}
type Card struct {
	Number rune
	Suit   string
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
	return suitpoints[deck[first].Suit] > suitpoints[deck[second].Suit] &&
		numberpoints[deck[first].Number] > numberpoints[deck[second].Number]
}

// Shuffle randomly shuffles the deck
func Shuffle() Options {
	return func(deck *Deck) (*Deck, error) {
		rand.Seed(time.Now().UnixNano())
		tmp := *deck
		rand.Shuffle(len(*deck), func(i, j int) {
			tmp[i], tmp[j] = tmp[j], tmp[i]
		})
		return &tmp, nil
	}
}

// JokerAdd adds n jokers to the deck
func JokerAdd(n int) Options {
	return func(deck *Deck) (*Deck, error) {
		tmp := *deck
		for i := 0; i < n; i++ {
			tmp = append(tmp, Card{0, "joker"})
		}
		return &tmp, nil
	}
}

func del(i int, deck Deck) Deck {
	copy(deck[i:], deck[i+1:])
	deck[len(deck)-1] = Card{}
	deck = deck[:len(deck)-1]
	return deck
}

// Remove removes a set of cards from the deck
func Remove(matches ...interface{}) Options {
	return func(deck *Deck) (*Deck, error) {
		for _, match := range matches {
			slow := 0
			for i := range *deck {
				switch v := match.(type) {
				case rune:
					if (*deck)[i].Number != match {
						(*deck)[slow] = (*deck)[i]
						slow++
					}
				case string:
					if (*deck)[i].Suit != match {
						(*deck)[slow] = (*deck)[i]
						slow++
					}
				case int:
					if match != 0 {
						slow++
						break
					}
					if (*deck)[i].Number != rune(0) {
						(*deck)[slow] = (*deck)[i]
						slow++
					}
				default:
					// Sneak tricks
					_ = v
					slow++
				}
			}
			// fix mem leaks
			for i := slow; i < len(*deck); i++ {
				(*deck)[i] = Card{}
			}
			(*deck) = (*deck)[:slow]
		}
		return deck, nil
	}
}

// Duplicate duplicates the deck n times
func Duplicate(n int) Options {
	return func(deck *Deck) (*Deck, error) {
		tmp := *deck
		for i := 0; i < n-1; i++ {
			tmp = append(tmp, tmp...)
		}
		return &tmp, nil
	}
}

// New creates a new deck of cards and returns them
func New(opts ...Options) (Deck, error) {
	var deck = make(Deck, 0)
	numbers := []rune{'A', '2', '3', '4', '5', '6', '7', '8', '9', 'J', 'Q', 'K'}
	suits := []string{"spade", "diamond", "club", "heart"}
	for _, suit := range suits {
		for _, number := range numbers {
			deck = append(deck, Card{
				Number: number,
				Suit:   suit,
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
