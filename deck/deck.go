package deck

// Card is the primitive type of the deck
type Card struct {
	number rune
	suit   string
}

// Deck is the type that is visible to the user
type Deck []Card

// New creates a new deck of cards and returns them
func New() Deck {
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
	return deck
}
