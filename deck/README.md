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