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
