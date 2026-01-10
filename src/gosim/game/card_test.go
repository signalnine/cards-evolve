package game

import "testing"

func TestCard_String(t *testing.T) {
	tests := []struct {
		card Card
		want string
	}{
		{Card{Rank: Ace, Suit: Hearts}, "AH"},
		{Card{Rank: King, Suit: Spades}, "KS"},
		{Card{Rank: Two, Suit: Diamonds}, "2D"},
	}

	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			if got := tt.card.String(); got != tt.want {
				t.Errorf("Card.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewDeck(t *testing.T) {
	deck := NewDeck()

	if len(deck) != 52 {
		t.Errorf("NewDeck() returned %d cards, want 52", len(deck))
	}

	// Check for duplicates
	seen := make(map[string]bool)
	for _, card := range deck {
		key := card.String()
		if seen[key] {
			t.Errorf("Duplicate card found: %s", key)
		}
		seen[key] = true
	}
}
