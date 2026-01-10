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
go test ./game -v
go test ./game -bench=. -benchtime=10s
```

### Benchmarks

```bash
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
