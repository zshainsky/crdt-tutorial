package todolistsolution

import (
	"fmt"

	registers "github.com/zshainsky/crdt-tutorial/module-03-registers-and-sets"
)

// Task is a snapshot of a single task's current state.
type Task struct {
	ID        string
	Title     string
	Completed bool
}

// TodoList is a CRDT-based collaborative todo list.
// It composes an OR-Set (for task membership) and LWW-Registers
// (for each task's mutable fields) to achieve convergent merge.
type TodoList struct {
	replicaID string
	counter   int
	taskIDs   *registers.ORSet
	titles    map[string]*registers.LWWRegister // task ID → title register
	completed map[string]*registers.LWWRegister // task ID → "true"/"false" register
}

// NewTodoList creates a new TodoList for the given replica.
func NewTodoList(replicaID string) *TodoList {
	return &TodoList{
		replicaID: replicaID,
		taskIDs:   registers.NewORSet(replicaID),
		titles:    make(map[string]*registers.LWWRegister),
		completed: make(map[string]*registers.LWWRegister),
	}
}

// AddTask creates a new task with the given title and returns its generated ID.
// Task IDs are globally unique: replicaID:counter.
func (t *TodoList) AddTask(title string) string {
	t.counter++
	id := fmt.Sprintf("%s:%d", t.replicaID, t.counter)
	t.taskIDs.Add(id)
	t.titles[id] = registers.NewLWWRegister(t.replicaID)
	t.titles[id].Set(title)
	t.completed[id] = registers.NewLWWRegister(t.replicaID)
	t.completed[id].Set("false")
	return id
}

// RemoveTask removes a task by ID using OR-Set semantics (add-wins on concurrent add+remove).
func (t *TodoList) RemoveTask(id string) {
	t.taskIDs.Remove(id)
}

// SetTitle updates the title of an existing task.
// Silently ignores unknown task IDs.
func (t *TodoList) SetTitle(id, title string) {
	if t.titles[id] == nil {
		return
	}
	t.titles[id].Set(title)
}

// SetTitleAt updates the title with an explicit timestamp.
// Use in tests for deterministic conflict resolution.
func (t *TodoList) SetTitleAt(id, title string, ts int64) {
	if t.titles[id] == nil {
		return
	}
	t.titles[id].SetAt(title, ts)
}

// SetCompleted marks a task as completed or not.
// Silently ignores unknown task IDs.
func (t *TodoList) SetCompleted(id string, done bool) {
	if t.completed[id] == nil {
		return
	}
	val := "false"
	if done {
		val = "true"
	}
	t.completed[id].Set(val)
}

// SetCompletedAt marks a task completed with an explicit timestamp.
// Use in tests for deterministic conflict resolution.
func (t *TodoList) SetCompletedAt(id string, done bool, ts int64) {
	if t.completed[id] == nil {
		return
	}
	val := "false"
	if done {
		val = "true"
	}
	t.completed[id].SetAt(val, ts)
}

// Tasks returns all live tasks, ordered by task ID.
// Tasks whose IDs have been removed from the OR-Set are excluded.
func (t *TodoList) Tasks() []Task {
	ids := t.taskIDs.Elements() // sorted by OR-Set
	tasks := make([]Task, 0, len(ids))
	for _, id := range ids {
		title := ""
		if t.titles[id] != nil {
			title = t.titles[id].Get()
		}
		done := false
		if t.completed[id] != nil {
			done = t.completed[id].Get() == "true"
		}
		tasks = append(tasks, Task{ID: id, Title: title, Completed: done})
	}
	return tasks
}

// Merge combines another TodoList's state into this one.
// Merges the OR-Set of IDs and all per-task LWW-Registers.
func (t *TodoList) Merge(other *TodoList) {
	// Merge the set of live task IDs
	t.taskIDs.Merge(other.taskIDs)

	// Merge title registers for every task the other side knows about
	for id, otherTitle := range other.titles {
		if t.titles[id] == nil {
			t.titles[id] = registers.NewLWWRegister(t.replicaID)
		}
		t.titles[id].Merge(otherTitle)
	}

	// Merge completed registers for every task the other side knows about
	for id, otherCompleted := range other.completed {
		if t.completed[id] == nil {
			t.completed[id] = registers.NewLWWRegister(t.replicaID)
		}
		t.completed[id].Merge(otherCompleted)
	}
}
