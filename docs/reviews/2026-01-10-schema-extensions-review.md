# Multi-Agent Consensus Analysis

**Mode:** general-prompt
**Date:** 2026-01-10 11:23:08
**Agents Succeeded:** 3/3
**Chairman:** Claude

---

## Original Question

Review the extended genome schema for a card game evolution system.

SCHEMA LOCATION: /home/gabe/cards-playtest/docs/genome-schema-examples.md

CONTEXT:
This schema represents card games as structured genomes (Python dataclasses) that can be:
- Simulated by a game engine
- Compiled to bytecode for a Golang performance core
- Evolved using genetic algorithms
- Validated against real games from Hoyle's Encyclopedia

BASE SCHEMA (Core types):
- Enumerations: Rank, Suit, Location, ConditionType, Operator, ActionType
- Condition system: Composable predicates for game logic
- Action system: Executable game actions
- Game structure: SetupRules, TurnStructure (phases), SpecialEffects, WinConditions, ScoringRules
- Phases: DrawPhase, PlayPhase, DiscardPhase

OPTIONAL EXTENSIONS (added today):
1. Opponent Interaction:
   - Location.OPPONENT_HAND, Location.OPPONENT_DISCARD
   - ActionType.DRAW_FROM_OPPONENT, ActionType.DISCARD_PAIRS
   - SetupRules.post_deal_actions (for Old Maid initial pairing)

2. Set/Collection Detection:
   - ConditionType.HAS_SET_OF_N (Go Fish books of 4)
   - ConditionType.HAS_RUN_OF_N (Gin Rummy sequential cards)
   - ConditionType.HAS_MATCHING_PAIR (Old Maid pairs)
   - DiscardPhase.matching_condition (constrain discards)

3. Betting/Wagering:
   - ResourceRules (chip tracking, starting_chips, ante, blinds)
   - BettingPhase (min_bet, max_bet, allow_check/raise/fold, raise limits)
   - Actions: BET, CALL, RAISE, FOLD, CHECK, ALL_IN
   - Conditions: CHIP_COUNT, POT_SIZE, CURRENT_BET, CAN_AFFORD

4. Bluffing/Challenges:
   - ClaimPhase (claim_types, can_lie, penalties)
   - Actions: CLAIM, CHALLENGE, REVEAL
   - Conditional transfers based on claim truthfulness

EXAMPLE GAMES ENCODED:
- Crazy 8s (base schema)
- War (base schema)
- Gin Rummy (base schema with set/run detection)
- Old Maid (opponent interaction + pairing detection)
- Go Fish (set detection + opponent interaction)
- Betting War (chip tracking + betting phase)
- I Doubt It / Cheat (claim phase + challenge mechanics)

SCHEMA COVERAGE:
- Base schema: 60-70% of simple card games
- With extensions: 85-90% of simple card games

DESIGN QUESTIONS:

1. COMPLETENESS: Are the extensions sufficient for the stated coverage (85-90%)? Are there obvious gaps?

2. BACKWARD COMPATIBILITY: Are extensions truly optional and non-breaking? Can base games (War, Crazy 8s) still work without extensions?

3. EVOLUTION FRIENDLINESS: Can genetic algorithms effectively explore this schema space? Are the structures too rigid or too flexible?

4. IMPLEMENTATION COMPLEXITY: Are the extensions implementable in both Python and Go (via bytecode)? Are there hidden complexities?

5. VALIDATION: Can we effectively validate genomes to ensure they're playable (no infinite loops, termination guaranteed)?

6. SEMANTIC CLARITY: Are the extension semantics clear? For example, does OPPONENT_HAND clearly mean 'draw from opponent's hand'?

7. MISSING PATTERNS: What common card game patterns are still NOT representable? What would be needed to support them?

8. SCHEMA DESIGN: Is the dataclass-based approach (Path A) better than alternatives (DSL, behavior trees, bytecode-first)?

Please provide:
A. Strengths of the current design
B. Weaknesses or gaps identified
C. Ambiguities that need clarification
D. Recommendations for improvements (if any)
E. Assessment of whether this schema can achieve the 85-90% coverage claim


## Context

Full schema file available at /home/gabe/cards-playtest/docs/genome-schema-examples.md

Example games documented in /home/gabe/cards-playtest/docs/hoyles-game-examples.md

Key constraint: Must be evolvable (genetic algorithms can mutate/crossover)
Key constraint: Must compile to bytecode for Go execution
Key constraint: Must be playable with standard 52-card deck


---

## Stage 1: Independent Analyses

### Claude Review

# Independent Analysis: Extended Genome Schema for Card Game Evolution System

## Strong Points

**1. Composable Condition/Action Architecture**
The separation of `Condition` and `Action` systems with composable predicates is excellent for genetic algorithms. This provides clear mutation points (swap conditions, modify operators) and natural crossover boundaries (exchange entire phase definitions). The enumeration-based design makes the search space discrete and bounded.

**2. Backward Compatibility Design**
The extensions appear genuinely optional—base games using only core `Location`, `ActionType`, and `ConditionType` enums shouldn't require any extension-specific logic. The `post_deal_actions` as an optional field in `SetupRules` is a clean pattern.

**3. Phase-Based Turn Structure**
`DrawPhase`, `PlayPhase`, `DiscardPhase` as explicit structures provide semantic clarity for both human understanding and bytecode compilation. This maps well to real card games and provides natural phase boundaries for validation.

**4. Betting Extension Modularity**
`ResourceRules` and `BettingPhase` as separate structures (rather than baked into core types) is smart. Poker variants can reuse the same betting mechanics with different hand evaluation—good separation of concerns.

**5. Claim/Challenge Mechanics**
The `ClaimPhase` with `can_lie` and `penalties` captures bluffing games (I Doubt It, Cheat) cleanly. The truthfulness-conditional transfers are a sophisticated pattern that many schemas miss.

---

## Moderate Points

**1. Set/Run Detection Abstraction**
`HAS_SET_OF_N` and `HAS_RUN_OF_N` are useful but may need parameterization beyond N. Questions arise:
- Does `HAS_RUN_OF_N` require same suit? (Gin vs. some variants)
- Can you specify "run in hand" vs "run on table"?
- Are wildcards supported in set detection?

The semantic clarity is moderate—needs explicit documentation of suit constraints and location scope.

**2. Opponent Interaction Scope**
`OPPONENT_HAND` and `DRAW_FROM_OPPONENT` work for 2-player games but raise questions for multiplayer:
- Which opponent? (left, right, chosen, random?)
- Does the schema support directed actions (ask specific player in Go Fish)?

This may limit coverage to 2-player variants without additional parameterization.

**3. Coverage Claim (85-90%)**
Plausible for "simple" card games but depends heavily on definition of "simple." The schema appears to cover:
- Shedding games (Crazy 8s, Uno patterns) ✓
- Trick-taking basics (War) ✓
- Collection games (Go Fish, Gin) ✓
- Bluffing games (I Doubt It) ✓
- Simple betting (Betting War) ✓

But may struggle with: partnership games, trump systems, auction/bidding, tableau building.

**4. Bytecode Compilation Path**
Dataclass → bytecode → Go is implementable, but the complexity depends on:
- How conditions are evaluated (interpreter vs compiled predicates)
- Whether genetic mutations can produce invalid bytecode
- Stack depth guarantees for nested conditions

Moderate confidence—workable but will need careful implementation.

---

## Weak Points / Concerns

**1. No Explicit Trick-Taking Mechanics**
The schema examples include War but not true trick-taking games (Hearts, Spades, Bridge). Missing patterns:
- **Lead/follow suit obligations** (ConditionType.MUST_FOLLOW_SUIT?)
- **Trump hierarchy** (dynamic rank ordering)
- **Trick collection and scoring** (who wins the trick?)
- **Round vs trick vs hand distinction**

This is a significant gap—trick-taking represents a major card game family.

**2. Player Targeting Ambiguity**
`Location.OPPONENT_HAND` and `ActionType.DRAW_FROM_OPPONENT` don't specify:
- How to select which opponent in 3+ player games
- Whether selection is player choice, random, or deterministic
- Turn order and direction (clockwise? variable?)

For Go Fish specifically: "Ask any player for a rank" requires player selection + rank specification, which seems underspecified.

**3. Hidden Information Modeling**
The schema doesn't explicitly model:
- What information each player can see
- Face-up vs face-down cards in shared locations
- Revealed hands (showdown in poker)

This could cause issues where genomes specify actions requiring hidden knowledge (draw specific card from opponent's hand) or where validation can't distinguish legal from illegal information access.

**4. Termination Guarantee Gaps**
The question asks about validation for infinite loop prevention, but I don't see:
- Maximum turn limits
- Guaranteed deck depletion mechanics
- Progress invariants (some card must leave the system each turn)

A genome could specify: "Draw 1, discard 0, no win condition triggers" → infinite game. Validation rules need explicit termination criteria.

**5. Scoring/Win Condition Interaction**
`ScoringRules` and `WinConditions` are listed but their interaction isn't clear:
- Is scoring incremental or end-game?
- Can win conditions reference cumulative scores?
- How do multi-round games aggregate?

Games like Gin Rummy score across multiple hands—this temporal structure may be missing.

**6. Wildcards and Special Cards**
No mention of:
- Wild card designation (which cards are wild?)
- Special card effects beyond phase constraints
- Card transformation (wild becomes specified rank)

Crazy 8s is included but the "8 is wild/changes suit" mechanic's representation isn't clear.

**7. Evolution Friendliness Tension**
The schema may be **too structured** for open-ended evolution:
- Fixed phase order (Draw→Play→Discard) may be too rigid
- Can GA discover that some games skip phases?
- Can GA create new phase types or just parameterize existing ones?

Alternatively, it may be **too flexible** in places:
- Arbitrary condition composition could create nonsensical predicates
- Need fitness landscapes that guide toward playable games

---

## Recommendations

1. **Add Trick-Taking Extension**: `TrickPhase`, `ConditionType.MUST_FOLLOW_SUIT`, `ActionType.WIN_TRICK`, trump designation in `SetupRules`.

2. **Parameterize Player Targeting**: Add `TargetSelector` (OPPONENT_LEFT, OPPONENT_CHOICE, OPPONENT_RANDOM, ALL_OPPONENTS) to opponent-directed actions.

3. **Add Termination Constraints**: Require `max_turns` field or prove progress invariant (deck monotonically decreases, or chip total conserved with bankruptcy elimination).

4. **Clarify Wildcard Representation**: Add `SetupRules.wild_cards: List[Rank]` and `ConditionType.MATCHES_OR_WILD`.

5. **Add Visibility Annotations**: `Location` should pair with visibility (FACE_UP, FACE_DOWN, OWNER_ONLY).

---

## Assessment: 85-90% Coverage Claim

**Verdict: Optimistic but plausible with caveats**

For 2-player games without trick-taking, auction, or complex tableau mechanics, 85% seems achievable. The extensions address real gaps (betting, bluffing, collections).

However, the claim becomes strained when considering:
- Trick-taking family (Hearts, Spades, Bridge, Euchre) - major omission
- Partnership games
- Auction mechanics (Bridge bidding)
- Tableau/layout games (Solitaire variants, Canasta melds on table)

**Revised estimate**: 70-75% of Hoyle's simple card games without trick-taking extension; potentially 85% with it added.

The schema is a solid foundation but needs trick-taking, clearer multiplayer semantics, and termination guarantees to hit the high end of the coverage claim.

---

### Gemini Review

Here is my independent analysis of the extended genome schema.

## Strong Points

*   **Modular Extensibility via Enums:** The use of `Enum` types (ConditionType, ActionType) as the primary drivers for logic allows for genetic mutation without breaking the structural integrity of the schema. Adding `HAS_SET_OF_N` or `BET` as atomic opcodes makes the search space navigable for a Genetic Algorithm (GA).
*   **High-Value Genre Capture:** The specific extensions chosen (Set Detection, Betting, Bluffing) effectively unlock three massive branches of the card game evolutionary tree: Rummy-style games, Poker-style games, and Social Deduction games. This provides the highest ROI for schema complexity.
*   **Semantic "Sugar" for Complex Logic:** Abstracting complex logic like `HAS_RUN_OF_N` into a single `ConditionType` is excellent for evolution. If the GA had to construct a run detector using raw comparison operators (e.g., `card[i].val == card[i+1].val - 1`), it would likely never converge. High-level primitives act as "functional genes."
*   **Strict Phase Separation:** Distinguishing between `BettingPhase`, `ClaimPhase`, and `PlayPhase` simplifies the simulation engine's state machine. It prevents invalid states (e.g., trying to bet while resolving a discard) and makes validation easier.

## Moderate Points

*   **Evolutionary "Fragility" of Composite Conditions:** While individual conditions are robust, deeply nested logic (e.g., complex predicates combining `AND`, `OR`, and specific game states) can be brittle during crossover operations. The schema relies on the GA maintaining valid references within these trees.
*   **Bytecode Complexity for Set Detection:** While `HAS_SET_OF_N` is clean in the schema, implementing it efficiently in a Go-based VM requires the bytecode interpreter to perform computationally expensive combinatorial checks on every frame/action validation. This could impact the speed of training generations.
*   **Multiplayer Targeting Ambiguity:** The schema introduces `OPPONENT_HAND`, but in a game with 3+ players, "Opponent" is ambiguous. Does this default to the "Next" player? Does it require a targeting selector (e.g., "Choose Player")? The current schema implies a 2-player focus or a rigid rotation (Left/Right) which might limit the evolution of games like "Go Fish" (where you ask a specific person).
*   **Validation of Terminators:** With the introduction of Betting and Bluffing loops (Raise/Re-raise cycles, Challenge/Counter-challenge), guaranteeing game termination becomes significantly harder than in draw-and-discard games. The validator needs heuristic loop detection.

## Weak Points / Concerns

*   **Trick-Taking Structural Gap:** The extensions cover Rummy and Poker well, but **Trick-Taking** (Whist, Hearts, Spades, Bridge) appears to be a significant gap. These games require:
    *   A concept of a "Trick" (a sub-cycle within a round).
    *   "Lead suit" state tracking.
    *   "Must follow suit" constraints (which are different than standard play conditions).
    *   "Trump" logic (comparing cards based on a global variable).
    *   *Assessment:* Without specific structures for Tricks, the 85-90% coverage claim is likely overstated, as trick-taking is a pillar of Hoyle's games.
*   **Redundancy in Action/Location Logic:** There is a conflict between `ActionType.DRAW_FROM_OPPONENT` and the generic `ActionType.DRAW` combined with `Location.OPPONENT_HAND`.
    *   *Risk:* If the GA evolves `Action: DRAW, Source: OPPONENT_HAND` and also `Action: DRAW_FROM_OPPONENT, Source: DECK`, the engine semantics become muddy. Redundancy increases the search space unnecessarily.
*   **Melding Mechanics:** The schema supports *detecting* sets (`HAS_SET_OF_N`), and *discarding* pairs (`DISCARD_PAIRS`), but it is unclear if it supports **Melding** (playing a set to the table to score points without discarding). Rummy requires cards to move from Hand -> Table (Melds) separately from Hand -> Discard. The current `PlayPhase` actions need to explicitly support `PLAY_MELD`.
*   **State Vector Complexity for AI:** Introducing Betting (Chips) and Bluffing (Hidden Information) drastically changes the input requirements for the AI playing the game. The AI agents used for fitness evaluation must be upgraded to handle "Pot Odds" and "Probability of Lie," or the evolved games will be evaluated poorly.

## Recommendations

1.  **Add "Trick" Infrastructure:** Introduce a `TrickPhase` or a `trick_rules` definition that handles leading, following suit, and trump comparisons. This is essential for coverage.
2.  **Consolidate Actions:** Deprecate specific actions like `DRAW_FROM_OPPONENT` in favor of a robust `Action(Type=MOVE, Source=OPPONENT, Dest=SELF)` system to reduce redundancy.
3.  **Explicit "Tableau/Meld" Support:** Ensure there is a `Location.TABLEAU` or `Location.MELDS` where sets can be played and scored, distinct from the Discard pile.
4.  **Targeting Selector:** Add a `Target` field to actions (Target: NEXT_PLAYER, PREV_PLAYER, ANY_OPPONENT) to support 3+ player games effectively.

## Coverage Assessment
**Can this schema achieve 85-90% coverage?**
**No.**
It likely achieves **85-90% coverage of "Draw-and-Discard" and "Comparison" games** (War, Rummy, Poker variants).
However, it currently misses the structural logic for **Trick-Taking games** (Hearts, Spades) and **Solitaire-style Tableau games** (which require spatial column movement).
*Revised Estimate:* **60-70%** of Hoyle's games, rising to **80%** if Trick-Taking logic is added.

---

### Codex Review

## Strong Points
- The core is decomposed into enumerated types, condition/action systems, and phase structures, which gives the evolutionary process a constrained, composable search space rather than arbitrary code generation.
- The optional extensions cover several important families beyond the base schema: opponent interaction (Old Maid, Go Fish), set/run detection (Rummy variants), betting/wagering (poker‑like), and claim/challenge mechanics (Cheat/BS), which are common in “simple” card games.
- Making extensions optional means existing definitions for War, Crazy 8s, etc., can remain unchanged; you aren’t forcing all games to include betting or bluffing to compile.
- The fact that you’ve already encoded a diverse set of examples (Crazy 8s, War, Gin Rummy, Old Maid, Go Fish, Betting War, I Doubt It) suggests the schema is expressive enough for a broad swath of casual games.
- A dataclass‑based approach maps cleanly to both Python and a bytecode compiler for Go; the finite set of enums/actions should be straightforward to implement on both sides.

## Moderate Points
- The genetic algorithm friendliness is mixed: having well‑typed fields and enums curbs nonsensical mutations, but the number of knobs (phases, conditions, resource rules, claim types) could make the search space very large. You might need mutation/crossover operators tailored to this structure.
- Backward compatibility appears intact in spirit, but it hinges on defaulting extension fields sensibly. If defaults are not set, older genomes could break when the class definitions change.
- The semantics of some new elements are intuitive (e.g. `Location.OPPONENT_HAND` implies drawing from an opponent), but might need formal clarification when multiple opponents exist, or when multiple hands/discards are in play.
- Validation/termination checking is feasible for phase‑based turn structures, but bluffing/betting loops can introduce unbounded interaction unless you add explicit turn/round caps in the schema.
- Implementation of set/run detection and betting logic is conceptually simple, but edge cases (runs across suits, Aces high/low, minimum raises, all‑in side pots) could complicate the Go bytecode core if not nailed down.

## Weak Points / Concerns
- The claimed 85–90 % coverage is optimistic if you consider “simple card games” to include trick‑taking and partnership games (Hearts, Spades, Euchre, Bridge), bidding/contract phases, trump suits, trick evaluation, and role asymmetry – none of which are addressed in the current schema.
- Other common patterns not covered include: UNO‑style special card effects (skip, reverse, wild), shedding/accumulation games with action chaining, tableau/stack movement (solitaire), multi‑deck/shoe games, wild cards/jokers, and team scoring. Supporting these would require new condition/action primitives beyond the current extensions.
- Ambiguities remain around opponent interaction: does `DRAW_FROM_OPPONENT` choose a random opponent in multiplayer? Can you target specific opponents? Is `OPPONENT_HAND` a single location or per‑opponent?
- Resource/betting rules are sketched, but without a clear model for turn order in betting rounds, side pots, or what constitutes a round’s end. Those are non‑trivial to model and validate for termination.
- The schema doesn’t explicitly guard against infinite loops (e.g. a ClaimPhase that can be challenged ad infinitum or betting rounds with no cap). Without explicit iteration limits in the genome, playability validation will be hard.
- While dataclasses are convenient, the rigidity of predefined phases/actions could make it hard for evolution to discover novel mechanics; conversely, if fields like “post_deal_actions” accept arbitrary lists, the space may be too unconstrained without strong validation.

---

## Stage 2: Chairman Consensus (by Claude)

# Consensus Synthesis: Extended Genome Schema Analysis

## Areas of Agreement

**All three analyses converge on these points:**

1. **Trick-Taking is a Critical Gap** (UNANIMOUS, HIGH CONFIDENCE)
   - All reviewers independently identified the absence of trick-taking mechanics (Hearts, Spades, Bridge, Euchre) as the most significant omission
   - Missing: lead/follow suit obligations, trump hierarchy, trick collection/evaluation
   - This single gap substantially undermines the 85-90% coverage claim

2. **Core Architecture is Sound** (UNANIMOUS)
   - Enum-based condition/action system praised by all
   - Phase separation (Draw→Play→Discard) provides clear structure
   - Dataclass approach maps well to both Python and Go bytecode
   - Composable predicates enable meaningful genetic operations

3. **Backward Compatibility Appears Intact** (UNANIMOUS)
   - Extensions are genuinely optional
   - Base games (War, Crazy 8s) should work without modification
   - *Caveat from Codex*: depends on sensible defaults for new fields

4. **Multiplayer/Opponent Targeting is Ambiguous** (UNANIMOUS)
   - `OPPONENT_HAND` and `DRAW_FROM_OPPONENT` don't specify which opponent in 3+ player games
   - Go Fish requires asking *specific* players for ranks—currently underspecified
   - All recommend adding targeting selector (NEXT_PLAYER, CHOICE, RANDOM, etc.)

5. **Termination Guarantees are Missing** (UNANIMOUS)
   - No explicit max_turns or progress invariants
   - Betting/bluffing loops could run infinitely
   - Validation for playability needs explicit termination criteria

6. **Coverage Claim is Overstated** (UNANIMOUS)
   - 85-90% is optimistic without trick-taking support
   - Claude estimates 70-75% without trick-taking extension
   - Gemini estimates 60-70% current, 80% with trick-taking
   - Codex: "optimistic" for anything beyond draw-and-discard games

## Areas of Disagreement

### 1. Severity of Melding Gap
| Reviewer | Position |
|----------|----------|
| **Gemini** | Explicit concern: Schema detects sets but lacks PLAY_MELD action for Rummy-style tableau play. Recommends `Location.TABLEAU` or `Location.MELDS` |
| **Claude** | Mentions "tableau building" as unrepresentable but doesn't emphasize melding specifically |
| **Codex** | Lists "tableau/stack movement" as missing but focuses on solitaire rather than Rummy melds |

**Conflict**: Gemini sees melding as a blocking issue for Rummy coverage; others treat it as lower priority.

### 2. Action Redundancy Concerns
| Reviewer | Position |
|----------|----------|
| **Gemini** | Flags explicit conflict: `DRAW_FROM_OPPONENT` vs generic `DRAW` + `Location.OPPONENT_HAND` creates semantic ambiguity. Recommends consolidating to `Action(Type=MOVE, Source=OPPONENT, Dest=SELF)` |
| **Claude** | Does not mention this redundancy |
| **Codex** | Does not mention this redundancy |

**Conflict**: Only Gemini identifies action/location redundancy as a design problem.

### 3. Evolution Friendliness Assessment
| Reviewer | Position |
|----------|----------|
| **Claude** | Schema may be "too structured" (fixed phase order) OR "too flexible" (arbitrary condition composition). Tension exists. |
| **Gemini** | Praises enum-based mutations but warns composite conditions are "brittle during crossover" |
| **Codex** | "Mixed"—well-typed fields curb nonsense but large knob count expands search space. Needs tailored mutation operators. |

**Nuance**: All see trade-offs, but Claude emphasizes structural rigidity while Gemini/Codex focus on search space complexity.

### 4. Wildcard/Special Card Treatment
| Reviewer | Position |
|----------|----------|
| **Claude** | Explicitly recommends `SetupRules.wild_cards: List[Rank]` and `ConditionType.MATCHES_OR_WILD` |
| **Codex** | Lists "wild cards/jokers" as missing pattern but no specific recommendation |
| **Gemini** | Does not mention wildcards |

**Conflict**: Varying levels of concern about wildcard representation.

### 5. Bytecode Implementation Complexity
| Reviewer | Position |
|----------|----------|
| **Claude** | "Moderate confidence—workable but careful implementation needed" for condition evaluation |
| **Gemini** | Warns `HAS_SET_OF_N` requires expensive combinatorial checks in Go VM, could impact training speed |
| **Codex** | "Conceptually simple" but edge cases (Aces high/low, side pots) could complicate Go core |

**Conflict**: Gemini most concerned about performance; Claude and Codex focus on correctness edge cases.

### 6. Revised Coverage Estimates
| Reviewer | Current | With Trick-Taking |
|----------|---------|-------------------|
| **Claude** | 70-75% | ~85% |
| **Gemini** | 60-70% | ~80% |
| **Codex** | Not quantified but "optimistic" | Not quantified |

**Conflict**: 10-15 percentage point spread in current coverage estimates.

## Confidence Level

**MEDIUM-HIGH**

- High confidence on architectural strengths and trick-taking gap
- Medium confidence on coverage percentages (definitional ambiguity around "simple card games")
- Medium confidence on evolution friendliness (requires empirical testing)

## Synthesized Recommendation

### Must Address (Consensus Critical)

1. **Add Trick-Taking Extension**
   - `TrickPhase` with lead/follow mechanics
   - `ConditionType.MUST_FOLLOW_SUIT`
   - Trump designation in `SetupRules`
   - Trick evaluation/collection actions
   - *Without this, coverage claim should be revised down to 65-75%*

2. **Add Player Targeting System**
   - `TargetSelector` enum: `NEXT_PLAYER`, `PREV_PLAYER`, `CHOICE`, `RANDOM`, `ALL_OPPONENTS`
   - Attach to opponent-directed actions
   - Required for accurate Go Fish, 3+ player games

3. **Add Termination Guarantees**
   - `max_turns` field OR
   - Provable progress invariant (deck depletion, chip conservation with elimination)
   - Essential for validation

### Should Address (Majority Concern)

4. **Clarify Visibility/Hidden Information**
   - Add `FACE_UP`, `FACE_DOWN`, `OWNER_ONLY` annotations to locations
   - Prevents illegal information access in genomes
   - *Claude and Codex flagged; Gemini implied via AI state vector concerns*

5. **Add Wildcard Support**
   - `SetupRules.wild_cards: List[Rank]`
   - `ConditionType.MATCHES_OR_WILD`
   - *Claude explicit; Codex mentioned; Gemini silent*

6. **Resolve Melding Ambiguity** (Gemini's concern)
   - Clarify whether `Location.TABLE` supports per-player melds
   - Add `PLAY_MELD` action if needed for Gin Rummy scoring
   - *May already be representable—needs documentation*

### Consider Addressing (Single Reviewer or Lower Priority)

7. **Consolidate Action Redundancy** (Gemini only)
   - Evaluate whether `DRAW_FROM_OPPONENT` should be deprecated in favor of generic `DRAW` + source location
   - Trade-off: specificity vs. search space size

8. **Parameterize Set/Run Detection** (Claude only)
   - Same-suit requirement for runs?
   - Location scope (hand vs. table)?
   - May be documentation issue rather than schema change

### Coverage Assessment Revision

**Synthesized estimate:**
- **Current schema**: 65-75% of Hoyle's simple card games
- **With trick-taking extension**: 80-85%
- **With all recommended additions**: Potentially 85-90%

The original 85-90% claim is achievable but requires the trick-taking extension at minimum. The schema is a solid foundation that needs one major addition (tricks) and several clarifications to meet its stated goals.

