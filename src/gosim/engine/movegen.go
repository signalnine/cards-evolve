package engine

import "encoding/binary"

// LegalMove represents a possible action
type LegalMove struct {
	PhaseIndex int
	CardIndex  int // -1 if not card-specific
	TargetLoc  Location
}

// GenerateLegalMoves returns all valid moves for current player
func GenerateLegalMoves(state *GameState, genome *Genome) []LegalMove {
	moves := make([]LegalMove, 0, 10)
	currentPlayer := state.CurrentPlayer

	for phaseIdx, phase := range genome.TurnPhases {
		switch phase.PhaseType {
		case 1: // DrawPhase
			if len(phase.Data) < 6 {
				continue
			}
			source := Location(phase.Data[0])
			mandatory := phase.Data[5] == 1

			// Check if can draw
			canDraw := false
			switch source {
			case LocationDeck:
				canDraw = len(state.Deck) > 0
			case LocationDiscard:
				canDraw = len(state.Discard) > 0
			case LocationOpponentHand:
				opponentID := 1 - currentPlayer
				canDraw = len(state.Players[opponentID].Hand) > 0
			}

			if canDraw || mandatory {
				moves = append(moves, LegalMove{
					PhaseIndex: phaseIdx,
					CardIndex:  -1,
					TargetLoc:  source,
				})
			}

		case 2: // PlayPhase
			if len(phase.Data) < 3 {
				continue
			}
			target := Location(phase.Data[0])
			minCards := int(phase.Data[1])
			maxCards := int(phase.Data[2])

			// For now, only support single-card plays
			if minCards <= 1 && maxCards >= 1 {
				// Check each card in hand
				for cardIdx := range state.Players[currentPlayer].Hand {
					// TODO: Evaluate valid_play_condition from phase.Data
					// For now, allow all cards
					moves = append(moves, LegalMove{
						PhaseIndex: phaseIdx,
						CardIndex:  cardIdx,
						TargetLoc:  target,
					})
				}
			}

		case 3: // DiscardPhase
			// Always allow discard if have cards
			if len(state.Players[currentPlayer].Hand) > 0 {
				for cardIdx := range state.Players[currentPlayer].Hand {
					moves = append(moves, LegalMove{
						PhaseIndex: phaseIdx,
						CardIndex:  cardIdx,
						TargetLoc:  LocationDiscard,
					})
				}
			}
		}
	}

	return moves
}

// ApplyMove executes a legal move, mutating state
func ApplyMove(state *GameState, move *LegalMove, genome *Genome) {
	if move.PhaseIndex >= len(genome.TurnPhases) {
		return
	}

	phase := genome.TurnPhases[move.PhaseIndex]
	currentPlayer := state.CurrentPlayer

	switch phase.PhaseType {
	case 1: // DrawPhase
		if len(phase.Data) >= 5 {
			count := int(binary.BigEndian.Uint32(phase.Data[1:5]))
			for i := 0; i < count; i++ {
				state.DrawCard(currentPlayer, move.TargetLoc)
			}
		}

	case 2: // PlayPhase
		if move.CardIndex >= 0 {
			state.PlayCard(currentPlayer, move.CardIndex, move.TargetLoc)

			// War-specific logic: if playing to tableau in 2-player game
			if move.TargetLoc == LocationTableau && len(state.Players) == 2 {
				resolveWarBattle(state)
			}
		}

	case 3: // DiscardPhase
		if move.CardIndex >= 0 {
			state.PlayCard(currentPlayer, move.CardIndex, LocationDiscard)
		}
	}

	// Advance turn
	state.CurrentPlayer = 1 - state.CurrentPlayer
	state.TurnNumber++
}

// resolveWarBattle handles War game card comparison
func resolveWarBattle(state *GameState) {
	// Check if both players have played (tableau has 2 cards)
	if len(state.Tableau) == 0 || len(state.Tableau[0]) < 2 {
		return
	}

	tableau := state.Tableau[0]
	card1 := tableau[len(tableau)-2] // Second-to-last card (player 0's card)
	card2 := tableau[len(tableau)-1] // Last card (player 1's card)

	// Compare ranks (Ace high: A=12, K=11, ..., 2=0)
	var winner uint8
	if card1.Rank > card2.Rank {
		winner = 0
	} else if card2.Rank > card1.Rank {
		winner = 1
	} else {
		// Tie - in simplified War, alternate who wins ties
		winner = state.CurrentPlayer
	}

	// Winner takes all cards from tableau
	for _, card := range tableau {
		state.Players[winner].Hand = append(state.Players[winner].Hand, card)
	}

	// Clear tableau
	state.Tableau[0] = state.Tableau[0][:0]
}

// CheckWinConditions evaluates win conditions, returns winner ID or -1
// Exported so mcts package can use it
func CheckWinConditions(state *GameState, genome *Genome) int8 {
	for _, wc := range genome.WinConditions {
		switch wc.WinType {
		case 0: // empty_hand
			for playerID, player := range state.Players {
				if len(player.Hand) == 0 {
					return int8(playerID)
				}
			}
		case 1: // high_score
			// TODO: Implement score-based win
		case 2: // first_to_score
			for playerID, player := range state.Players {
				if player.Score >= wc.Threshold {
					return int8(playerID)
				}
			}
		case 3: // capture_all
			for playerID, player := range state.Players {
				if len(player.Hand) == 52 {
					return int8(playerID)
				}
			}
		}
	}
	return -1
}
