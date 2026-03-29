# Module 04 — TodoList CRDT
## Application-Layer CRDT Composition

**Concepts:** Map-of-CRDTs pattern, application-layer composition, convergent data models
**Challenge:** Build a collaborative todo list from the CRDTs you've already implemented
**Time:** ~2-3 hours

---

## Introduction

You've built four CRDTs from scratch:
- **G-Counter** — grow-only counter
- **PN-Counter** — increment and decrement
- **LWW-Register** — single value, last write wins
- **OR-Set** — add-wins membership set

Now you compose them into something real: a **collaborative todo list** that multiple replicas can edit simultaneously and always converge to the same state—no coordination required.

This is the payoff module. Everything you've built gets used.

---

## The Design

A TodoList needs to support:
- Adding tasks (with a title)
- Removing tasks
- Editing a task's title
- Marking a task complete or incomplete
- Merging two lists from different replicas

Here's the data model:

```
TodoList
├── taskIDs: ORSet          ← which task IDs are "live"
├── titles:  map[ID] → LWWRegister   ← title per task
└── completed: map[ID] → LWWRegister ← "true"/"false" per task
```

**Why OR-Set for IDs?** Because we want add-wins semantics. If replica A removes a task while replica B adds it concurrently, B's add survives.

**Why LWW-Register for fields?** Because each field is a single value with concurrent-write conflicts. Last write wins — the higher timestamp takes it.

**Why map-of-registers instead of a single struct?** Because we need to merge fields independently. If A edits the title and B edits the completed state at the same time, both changes should survive — they touch different registers.

### ASCII Data Model

```
Replica A                           Replica B
─────────────────────────────────   ─────────────────────────────────
taskIDs (ORSet):                    taskIDs (ORSet):
  "A:1" → [tag: A:1]                 "A:1" → [tag: A:1] (via merge)
  "A:2" → [tag: A:2]                 "B:1" → [tag: B:1]

titles:                             titles:
  "A:1" → LWW{"buy milk", t=100}     "A:1" → LWW{"buy OAT milk", t=200}
  "A:2" → LWW{"walk dog", t=110}     "B:1" → LWW{"read book", t=150}

After A.Merge(B):
  taskIDs: {"A:1", "A:2", "B:1"}
  titles:
    "A:1" → LWW{"buy OAT milk", t=200}  ← B's edit wins (higher ts)
    "A:2" → LWW{"walk dog", t=110}
    "B:1" → LWW{"read book", t=150}
```

---

## The Problem: Merging a Map of CRDTs

The tricky part isn't merging a single CRDT — you already know how. The challenge is merging a *dynamic collection* of them.

When merging two TodoLists:

1. **New tasks on the other side:** The other replica has registers for task IDs we've never seen. We need to adopt those registers.

2. **Shared tasks with conflicting edits:** Both replicas have registers for the same task ID, but they may have different values (concurrent writes). Merge each register using LWW semantics.

3. **Orphan registers:** A task may have been removed from the OR-Set but its registers still exist in the map. That's fine — `Tasks()` filters by what the OR-Set says is live. The stale registers are harmless.

### Merge Algorithm (in English)

```
Merge(other):
  1. Merge the OR-Sets (union tag sets, union tombstones)
  2. For each task ID in other.titles:
       If we don't have a title register for this ID, create one
       Merge our register with theirs
  3. Same for other.completed
```

---

## Getting Started

Copy the starter templates:

```bash
# From the module-04-todolist directory
cp ../starters/module-04-todolist/*.tmpl .
for f in *.tmpl; do mv "$f" "${f%.tmpl}"; done
```

This gives you:
- `todolist.go` — struct definitions and function stubs
- `todolist_test.go` — test suite with TODOs

### Need to Reset?

```bash
rm *.go
cp ../starters/module-04-todolist/*.tmpl .
for f in *.tmpl; do mv "$f" "${f%.tmpl}"; done
```

---

## The Challenge

### Types

```go
type Task struct {
    ID        string
    Title     string
    Completed bool
}

type TodoList struct {
    replicaID string
    counter   int
    // OR-Set for live task IDs
    // map[taskID]*LWWRegister for titles
    // map[taskID]*LWWRegister for completed ("true"/"false")
}
```

### API

```go
func NewTodoList(replicaID string) *TodoList

func (t *TodoList) AddTask(title string) string          // returns task ID
func (t *TodoList) RemoveTask(id string)
func (t *TodoList) SetTitle(id, title string)
func (t *TodoList) SetTitleAt(id, title string, ts int64)     // for deterministic tests
func (t *TodoList) SetCompleted(id string, done bool)
func (t *TodoList) SetCompletedAt(id string, done bool, ts int64) // for deterministic tests
func (t *TodoList) Tasks() []Task                        // live tasks, sorted by ID
func (t *TodoList) Merge(other *TodoList)
```

### Requirements

1. **Task IDs** must be globally unique — use `replicaID:counter` format
2. **AddTask** adds the ID to the OR-Set and creates LWW-Registers for title and completed
3. **RemoveTask** removes the ID from the OR-Set; registers can remain (orphans are fine)
4. **SetTitle / SetCompleted** silently no-op for unknown task IDs
5. **SetTitleAt / SetCompletedAt** call the underlying register's `SetAt` — needed for deterministic tests
6. **Tasks()** returns only tasks whose IDs are live in the OR-Set, in sorted order
7. **Merge** must handle task IDs the other side has that we don't (adopt their registers)
8. **Convergence:** all replicas that see the same operations must produce the same `Tasks()` output

### Import

```go
import registers "github.com/zshainsky/crdt-tutorial/module-03-registers-and-sets"
```

---

## Workflow

### Step 1: Set Up Your Files

```bash
cp ../starters/module-04-todolist/*.tmpl .
for f in *.tmpl; do mv "$f" "${f%.tmpl}"; done
```

### Step 2: Add the Package Declaration

Both files need:
```go
package todolist
```

### Step 3: Define the Struct

Think through what fields `TodoList` needs before writing any methods. Sketch the data model on paper first if it helps.

### Step 4: Implement Methods One at a Time

Start with `NewTodoList` → `AddTask` → `Tasks()`. Get those working together first, then add `SetTitle`, `SetCompleted`, `RemoveTask`. Save `Merge` for last.

```bash
go test . -run TestTodoListAddTask
```

### Step 5: Implement Merge

This is the hardest part. Re-read "The Problem: Merging a Map of CRDTs" above before writing this function.

### Step 6: Run All Tests

```bash
go test .
```

Fix until all tests pass.

### Step 7: Compare with Solution

```bash
cat ../solutions/module-04-todolist/todolist_solution.go
cd ../solutions/module-04-todolist && go test -v .
cd ../../module-04-todolist
```

All green = module complete ✅

---

## Hints

<details>
<summary><strong>Hint 1:</strong> How do I generate unique task IDs?</summary>

You need something that is unique across all replicas. You already have a `replicaID` that identifies this node, and a `counter` that increments with every add. Combine them. The same strategy is used internally by OR-Set to generate its tags.

</details>

<details>
<summary><strong>Hint 2:</strong> How do I store "true" or "false" in a LWW-Register?</summary>

`LWWRegister` stores strings. Represent the boolean as the string `"true"` or `"false"`. In `Tasks()`, compare the register's value to `"true"` to get a bool back.

</details>

<details>
<summary><strong>Hint 3:</strong> My Merge doesn't handle new tasks from the other side</summary>

When merging, the other replica may have task IDs (and registers) that you've never seen. Before calling `Merge` on a register, check whether you have that register. If not, create a fresh one for this replica and then merge. A fresh `NewLWWRegister` has timestamp 0, so the other side's data will always win on the first merge.

</details>

<details>
<summary><strong>Hint 4:</strong> Tasks() shows removed tasks</summary>

`Tasks()` should only show tasks that are currently live in the OR-Set. Use the OR-Set's `Elements()` method — it returns only the live element IDs. Iterate *that list*, not the raw map of registers.

</details>

<details>
<summary><strong>Hint 5:</strong> The concurrent title edit test fails</summary>

For `TestTodoListConcurrentTitleEdit`, the test calls `SetTitleAt` on the underlying register with explicit timestamps. Make sure your `SetTitleAt` implementation delegates to the register's `SetAt` method. If you're using `Set` (which uses wall-clock time), fast test execution may not preserve the expected ordering.

</details>

---

## Questions to Ponder

- **Orphan registers.** When a task is removed, its LWW-Registers stay in the map forever. For a long-running application, is this a problem? How would you clean them up safely?

- **What if the same task is added by two different replicas?** They'll get different IDs (because IDs include replicaID), so they'll appear as two separate tasks after merge. Is that the right behavior?

- **Clock skew.** LWW-Register relies on timestamps. What happens if replica B's system clock is 10 minutes behind replica A's? Who wins concurrent writes?

- **Add-wins vs remove-wins.** This TodoList uses OR-Set (add-wins). Can you think of a use case where you'd want remove-wins instead?

- **Completed as LWW.** If two replicas toggle completed in opposite directions concurrently, LWW picks one arbitrarily (higher timestamp). For a "done" checkbox, is that acceptable? What alternative would you use?

---

## Extension Challenge (Optional)

1. **Add a `priority` field** to each task — an integer stored in its own LWW-Register. Implement `SetPriority(id string, p int)` and have `Tasks()` return tasks sorted by priority instead of ID.

2. **Observe orphan state** — add a `Debug()` method that prints all registers (including for removed tasks). Notice that remove doesn't delete the register data; it only updates the OR-Set.

3. **Compact tombstones** — implement a `Compact()` method that removes registers for task IDs that are fully tombstoned in the OR-Set and known by all replicas. When is it *safe* to compact?

---

## Next Module

Once you've completed this module, move on to [Module 05 - Real-Time Sync Service](../module-05-sync-service/) to expose your TodoList CRDT over WebSockets and build a collaborative browser client.
