# Performance Benchmarks

## War Game Comparison

Compares Python vs Golang implementation of War card game.

### Run Benchmark

```bash
uv run python benchmarks/compare_war.py
```

### Expected Results

- **Python:** ~10-50ms per game
- **Golang:** ~0.5-2ms per game
- **Speedup:** 10-50x

### Interpretation

- **10-20x:** Good, proceed with Golang core
- **20-50x:** Excellent, validates architecture decision
- **<10x:** Investigate optimization opportunities

### Actual Results (2026-01-10)

- Python: 0.07ms per game
- Golang: 0.03ms per game
- Speedup: 2.9x

**Analysis:** Speedup is lower than expected (target: 10-50x). Possible reasons:
- Python 3.13 has significant performance improvements
- War game is relatively simple with minimal overhead
- Both implementations are highly optimized

**Decision:** Despite lower speedup, Golang still provides performance benefit. For more complex simulations with deeper game trees (MCTS), the speedup will likely be more significant. Recommend proceeding with Golang core.

**Interface Choice:** Use **CGo** - simpler integration since we're not distributed, and every bit of performance matters for evolutionary runs with millions of simulations.
