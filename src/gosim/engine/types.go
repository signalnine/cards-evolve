package engine

import (
	"sync"
)

// Card represents a playing card (1 byte)
type Card struct {
	Rank uint8 // 0-12 (A,2-10,J,Q,K)
	Suit uint8 // 0-3 (H,D,C,S)
}

// Location enum
type Location uint8

const (
	LocationDeck Location = iota
	LocationHand
	LocationDiscard
	LocationTableau
	// Optional extensions
	LocationOpponentHand
	LocationOpponentDiscard
)

// PlayerState is mutable for performance
type PlayerState struct {
	Hand   []Card
	Score  int32
	Active bool // Still in the game (not folded/eliminated)
	// Optional extensions for betting games
	Chips      int32 // Chip/token count for betting games
	CurrentBet int32 // Current bet in this round
	HasFolded  bool  // Folded this round
}

// GameState is mutable and pooled
type GameState struct {
	Players       []PlayerState
	Deck          []Card
	Discard       []Card
	Tableau       [][]Card // For games like War, Gin Rummy
	CurrentPlayer uint8
	TurnNumber    uint32
	WinnerID      int8 // -1 = no winner yet, 0/1 = player ID
	// Optional extensions for betting games
	Pot        int32 // Current pot size
	CurrentBet int32 // Highest bet in current round
}

// StatePool manages GameState memory
var StatePool = sync.Pool{
	New: func() interface{} {
		return &GameState{
			Players: make([]PlayerState, 2),
			Deck:    make([]Card, 0, 52),
			Discard: make([]Card, 0, 52),
			Tableau: make([][]Card, 0, 10),
		}
	},
}

// GetState acquires a GameState from pool
func GetState() *GameState {
	state := StatePool.Get().(*GameState)
	state.Reset()
	return state
}

// PutState returns a GameState to pool
func PutState(state *GameState) {
	StatePool.Put(state)
}

// Reset clears state for reuse
func (s *GameState) Reset() {
	s.Players[0].Hand = s.Players[0].Hand[:0]
	s.Players[0].Score = 0
	s.Players[0].Active = true
	s.Players[0].Chips = 0
	s.Players[0].CurrentBet = 0
	s.Players[0].HasFolded = false

	s.Players[1].Hand = s.Players[1].Hand[:0]
	s.Players[1].Score = 0
	s.Players[1].Active = true
	s.Players[1].Chips = 0
	s.Players[1].CurrentBet = 0
	s.Players[1].HasFolded = false

	s.Deck = s.Deck[:0]
	s.Discard = s.Discard[:0]
	s.Tableau = s.Tableau[:0]
	s.CurrentPlayer = 0
	s.TurnNumber = 0
	s.WinnerID = -1
	s.Pot = 0
	s.CurrentBet = 0
}

// Clone creates a deep copy for MCTS tree search
func (s *GameState) Clone() *GameState {
	clone := GetState()

	clone.Players[0].Hand = append(clone.Players[0].Hand, s.Players[0].Hand...)
	clone.Players[0].Score = s.Players[0].Score
	clone.Players[0].Active = s.Players[0].Active
	clone.Players[0].Chips = s.Players[0].Chips
	clone.Players[0].CurrentBet = s.Players[0].CurrentBet
	clone.Players[0].HasFolded = s.Players[0].HasFolded

	clone.Players[1].Hand = append(clone.Players[1].Hand, s.Players[1].Hand...)
	clone.Players[1].Score = s.Players[1].Score
	clone.Players[1].Active = s.Players[1].Active
	clone.Players[1].Chips = s.Players[1].Chips
	clone.Players[1].CurrentBet = s.Players[1].CurrentBet
	clone.Players[1].HasFolded = s.Players[1].HasFolded

	clone.Deck = append(clone.Deck, s.Deck...)
	clone.Discard = append(clone.Discard, s.Discard...)

	for _, pile := range s.Tableau {
		tableuClone := make([]Card, len(pile))
		copy(tableuClone, pile)
		clone.Tableau = append(clone.Tableau, tableuClone)
	}

	clone.CurrentPlayer = s.CurrentPlayer
	clone.TurnNumber = s.TurnNumber
	clone.WinnerID = s.WinnerID
	clone.Pot = s.Pot
	clone.CurrentBet = s.CurrentBet

	return clone
}
