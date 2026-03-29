# Module 03 — Registers & Sets
## LWW-Register and OR-Set: Values and Membership

**Concepts:** Last-write-wins semantics, tombstones, add-wins conflict resolution
**Challenge:** Implement a single-value register and an add-wins set
**Time:** ~2-3 hours

---

## Introduction

Modules 01 and 02 taught you counters — CRDTs that only store aggregate numbers. But most real applications need to store *values* (like a user's name) and *collections* (like a set of tags). This module introduces two CRDTs that handle both.

**Part 1 — LWW-Register:** A register stores a single value. In a distributed system, two replicas might write different values simultaneously. The LWW-Register resolves this with a simple rule: the write with the highest timestamp wins.

**Part 2 — OR-Set:** A set supports add and remove. But remove is surprisingly hard in distributed systems. If replica A removes an element while replica B adds it concurrently, what should happen after merge? The OR-Set (Observed-Remove Set) answers: **the add wins**, using a technique called *unique tags*.

---

## Part 1: LWW-Register

### The Problem

Imagine two replicas storing a user's display name:

```
Replica A: name = "alice"   (written at t=100)
Replica B: name = "Alice"   (written at t=200, a correction)
```

When they sync, which value should survive? The LWW-Register says: **the one with the higher timestamp**. Replica A adopts "Alice" because t=200 > t=100.

Now what if they write at the exact same millisecond?

```
Replica A: name = "alice"   (t=100)
Replica B: name = "Alice"   (t=100)
```

We need a deterministic tie-breaking rule. In this implementation, the current value wins — merging doesn't change the register when timestamps are equal.

### How It Works

The register stores a value and a timestamp together:

```
state = { value: "alice", timestamp: 100 }
```

**Set:** Write a new value with the current time as its timestamp.

**Merge:** Compare timestamps. Keep whichever is higher. If equal, keep the current value.

### Example

```
Replica A writes "alice" at t=100:
  state = { value: "alice", ts: 100 }

Replica B writes "Alice" at t=200:
  state = { value: "Alice", ts: 200 }

A merges B:
  200 > 100 → adopt B's state
  A.state = { value: "Alice", ts: 200 } ✅

B merges A:
  200 > 100 → keep current state (B's)
  B.state = { value: "Alice", ts: 200 } ✅

Both converge to "Alice". ✅
```

### Tradeoffs

LWW is simple and efficient, but it relies on synchronized clocks. Two replicas writing at exactly the same logical time require a tie-breaking rule that is necessarily arbitrary. For many use cases (display names, settings, preferences) this is fine. For financial transactions or conflict-sensitive data, you'd want a different strategy.

---

## Part 2: OR-Set

### The Problem

Suppose replica A removes an element from a set at the same moment replica B adds it back. What's the right answer after merge?

```
Replica A: removes "apple"
Replica B: adds "apple"     (concurrent — before any sync)

After merge: should "apple" be in the set?
```

There's no objectively correct answer — but the **OR-Set chooses add-wins**. Concurrent adds survive concurrent removes. This tends to be the least surprising behavior for users.

### Why Not Just Track Elements?

The naive approach: track a set of element names. But then:

```
A adds "apple" → set = {"apple"}
A removes "apple" → set = {}
B adds "apple" → set = {"apple"}

If B's add happened before A's remove was synced:
  Merge: A's set = {}, B's set = {"apple"}
  Union = {"apple"} — ignores the remove! ❌
  Intersection = {} — ignores the re-add! ❌
```

Neither union nor intersection gives the right answer.

### How It Works: Unique Tags

Each `Add` operation generates a globally unique tag (e.g., `replicaID:counter`). The set's real state is a collection of `(element, tag)` pairs, plus a set of tombstoned (removed) tags.

- `Add("apple")` → creates tag `A:1`, records pair `("apple", "A:1")`
- `Remove("apple")` → moves all *currently known* tags for "apple" into the tombstone set
- `Contains("apple")` → true if any non-tombstoned tag for "apple" exists
- `Merge` → union of all tag sets, union of tombstone sets

The key insight: `Remove` only tombstones tags it has *observed*. A concurrent `Add` on another replica creates a new tag that the removing replica hasn't seen yet — so it survives the tombstone.

### Example

```
Replica A adds "apple" → tag A:1
  A.adds = {"apple": {"A:1"}}

Replica B adds "apple" → tag B:1  (concurrent — before merge)
  B.adds = {"apple": {"B:1"}}

Replica A removes "apple":
  A.tombstones = {"A:1"}
  (B:1 not yet observed by A)

Merge A into B:
  adds = {"apple": {"A:1", "B:1"}}  (union)
  tombstones = {"A:1"}               (union)

Contains("apple")?
  Check A:1 → tombstoned ✗
  Check B:1 → not tombstoned ✓
  → true ✅

"apple" survives because B:1 was never removed.
```

---

## Getting Started

Copy the starter templates:

```bash
# From the module-03-registers-and-sets directory
cp ../starters/module-03-registers-and-sets/*.tmpl .
for f in *.tmpl; do mv "$f" "${f%.tmpl}"; done
```

This gives you four files:
- `lwwregister.go` — LWW-Register struct and function stubs
- `lwwregister_test.go` — test suite with TODOs
- `orset.go` — OR-Set struct and function stubs
- `orset_test.go` — test suite with TODOs

**Why .tmpl?** Template files let you reset to the starting state anytime. They're ignored by the Go language server, preventing IDE conflicts with the solution files.

### Need to Reset?

```bash
rm *.go
cp ../starters/module-03-registers-and-sets/*.tmpl .
for f in *.tmpl; do mv "$f" "${f%.tmpl}"; done
```

---

## The Challenge

### Part 1: LWW-Register

Implement `LWWRegister` in `package registers`:

```go
type LWWRegister struct {
    replicaID string
    // stores the current value and its timestamp
}

func NewLWWRegister(replicaID string) *LWWRegister
func (r *LWWRegister) Set(value string)                 // use current time
func (r *LWWRegister) SetAt(value string, ts int64)     // explicit timestamp (for tests)
func (r *LWWRegister) Get() string
func (r *LWWRegister) Merge(other *LWWRegister)         // higher timestamp wins
```

### Part 2: OR-Set

Implement `ORSet` in the same `package registers`:

```go
type ORSet struct {
    replicaID string
    counter   int
    // adds: element → set of live tags
    // tombstones: set of removed tags
}

func NewORSet(replicaID string) *ORSet
func (s *ORSet) Add(element string)
func (s *ORSet) Remove(element string)
func (s *ORSet) Contains(element string) bool
func (s *ORSet) Elements() []string   // sorted
func (s *ORSet) Merge(other *ORSet)
```

**Requirements:**
- Tags must be globally unique per Add — use `replicaID:counter` format
- `Remove` only tombstones currently-observed tags; concurrent adds survive
- `Elements()` returns a sorted slice (use `sort.Strings`)
- `Merge` is idempotent, commutative, and associative

---

## Workflow

### Step 1: Set Up Your Files

```bash
cp ../starters/module-03-registers-and-sets/*.tmpl .
for f in *.tmpl; do mv "$f" "${f%.tmpl}"; done
```

### Step 2: Add the Package Declaration

All four files need:
```go
package registers
```

### Step 3: Tackle LWW-Register First

It's simpler and builds intuition for timestamps. Fill in `lwwregister.go`, then make `lwwregister_test.go` pass.

```bash
go test . -run TestLWW
```

### Step 4: Tackle OR-Set

The OR-Set is more complex. Think through the data structures before writing code:
- How do you represent "a set of tags per element"?
- How do you generate a globally unique tag?
- What does Contains need to check?

Fill in `orset.go`, then make `orset_test.go` pass.

```bash
go test . -run TestOR
```

### Step 5: Run All Tests

```bash
go test .
```

Fix until all tests pass.

### Step 6: Compare with Solution

```bash
# View the solution
cat ../solutions/module-03-registers-and-sets/lwwregister_solution.go
cat ../solutions/module-03-registers-and-sets/orset_solution.go

# Run solution tests
cd ../solutions/module-03-registers-and-sets && go test -v .
cd ../../module-03-registers-and-sets
```

All green = module complete ✅

---

## Hints

### LWW-Register Hints

<details>
<summary><strong>Hint 1:</strong> What fields does LWWRegister need?</summary>

Think about what you need to resolve a conflict: you need to know both the value and *when* it was written. Store them together.

</details>

<details>
<summary><strong>Hint 2:</strong> What does Merge compare?</summary>

Merge needs to decide which replica's value to keep. It has one piece of information to compare per replica. The comparison determines which replica "wins."

</details>

<details>
<summary><strong>Hint 3:</strong> How do I make tests deterministic?</summary>

`Set` uses real wall-clock time, which is hard to control in tests. `SetAt` lets you pass an exact timestamp. Use `SetAt` in all your tests with explicit integer timestamps — that way the test is deterministic and readable.

</details>

### OR-Set Hints

<details>
<summary><strong>Hint 1:</strong> What data structures does ORSet need?</summary>

You need to track two things: which tags are associated with each element (so you know what to check or tombstone), and which tags have been removed. Think about what Go types naturally represent a set and a mapping.

</details>

<details>
<summary><strong>Hint 2:</strong> How do I generate a unique tag?</summary>

The tag just needs to be unique across all replicas for all time. You have two pieces of information that together uniquely identify any single Add operation. Combine them into a string.

</details>

<details>
<summary><strong>Hint 3:</strong> Why does the concurrent Add+Remove test pass?</summary>

When replica A removes an element, it tombstones only the tags it currently knows about. Replica B's tag was generated independently and hasn't been merged into A yet — so A has no knowledge of it and can't tombstone it. After merge, B's tag is live and the element survives.

</details>

<details>
<summary><strong>Hint 4:</strong> My Merge isn't working — elements disappear</summary>

When merging the adds maps, make sure you're taking the **union** of tag sets for each element. If you overwrite one replica's tags with the other's, you lose data. For each element, every tag from both replicas should end up in the merged set.

</details>

---

## Questions to Ponder

- **What are the limits of LWW?** What happens if two replicas have clocks that drift significantly? What data would you lose?

- **Why is Remove hard in distributed systems?** If you could guarantee message delivery order, would you still need tombstones?

- **What would a Remove-wins set look like?** Some use cases prefer remove-wins semantics. How would you implement that instead?

- **Tombstones grow forever.** Every remove adds to the tombstone set and it never shrinks. For a long-lived system with frequent removes, what would you do about this?

- **OR-Set vs 2P-Set.** A 2P-Set uses two G-Sets: one for adds, one for removes, with remove-wins. What scenarios favor each?

---

## Extension Challenge (Optional)

1. **Generic LWW-Register:** Rewrite `LWWRegister` using Go generics so it can store any comparable type, not just strings.

2. **OR-Map:** Build an `ORMap[V]` where values are LWW-Registers. Keys follow OR-Set semantics (add-wins), values use LWW merge. This is the foundation of Module 04's TodoList.

3. **Tombstone compaction:** Add a `Compact()` method to `ORSet` that removes entries from both `adds` and `tombstones` for elements where every tag is tombstoned. When is this safe to call?

---

## Next Module

Once you've completed this module, move on to [Module 04 - TodoList CRDT](../module-04-todolist/) to compose everything you've built into a real application-layer CRDT.
