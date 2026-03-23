# Module 02 — PN-Counter
## Positive-Negative Counter: Adding Decrements

**Concept:** CRDT composition, negative operations, monotonicity preservation  
**Challenge:** Extend G-Counter to support decrement without breaking convergence  
**Time:** ~2 hours

---

## Introduction

In Module 01, you built a G-Counter that only grows. But what if you need a counter that can both increment **and** decrement?

The naive approach—storing a single integer and allowing both `+` and `-` operations—breaks down in distributed systems. You can't guarantee convergence when operations can cancel each other out.

The **PN-Counter** (Positive-Negative Counter) solves this elegantly by **composing two G-Counters**: one for increments (P), one for decrements (N). The value is simply `P - N`.

This module teaches you the most important CRDT pattern: **build complex CRDTs by composing simpler ones**.

---

## The Problem

Imagine three replicas tracking a bank account balance:

```
Replica A: balance = 100
Replica B: balance = 100
Replica C: balance = 100
```

Now:
- Replica A: `+50` → balance = 150
- Replica B: `-30` → balance = 70
- Replica C: `+20` → balance = 120

If they sync by just sending deltas, the order matters:
- A syncs to B first: 70 + 50 = 120, then C: 120 + 20 = 140 ✅
- C syncs to B first: 70 + 20 = 90, then A: 90 + 50 = 140 ✅

That looks fine! But what about this scenario with **out-of-order operations**:

```
Starting: both replicas at balance = 100

Replica A does: +50, then -30 → final balance = 120
Replica B does: -30, then +50 → final balance = 120
```

If they send **final balances** (120 and 120), they might converge. But what if they send **operations** that arrive out of order or get duplicated?

**Scenario 1:** B receives A's operations in wrong order
```
B starts at: 100
B applies: -30 (its own) → 70
B receives: -30 from A (arrived first!) → 40
B receives: +50 from A → 90  ❌ Wrong! Should be 120
```

**Scenario 2:** Operations get duplicated (common in unreliable networks)
```
A starts at: 100
A applies: +50, -30 → 120
A receives: -30 from B → 90
A receives: -30 from B AGAIN (duplicate!) → 60  ❌ Wrong!
```
```
A starts at: 100
A applies: +50 → 150
A applies: -30 → 120
A receives B's -30 → 90
A receives B's -30 again (duplicate message) → 60  ❌ Wrong!
```

**Scenario 3:** Total chaos
```
Each replica applies operations in different orders, gets different results:
- Replica A: 120
- Replica B: 90
- Network duplicate causes C: 60
❌ No convergence!
```

**The PN-Counter solves this** by tracking **how many times** each replica has incremented/decremented, not the operations themselves. Order doesn't matter because merge uses `max` (idempotent and commutative).

---

## How It Works

Instead of storing a single counter, store **two G-Counters**:

```go
type PNCounter struct {
    replicaID string
    p         *gcounter.GCounter  // Positive (increments)
    n         *gcounter.GCounter  // Negative (decrements)
}
```

### Operations

1. **Increment**: Increase the P counter
   ```go
   func (pn *PNCounter) Increment() {
       pn.p.Increment()
   }
   ```

2. **Decrement**: Increase the N counter (yes, *increase*!)
   ```go
   func (pn *PNCounter) Decrement() {
       pn.n.Increment()
   }
   ```

3. **Value**: Difference between P and N
   ```go
   func (pn *PNCounter) Value() int {
       return pn.p.Value() - pn.n.Value()
   }
   ```

4. **Merge**: Merge both P and N counters
   ```go
   func (pn *PNCounter) Merge(other *PNCounter) {
       pn.p.Merge(other.p)
       pn.n.Merge(other.n)
   }
   ```

### Example

```
Replica A increments 5 times:
  P = {"A": 5}, N = {"A": 0}
  Value() = 5 - 0 = 5

Replica B decrements 3 times:
  P = {"B": 0}, N = {"B": 3}
  Value() = 0 - 3 = -3

Merge A into B:
  P = {"A": 5, "B": 0}, N = {"A": 0, "B": 3}
  Value() = 5 - 3 = 2 ✅

Replica C increments 2 times, decrements 1 time:
  P = {"C": 2}, N = {"C": 1}
  Value() = 2 - 1 = 1

Merge C into B:
  P = {"A": 5, "B": 0, "C": 2}, N = {"A": 0, "B": 3, "C": 1}
  Value() = 7 - 4 = 3 ✅
```

No matter what order you merge, you always get the same result because both P and N are G-Counters with guaranteed convergence.

---

## Why This Works

The key insight: **decrement is just an increment of a different counter**.

Since both P and N are G-Counters (which we proved converge in Module 01), and subtraction is deterministic, the PN-Counter inherits convergence:

- **Commutativity:** Merging P and N in any order gives the same result
- **Idempotence:** Merging the same state twice doesn't change anything
- **Associativity:** Grouping doesn't matter

This is **CRDT composition**: combine proven-convergent CRDTs to build more powerful ones.

---

## Getting Started

Copy the starter templates:

```bash
# From the module-02-pncounter directory
cp ../starters/module-02-pncounter/*.tmpl .
for f in *.tmpl; do mv "$f" "${f%.tmpl}"; done
```

This gives you:
- `pncounter.go` — struct definition and function stubs
- `pncounter_test.go` — test suite with TODOs

**Why .tmpl?** Template files let you reset to the starting state anytime. They're ignored by the Go language server, preventing IDE conflicts with the solution files.

### Need to Reset?

Made a mistake and want to start over? Just re-copy the templates:

```bash
rm pncounter.go pncounter_test.go
cp ../starters/module-02-pncounter/*.tmpl .
for f in *.tmpl; do mv "$f" "${f%.tmpl}"; done
```

---

## The Challenge

Implement `PNCounter` with the following API:

```go
type PNCounter struct {
    replicaID string
    p         *gcounter.GCounter  // Positive counter
    n         *gcounter.GCounter  // Negative counter
}

// NewPNCounter creates a new PN-Counter for the given replica
func NewPNCounter(replicaID string) *PNCounter

// Increment increases the value by 1
func (pn *PNCounter) Increment()

// Decrement decreases the value by 1
func (pn *PNCounter) Decrement()

// Value returns the current value (P - N)
func (pn *PNCounter) Value() int

// Merge combines another PN-Counter's state into this one
func (pn *PNCounter) Merge(other *PNCounter)
```

### Requirements

1. **Reuse Module 01**: Import and use your `GCounter` from `module-01-gcounter`
2. **Composition**: PN-Counter should contain two G-Counters (P and N)
3. **Increment** only affects the P counter
4. **Decrement** only affects the N counter
5. **Value** returns `P.Value() - N.Value()`
6. **Merge** merges both P and N
7. **Convergence:** All replicas that have seen the same operations should have the same `Value()`

---

## Workflow

### Step 1: Set Up Your Files

Copy the starter templates (if you haven't already):

```bash
cp ../starters/module-02-pncounter/*.tmpl .
for f in *.tmpl; do mv "$f" "${f%.tmpl}"; done
```

### Step 2: Add Package Declaration

Both files need:
```go
package pncounter
```

And `pncounter.go` needs an import:
```go
import "github.com/zshainsky/crdt-tutorial/module-01-gcounter"
```

### Step 3: Write Your Tests

Open `pncounter_test.go` and fill in the test bodies. Each test has a comment describing what to assert.

```go
// Should return 0 for a new counter
func TestPNCounterInitialValue(t *testing.T) {
    // TODO: Create a PNCounter and assert Value() == 0
}
```

### Step 4: Implement

Edit `pncounter.go` and implement the five functions. Remember:
- Use `gcounter.NewGCounter()` to create the P and N counters
- Delegate to the underlying G-Counter methods

### Step 5: Run Your Tests

```bash
go test .
```

Fix until all tests pass.

### Step 6: Compare with Solution

Want to check your work against a reference implementation?

```bash
# View the solution
cat ../solutions/module-02-pncounter/pncounter_solution.go

# Run solution tests
cd ../solutions/module-02-pncounter
go test .
cd ../../module-02-pncounter
```

All green = module complete ✅

### Step 7: Stuck?

- **Check hints below** (expand as needed)
- **Review Module 01** — make sure you understand G-Counter
- **Review the solution** in `../solutions/module-02-pncounter/`
- **Ask questions** — composition can be tricky at first!

---

## Hints

<details>
<summary><strong>Hint 1:</strong> How do I import GCounter?</summary>

At the top of `pncounter.go`:

```go
import "github.com/zshainsky/crdt-tutorial/module-01-gcounter"
```

Then use `gcounter.NewGCounter(replicaID)` to create instances.

</details>

<details>
<summary><strong>Hint 2:</strong> How do I initialize P and N?</summary>

In `NewPNCounter`:

```go
func NewPNCounter(replicaID string) *PNCounter {
    return &PNCounter{
        replicaID: replicaID,
        p:         gcounter.NewGCounter(replicaID),
        n:         gcounter.NewGCounter(replicaID),
    }
}
```

Both P and N use the same `replicaID` because they're part of the same logical counter.

</details>

<details>
<summary><strong>Hint 3:</strong> Wait, you increment N to decrement?</summary>

Yes! Think about it:
- P tracks "how many times I've incremented"
- N tracks "how many times I've decremented"
- Value = P - N

So when you decrement, you're adding to the "decrement counter" (N), not subtracting from anything.

</details>

<details>
<summary><strong>Hint 4:</strong> My merge isn't working correctly</summary>

Make sure you merge **both** P and N:

```go
func (pn *PNCounter) Merge(other *PNCounter) {
    pn.p.Merge(other.p)
    pn.n.Merge(other.n)
}
```

If you only merge one, you lose half the state!

</details>

<details>
<summary><strong>Hint 5:</strong> Can Value() return negative numbers?</summary>

Absolutely! If you decrement more than you increment, N will be larger than P, so `P - N` is negative. That's expected behavior.

</details>

---

## Questions to Ponder

- **Why can't we just use a single counter with negative increments?** Because G-Counter uses `max` to merge, which assumes monotonic growth. A single counter with negatives would break convergence.

- **What if two replicas decrement at the same time?** No problem! Each replica's decrements are tracked independently in N (just like increments in P), and `max` merge ensures convergence.

- **Is this efficient for large numbers?** Yes! The vector size is O(replicas), not O(operations). Even after millions of increments/decrements, you only store one count per replica in both P and N.

- **What's the pattern here?** To add a "negative" or "inverse" operation to any grow-only CRDT, create a second instance and subtract. This pattern works for many CRDTs!

---

## Extension Challenge (Optional)

Once your tests pass, try adding:

1. **IncrementBy(n int)** — increment by an amount
2. **DecrementBy(n int)** — decrement by an amount
3. **Reset()** — semantically impossible for CRDTs! Why? (Hint: how would you represent "forget all history" in a distributed system?)

---

## Next Module

Once you've completed this module, move on to [Module 03 - LWW-Register & OR-Set](../module-03-lwwregister/) to learn about registers (single values) and sets with tombstones.

---

**Ready? Copy the templates and start implementing!**
