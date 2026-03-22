# Module 01 — G-Counter
## The Simplest CRDT: Grow-Only Counter

**Concept:** Vector clocks, commutativity, idempotence, and convergence  
**Challenge:** Implement a distributed counter that only increments  
**Time:** ~2 hours

---

## Introduction

Welcome to your first CRDT! The **G-Counter** (Grow-Only Counter) is the simplest possible conflict-free replicated data type, but it teaches you the three essential properties that make *all* CRDTs work:

1. **Commutativity** — merging A into B, then C is the same as merging C into B, then A
2. **Idempotence** — merging the same state twice has no effect
3. **Associativity** — the order of merges doesn't matter

If your data structure has these properties, it will **always converge** to the same state across all replicas, no matter what order updates arrive in.

---

## The Problem

Imagine three servers (replicas) all tracking the same counter:

```
Replica A: counter = 0
Replica B: counter = 0  
Replica C: counter = 0
```

Now:
- Replica A increments 3 times → counter = 3
- Replica B increments 2 times → counter = 2
- Replica C increments 5 times → counter = 5

If they sync by just sending their counter value, they'll get different results depending on who syncs last (last-write-wins). The "correct" answer should be **10** (3 + 2 + 5), but a naive sync will give you either 3, 2, or 5.

**The G-Counter solves this.**

---

## How It Works

Instead of storing a single integer, each replica tracks increments from *every* replica in a **vector**:

```go
type GCounter struct {
    replicaID string
    counts    map[string]int  // replicaID → increment count
}
```

When Replica A increments, it updates `counts["A"]`. When replicas merge, they take the **max** of each entry:

```go
func (g *GCounter) Merge(other *GCounter) {
    for replicaID, count := range other.counts {
        if count > g.counts[replicaID] {
            g.counts[replicaID] = count
        }
    }
}
```

The total value is the **sum** of all counts:

```go
func (g *GCounter) Value() int {
    total := 0
    for _, count := range g.counts {
        total += count
    }
    return total
}
```

### Example

```
Replica A increments 3 times:
  counts = {"A": 3}

Replica B increments 2 times:
  counts = {"B": 2}

Merge A into B:
  counts = {"A": 3, "B": 2}
  Value() = 5 ✅

Replica C increments 5 times:
  counts = {"C": 5}

Merge C into B:
  counts = {"A": 3, "B": 2, "C": 5}
  Value() = 10 ✅
```

No matter what order you merge, you always get 10.

---

## Why This Works

- **Commutativity:** `Merge(A, Merge(B, C))` = `Merge(B, Merge(A, C))` because `max` is commutative
- **Idempotence:** Merging the same state twice doesn't change anything (`max(3, 3) = 3`)
- **Associativity:** Grouping doesn't matter: `Merge(Merge(A, B), C)` = `Merge(A, Merge(B, C))`

These three properties → **guaranteed convergence**.

---

## The Challenge

Implement `GCounter` with the following API:

```go
type GCounter struct {
    replicaID string
    counts    map[string]int
}

// NewGCounter creates a new G-Counter for the given replica
func NewGCounter(replicaID string) *GCounter

// Increment increases this replica's count by 1
func (g *GCounter) Increment()

// Value returns the sum of all replica counts
func (g *GCounter) Value() int

// Merge combines another G-Counter's state into this one
func (g *GCounter) Merge(other *GCounter)
```

### Requirements

1. **Increment** only affects the local replica's count
2. **Value** returns the sum across all replicas
3. **Merge** takes the max of each replica's count
4. **Convergence:** After any sequence of increments and merges, all replicas that have seen the same updates should have the same `Value()`

---

## Workflow

### Step 1: Write Your Tests
Open `starter_test.go` and fill in the test bodies. Each test has a comment describing what to assert.

```go
// Should return 0 for a new counter
func TestGCounterInitialValue(t *testing.T) {
    // TODO: Create a GCounter and assert Value() == 0
}
```

### Step 2: Implement
Edit `gcounter.go` and implement the four functions.

### Step 3: Run Your Tests
```bash
go test .
```

Fix until all tests pass.

### Step 4: Verify Against Spec
When you think you're done:
```bash
cd solution
cp ../gcounter.go .
go test ./...
```

All green = module complete ✅

### Step 5: Stuck?
- **Check hints below** (expand as needed)
- **Ask questions** — I'm here to help
- **Peek at `solution/gcounter.go`** (only if truly stuck)

---

## Hints

<details>
<summary><strong>Hint 1:</strong> How do I initialize the counts map?</summary>

In `NewGCounter`, use `make(map[string]int)` and optionally set `counts[replicaID] = 0`.

</details>

<details>
<summary><strong>Hint 2:</strong> What does Increment do exactly?</summary>

It increments `counts[g.replicaID]`. Since `Increment` is always called on a specific replica, it only affects that replica's entry.

</details>

<details>
<summary><strong>Hint 3:</strong> How do I handle missing keys in Merge?</summary>

Go maps return `0` for missing keys, so `g.counts[replicaID]` is safe even if the key doesn't exist yet.

</details>

<details>
<summary><strong>Hint 4:</strong> My merge isn't idempotent—help!</summary>

Make sure you're using `max(a, b)` not `a + b`. The insight: each replica's count is monotonic (only goes up), so taking the max ensures you don't double-count.

</details>

---

## Questions to Ponder (Q&A Seeds)

- **Why can't we just send deltas (e.g., "+3")?** Because replicas might receive deltas out of order or multiple times. The vector clock approach is idempotent and order-independent.
  
- **Why max instead of sum?** Sum would double-count when you merge the same state twice. Max ensures idempotence.

- **Can G-Counter decrement?** No! That's why it's "grow-only." Decrement breaks the monotonicity that makes `max` work. (Next module: PN-Counter solves this.)

- **What if two replicas use the same ID?** Their increments will interfere (one will overwrite the other). Replica IDs must be unique (use UUIDs in production).

---

## Next Module

Once you've completed this module, move on to [Module 02 - PN-Counter](../module-02-pncounter/) to learn how to handle decrements by composing two G-Counters.

---

**Ready? Start by filling in `starter_test.go`!**
