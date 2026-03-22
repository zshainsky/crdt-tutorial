# Build Your Own CRDTs
## Conflict-Free Replicated Data Types from First Principles

**Welcome!** This tutorial teaches you to build production-ready CRDTs (Conflict-Free Replicated Data Types) from scratch in Go, culminating in a real-time collaborative todo list with WebSocket sync.

---

## What You'll Build

By the end of this tutorial, you'll have:

1. **A Go CRDT library** (`github.com/zshainsky/crdt-tutorial/crdts`) with:
   - G-Counter (grow-only counter)
   - PN-Counter (increment/decrement)
   - LWW-Register (last-write-wins register)
   - OR-Set (observed-remove set)
   - TodoList CRDT (composition of the above)

2. **A real-time collaborative todo list service** where multiple clients can simultaneously add, edit, complete, and reorder tasks with guaranteed eventual consistency—no coordination required.

---

## Prerequisites

- **Go 1.21+** installed
- Comfortable with Go syntax and basic data structures
- No distributed systems experience required—you'll learn as you build

---

## Modules

| Module | Topic | Time |
|--------|-------|------|
| [01 - G-Counter](module-01-gcounter/) | Grow-only counter, vector clocks, commutativity | ~2 hours |
| [02 - PN-Counter](module-02-pncounter/) | Increment & decrement, CRDT composition | ~2 hours |
| [03 - Registers & Sets](module-03-registers-and-sets/) | LWW-Register, OR-Set, tombstones | ~2-3 hours |
| [04 - Todo List CRDT](module-04-todolist/) | Application-layer CRDT composition | ~2-3 hours |
| [05 - Real-Time Sync Service](module-05-sync-service/) | WebSocket server + browser client | ~3-4 hours |
| **[06 - Canvas Extension (Bonus)]** | **Upgrade to drawing canvas** | **~4-6 hours** |

**Total:** 10-15 hours for core modules

---

## How to Use This Tutorial

Each module follows the same structure:

### 1. Read the Concept (README.md)
Start in `module-XX-.../README.md` to understand the theory and challenge.

### 2. Write Your Tests (`starter_test.go`)
Each module provides test stubs. Fill in the assertion logic:
```go
// Should return 5 after incrementing 5 times
func TestGCounterIncrement(t *testing.T) {
    // TODO: Your assertion here
}
```

### 3. Implement (`*.go`)
Build your implementation until your tests pass.

### 4. Verify Against the Spec (`solution_test.go`)
When you think you're done:
```bash
cd module-XX-.../solution
cp ../your-implementation.go .
go test ./...
```

All green = module complete ✅

### 5. Stuck? Check the Solution
`solution/*.go` contains the reference implementation. Try first, peek when stuck.

---

## Learning Philosophy

- **Learn by building.** Every module produces working code.
- **Tests as spec.** `solution_test` is the ground truth—no ambiguity.
- **Scaffold complexity.** Concepts introduced just-in-time.
- **Respect the engineer.** We assume competence—no "what is a variable" fluff.

---

## Getting Started

```bash
cd module-01-gcounter
cat README.md  # Read the challenge
```

---

## Optional Extension: Drawing Canvas

After completing Module 05, you'll have all the patterns needed to build a collaborative drawing canvas (Module 06). Same sync architecture, different data model:

- Replace `Task` → `Shape`
- Replace HTML list → HTML5 Canvas
- Add RGA (Replicated Growable Array) for freehand paths

This is the stepping stone to apps like Figma, Miro, Excalidraw.

---

## Questions?

As you work through modules, feel free to ask questions. I'm here to:
- Give Socratic hints (nudge you toward the answer)
- Review your code
- Explain concepts deeper
- Provide the solution if you're truly stuck

**Don't open `solution/` until you've genuinely tried!** The learning happens when you struggle.

---

## Ready?

Start here: **[Module 01 - G-Counter](module-01-gcounter/)**

Good luck! 🚀
