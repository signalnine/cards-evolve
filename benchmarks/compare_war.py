"""Compare Python vs Golang War game performance."""

import time
import subprocess
import statistics
from cards_evolve.simulation.war import play_war_game


def benchmark_python_war(iterations: int = 100) -> float:
    """Benchmark Python War implementation."""
    times = []

    for i in range(iterations):
        start = time.perf_counter()
        play_war_game(seed=i, max_turns=1000)
        elapsed = time.perf_counter() - start
        times.append(elapsed)

    return statistics.mean(times)


def benchmark_golang_war(iterations: int = 100) -> float:
    """Benchmark Golang War implementation via subprocess."""
    # Run Go benchmark and parse output
    result = subprocess.run(
        ["go", "test", "./game", "-bench=BenchmarkPlayWarGame",
         f"-benchtime={iterations}x"],
        cwd="src/gosim",
        capture_output=True,
        text=True
    )

    # Parse: "BenchmarkPlayWarGame-8   100   12345 ns/op"
    for line in result.stdout.split('\n'):
        if 'BenchmarkPlayWarGame' in line:
            parts = line.split()
            ns_per_op = float(parts[-2])
            return ns_per_op / 1e9  # Convert ns to seconds

    return 0.0


def main() -> None:
    """Run performance comparison."""
    print("Performance Comparison: Python vs Golang War Game")
    print("=" * 60)

    iterations = 100
    print(f"\nRunning {iterations} iterations of War game (max 1000 turns)...\n")

    print("Python implementation:")
    python_time = benchmark_python_war(iterations)
    print(f"  Average time: {python_time*1000:.2f}ms per game")

    print("\nGolang implementation:")
    golang_time = benchmark_golang_war(iterations)
    print(f"  Average time: {golang_time*1000:.2f}ms per game")

    speedup = python_time / golang_time
    print(f"\nSpeedup: {speedup:.1f}x")

    if speedup >= 10:
        print("✅ SUCCESS: Golang is 10x+ faster")
    else:
        print(f"⚠️  WARNING: Speedup is only {speedup:.1f}x (target: 10x+)")


if __name__ == "__main__":
    main()
