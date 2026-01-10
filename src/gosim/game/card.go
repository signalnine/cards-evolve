package game

import "fmt"

// Rank represents a card rank
type Rank int

const (
	Ace Rank = iota + 1
	Two
	Three
	Four
	Five
	Six
	Seven
	Eight
	Nine
	Ten
	Jack
	Queen
	King
)

// String returns the rank as a string
func (r Rank) String() string {
	ranks := []string{"", "A", "2", "3", "4", "5", "6", "7", "8", "9", "10", "J", "Q", "K"}
	return ranks[r]
}

// Suit represents a card suit
type Suit int

const (
	Hearts Suit = iota + 1
	Diamonds
	Clubs
	Spades
)

// String returns the suit as a string
func (s Suit) String() string {
	suits := []string{"", "H", "D", "C", "S"}
	return suits[s]
}

// Card represents a playing card
type Card struct {
	Rank Rank
	Suit Suit
}

// String returns the card as a string (e.g., "AH")
func (c Card) String() string {
	return fmt.Sprintf("%s%s", c.Rank.String(), c.Suit.String())
}

// NewDeck creates a standard 52-card deck
func NewDeck() []Card {
	deck := make([]Card, 0, 52)
	for suit := Hearts; suit <= Spades; suit++ {
		for rank := Ace; rank <= King; rank++ {
			deck = append(deck, Card{Rank: rank, Suit: suit})
		}
	}
	return deck
}
