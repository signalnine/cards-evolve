# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is an evolutionary computation system that uses genetic algorithms and Monte Carlo simulations to evolve novel card games playable with a standard 52-card deck. The system optimizes for configurable parameters like complexity, session length, skill/luck ratio, and player count.

## Core Architecture Concepts

### Game Representation
Games are encoded as genomes containing:
- **Setup rules**: Deal patterns, hand sizes, tableau configuration
- **Turn structure**: Draw/play/pass rules
- **Valid moves**: Card matching logic (suit, rank, sequence), special card effects
- **Win conditions**: Empty hand, point thresholds, capture goals
- **Scoring system**: Point values, bonuses, penalties

### Three-Layer System
1. **Rule DSL**: Domain-specific language for defining card game rules
2. **Simulation Engine**: Monte Carlo runner with multiple AI player types (random, greedy, MCTS)
3. **Genetic Algorithm**: Evolution engine with mutation, crossover, and selection operators

### Fitness Evaluation
Games are scored on measurable proxies for "fun":
- Decision density (meaningful choices vs forced plays)
- Comeback potential (trailing players can recover)
- Tension curve (uncertainty over time)
- Interaction frequency (actions affecting opponents)
- Rules complexity and session length
- Skill vs luck ratio (MCTS win rate vs random)

## Key Constraints

All evolved games must be:
- **Playable**: No infinite loops or unreachable win states
- **Terminable**: Enforce maximum turn limits
- **Agentic**: Contain non-random decision points

## Development Approach

When implementing, follow this sequence:
1. Rule DSL design and parser
2. Simulation harness with random AI baseline
3. Game state representation and validation
4. Genetic operators (mutation, crossover, selection)
5. Advanced AI players (greedy heuristics, MCTS)
6. Fitness function implementation
7. Natural language rule generator for human playtesting

## Validation Strategy

- Random AI validates games are mechanically playable
- Greedy AI measures obvious strategy effectiveness
- MCTS approximates skilled play
- Skill gap = MCTS win rate differential vs random baseline
- Human playtesting validates proxy metrics correlate with actual enjoyment

## Performance Benchmarks

### Golang Core Validation (Phase 1)

**War Game Benchmark:**
- Python implementation: 0.07ms per game
- Golang implementation: 0.03ms per game
- Measured speedup: 2.9x

**Interface Decision:** CGo

**Rationale:** Despite modest speedup on simple War game, Golang provides measurable performance benefit. More complex simulations (MCTS, deep game trees) will show greater advantages. CGo chosen for tight integration without serialization overhead, critical for millions of evolutionary iterations.

### Golang Performance Core (Phase 3)

**Architecture:** Python→Flatbuffers→CGo→Go

**Implementation Status:** Core architecture complete, simulation logic in progress

**Components Implemented:**
1. **Bytecode Compiler** (`src/cards_evolve/genome/bytecode.py`)
   - Compiles GameGenome to 36-byte header + phase data
   - OpCodes for conditions (0-19), actions (20-39), control flow (40-49), operators (50-55)
   - War genome compiles to 77 bytes
   - Deterministic compilation (same genome = same bytecode)

2. **Flatbuffers Schema** (`schema/simulation.fbs`)
   - Zero-copy binary serialization for Python↔Go
   - BatchRequest/BatchResponse for bulk simulation
   - Supports Random/Greedy/MCTS AI types (100/500/1000/2000 iterations)

3. **CGo Bridge** (`src/gosim/cgo/bridge.go`, `src/cards_evolve/bindings/cgo_bridge.py`)
   - `libcardsim.so` shared library (1.8MB)
   - `SimulateBatch` entry point for batch processing
   - Memory management via `FreeCString`
   - Built with: `make build-cgo`

4. **Mutable GameState** (`src/gosim/engine/types.go`)
   - sync.Pool memory pooling for zero-allocation reuse
   - Card, PlayerState, GameState types with betting extensions
   - DrawCard, PlayCard, ShuffleDeck operations
   - Clone() for tree search
   - 3/3 unit tests passing

5. **Genome Interpreter** (`src/gosim/engine/bytecode.go`, `conditions.go`, `movegen.go`)
   - ParseGenome extracts header + phases + win conditions
   - EvaluateCondition handles hand size, location size, card rank/suit, betting
   - GenerateLegalMoves produces move list for current phase
   - ApplyMove mutates state in-place
   - CheckWinConditions returns winner ID
   - 7/7 unit tests passing

6. **MCTS Engine** (`src/gosim/mcts/node.go`, `search.go`)
   - MCTSNode with sync.Pool memory pooling
   - UCB1 selection algorithm (configurable exploration parameter)
   - Selection → Expansion → Simulation → Backpropagation
   - Returns most-visited child as best move
   - Benchmark: ~3ms per search (100 iterations)
   - 8/8 unit tests passing

7. **Batch Simulation Engine** (`src/gosim/simulation/runner.go`)
   - RunBatch: Execute N games with same genome/AI config
   - Random AI: Uniform selection from legal moves
   - Greedy AI: Heuristic-based move scoring
   - MCTS AI: Tree search with configurable iterations
   - Aggregated statistics: wins, avg/median turns, duration
   - Turn limit protection against infinite loops

8. **Golden Test Suite** (`tests/integration/test_bytecode_equivalence.py`, `src/gosim/engine/bytecode_test.go`)
   - Python bytecode compilation tests (5 passing)
   - Go bytecode parsing tests (4 passing)
   - `tests/golden/war_genome.bin` (77 bytes) validates Python↔Go equivalence
   - Determinism verification

**Performance Status:**
- MCTS benchmark: ~3ms per search (100 iterations) on Intel N100
- Full end-to-end benchmark pending simulation logic fixes

**Known Issues:**
- War simulation ending prematurely with "no legal moves"
- Move generation for PlayPhase needs validation fixes
- Integration tests require Python environment setup fixes

**Next Steps:**
1. Debug War simulation move generation logic
2. Verify win condition triggers correctly
3. Run full end-to-end benchmark suite
4. Compare Python vs Go performance on 1000+ game batches
5. Validate MCTS vs Random AI skill differential

## Development Commands

### Run Tests

**Python:**
```bash
uv run pytest tests/ -v
uv run pytest tests/unit/test_specific.py -v  # Single file
```

**Golang:**
```bash
cd src/gosim
go test ./engine -v
go test ./mcts -v
go test ./simulation -v

# Benchmarks
go test ./mcts -bench=. -benchtime=3s
go test ./simulation -bench=. -benchtime=10s
```

### Build CGo Library

```bash
make build-cgo  # Builds libcardsim.so
```

### Benchmarks

```bash
# Phase 3: Golang performance core
uv run python benchmarks/benchmark_golang.py

# Phase 1: Python vs Go comparison
uv run python benchmarks/compare_war.py
```

### Code Quality

```bash
uv run black src/ tests/
uv run mypy src/
```

## Project Structure

```
cards-evolve/
├── src/
│   ├── cards_evolve/          # Python package
│   │   ├── genome/            # Game genome representation
│   │   ├── simulation/        # Game simulation engines
│   │   ├── evolution/         # Genetic algorithm
│   │   └── cli/               # Command-line interface
│   └── gosim/                 # Golang simulation core
│       └── game/              # Card game primitives
├── tests/
│   ├── unit/                  # Unit tests
│   ├── integration/           # Integration tests
│   └── property/              # Property-based tests (Hypothesis)
├── benchmarks/                # Performance comparisons
└── docs/
    ├── architecture/          # Architecture decisions
    └── plans/                 # Implementation plans
```

## Phase 2: Python Simulation Core Patterns

### Immutability Patterns

**Always use tuples for nested structures:**
```python
# ✅ Correct
@dataclass(frozen=True)
class GameState:
    deck: tuple[Card, ...]
    hands: tuple[tuple[Card, ...], ...]

# ❌ Wrong - frozen wrapper around mutable list
@dataclass(frozen=True)
class GameState:
    deck: list[Card]  # NOT immutable!
```

**Use copy_with for state transitions:**
```python
new_state = old_state.copy_with(
    turn=old_state.turn + 1,
    active_player=(old_state.active_player + 1) % 2
)
```

### Interpreter Pattern

**GenomeInterpreter converts data to logic (NOT code generation):**
- Genomes are pure dataclasses
- GameLogic instantiates based on genome fields
- Safe for pickling across processes
- No `exec()`, no security risks

### Testing Strategies

**Property-based tests with Hypothesis:**
```python
@given(seed=st.integers(min_value=0, max_value=10000))
def test_determinism_property(seed: int) -> None:
    result1 = engine.simulate_game(genome, players, seed)
    result2 = engine.simulate_game(genome, players, seed)
    assert result1.winner == result2.winner
```

**Sanity checks for metrics:**
- War game should have decision_density ≈ 0.0
- If not, metric calculation is broken

### Two-Phase Fitness Evaluation

**Phase 1 (cheap, all candidates):**
- Game length, completion rate, decision branch factor

**Phase 2 (expensive, top 10-20%):**
- Tension curve, comeback potential (deferred to Phase 3)

### Rank Comparison

**Use numeric mapping for Rank enum:**
```python
# Rank enum values are strings ("A", "2", "K")
# For comparison, use numeric mapping:
RANK_VALUES = {
    Rank.TWO: 2,
    Rank.THREE: 3,
    # ...
    Rank.ACE: 14,  # Ace high
}

def get_rank_value(card: Card) -> int:
    return RANK_VALUES[card.rank]
```
