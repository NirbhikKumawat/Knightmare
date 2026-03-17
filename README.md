# GoChess ♟️

A blazingly fast, from-scratch chess engine written entirely in Go.

**GoChess** is built on a highly optimized 64-bit **Bitboard** architecture. It is designed for maximum node throughput and strict mathematical accuracy, clocking in at over **5.5 Million Nodes Per Second (NPS)** on a single thread during move generation and legality verification.

## ⚡ Performance: The Perft Benchmark

The engine's move generator has been rigorously tested against the standard Perft (Performance Test) benchmarks from the starting position. It successfully traverses all **119,060,324 leaf nodes at Depth 6 in ~21 seconds** (single-threaded), proving 100% compliance with the rules of chess (including complex edge cases like en passant discovered checks, castling rights destruction, and promotions).

| Depth | Nodes       | Accuracy |
|:------|:------------|:---------|
| 1     | 20          | ✅ Pass   |
| 2     | 400         | ✅ Pass   |
| 3     | 8,902       | ✅ Pass   |
| 4     | 197,281     | ✅ Pass   |
| 5     | 4,865,609   | ✅ Pass   |
| 6     | 119,060,324 | ✅ Pass   |

## 🧠 Core Architecture & Technical Highlights

This engine avoids slow arrays and loops in favor of raw bitwise arithmetic and pre-calculated lookup tables.

* **Bitboards:** The board is represented by arrays of `uint64` integers. Piece movement, captures, and raycasting are resolved in fractions of a nanosecond using bitwise `AND`, `OR`, `XOR`, and bit-shifting.
* **Custom Magic Bitboards:** Sliding pieces (Rooks, Bishops, Queens) use generated Magic Bitboards to instantly look up attack rays. The engine includes a custom sparse-random brute-force generator to perfectly map blocker permutations to array indices, entirely bypassing expensive on-the-fly ray calculations.
* **Pseudo-Legal Move Generation:** Moves are generated in bulk (e.g., shifting entire pawn bitboards at once) rather than square-by-square.
* **Reverse Attack Legality Checking:** To verify King safety, the engine utilizes the "Super Piece" concept—projecting attacks outward from the King's square to detect overlapping enemy pieces, ensuring blistering fast legality checks.
* **Copy-Make Paradigm:** Taking advantage of Go's highly efficient struct copying, the engine uses a `MakeMove` approach on copied board states rather than a cumbersome `UnmakeMove` function, keeping the state mutation clean and fast.

## 📦 Features

Beyond the core rules, the engine is fully equipped to interact with the standard chess data ecosystem:
* **FEN Support:** Flawless parsing and generation of Forsyth–Edwards Notation strings.
* **SAN Disambiguation:** Complete implementation of Standard Algebraic Notation (e.g., `Nbd2`, `exd8=Q#`). The engine uses a generator-matching algorithm to perfectly resolve complex disambiguation conflicts.
* **PGN Parsing:** A custom, regex-free PGN parser capable of ingesting entire historical grandmaster games, extracting metadata headers, and mapping the move sequence to internal bitboard states.

## 🚀 Getting Started

Ensure you have Go installed (1.18+ recommended).

### Installation
`git clone https://github.com/NirbhikKumawat/GoChess.git`
`cd GoChess`

### Running the Tests
To verify the move generator's accuracy and benchmark the speed on your hardware:
`go test -v ./chess`

## 🛣️ Roadmap
- [x] Bitboard Representation
- [x] Magic Bitboard Generation
- [x] Pseudo-Legal / Legal Move Generation
- [x] FEN / SAN / PGN I/O Layer
- [ ] Static Evaluation (Piece-Square Tables, Material Weights)
- [ ] Minimax Search with Alpha-Beta Pruning
- [ ] Zobrist Hashing & Transposition Tables
- [ ] UCI (Universal Chess Interface) Protocol Support


---
Made by [NirbhikTheNice](https://github.com/NirbhikKumawat/GoChess.git)