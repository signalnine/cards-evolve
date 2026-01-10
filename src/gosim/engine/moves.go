package engine

// DrawCard moves a card from source to player hand
func (s *GameState) DrawCard(playerID uint8, source Location) bool {
	var srcPile *[]Card

	switch source {
	case LocationDeck:
		srcPile = &s.Deck
	case LocationDiscard:
		srcPile = &s.Discard
	case LocationOpponentHand:
		// Optional extension: draw from opponent's hand
		opponentID := 1 - playerID
		srcPile = &s.Players[opponentID].Hand
	case LocationOpponentDiscard:
		// Optional extension: draw from opponent's discard (not standard)
		// Would need per-player discard piles
		return false
	default:
		return false
	}

	if len(*srcPile) == 0 {
		return false
	}

	// Pop from source
	card := (*srcPile)[len(*srcPile)-1]
	*srcPile = (*srcPile)[:len(*srcPile)-1]

	// Add to player hand
	s.Players[playerID].Hand = append(s.Players[playerID].Hand, card)
	return true
}

// PlayCard moves a card from player hand to target location
func (s *GameState) PlayCard(playerID uint8, cardIndex int, target Location) bool {
	hand := &s.Players[playerID].Hand

	if cardIndex < 0 || cardIndex >= len(*hand) {
		return false
	}

	// Remove from hand
	card := (*hand)[cardIndex]
	*hand = append((*hand)[:cardIndex], (*hand)[cardIndex+1:]...)

	// Add to target
	switch target {
	case LocationDiscard:
		s.Discard = append(s.Discard, card)
	case LocationTableau:
		if len(s.Tableau) == 0 {
			s.Tableau = append(s.Tableau, make([]Card, 0, 10))
		}
		s.Tableau[0] = append(s.Tableau[0], card)
	default:
		return false
	}

	return true
}

// ShuffleDeck randomizes deck order (in-place)
func (s *GameState) ShuffleDeck(seed uint64) {
	// Simple LCG for deterministic shuffle
	rng := seed
	n := len(s.Deck)

	for i := n - 1; i > 0; i-- {
		rng = rng*6364136223846793005 + 1442695040888963407
		j := int(rng % uint64(i+1))
		s.Deck[i], s.Deck[j] = s.Deck[j], s.Deck[i]
	}
}
