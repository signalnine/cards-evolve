package engine

import "encoding/binary"

// EvaluateCondition checks if condition is true for given state
func EvaluateCondition(state *GameState, playerID uint8, conditionBytes []byte) bool {
	if len(conditionBytes) < 7 {
		return false
	}

	opcode := OpCode(conditionBytes[0])
	operator := conditionBytes[1]
	value := int32(binary.BigEndian.Uint32(conditionBytes[2:6]))
	reference := conditionBytes[6]

	var actual int32

	switch opcode {
	case OpCheckHandSize:
		actual = int32(len(state.Players[playerID].Hand))

	case OpCheckLocationSize:
		switch Location(reference) {
		case LocationDeck:
			actual = int32(len(state.Deck))
		case LocationDiscard:
			actual = int32(len(state.Discard))
		case LocationTableau:
			if len(state.Tableau) > 0 {
				actual = int32(len(state.Tableau[0]))
			}
		}

	case OpCheckCardRank:
		// Check if card at index matches rank
		refCard := getReferencedCard(state, reference)
		if refCard != nil && int(refCard.Rank) == int(value) {
			return true
		}
		return false

	case OpCheckCardSuit:
		refCard := getReferencedCard(state, reference)
		if refCard != nil && int(refCard.Suit) == int(value) {
			return true
		}
		return false

	// Optional extensions: betting conditions
	case OpCheckChipCount:
		actual = state.Players[playerID].Chips

	case OpCheckPotSize:
		actual = state.Pot

	case OpCheckCurrentBet:
		actual = state.CurrentBet

	case OpCheckCanAfford:
		actual = state.Players[playerID].Chips
		// Check if player can afford the value
		return actual >= value

	default:
		return false
	}

	// Apply operator
	switch OpCode(operator + 50) {
	case OpEQ:
		return actual == value
	case OpNE:
		return actual != value
	case OpLT:
		return actual < value
	case OpGT:
		return actual > value
	case OpLE:
		return actual <= value
	case OpGE:
		return actual >= value
	default:
		return false
	}
}

func getReferencedCard(state *GameState, reference uint8) *Card {
	switch reference {
	case 1: // top_discard
		if len(state.Discard) > 0 {
			return &state.Discard[len(state.Discard)-1]
		}
	case 2: // last_played (tableau top)
		if len(state.Tableau) > 0 && len(state.Tableau[0]) > 0 {
			pile := state.Tableau[0]
			return &pile[len(pile)-1]
		}
	}
	return nil
}
