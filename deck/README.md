# Deck - making decks of cards
A package to created specific card decks.  
Created as a part of [Gophercises](https://gophercises.com)

## Details:
* A chance to use functional options in full [Tutorial](https://www.calhoun.io/using-functional-options-instead-of-method-chaining-in-go/)
* Check [Description](Description.md) for details

## Usage:
1. Import as `soul-sapphire/deck`
2. Card is the primary type. Deck is a slice of cards

* ### With functional options:
    Custom functions can be provided within the `New` function. Such as:
    * `Sort` to sort the deck with a custom function, e.g. `deck.Sort(deck.DefaultCompare)`
    * `Shuffle` to shuffle the deck of cards
    * `JokerAdd` adds n number of jokers to the deck
    * `Remove` allows to remove sequences/particular cards from the deck e.g. `[]rune{'2', '3', 0}, []string{"diamond", "club"}, []Card{Card{0, "joker"}}`
    * `Duplicate` duplicates the deck n times

3. The functional options can be called in any order
4. An example call
```go
deck, err := deck.New(
		deck.Sort(deck.DefaultCompare),
		deck.Shuffle(),
		deck.JokerAdd(2),
		deck.Remove([]rune{'2', '3', 0}, []string{"diamond", "club"}, nil),
		deck.Duplicate(2),
	)
```