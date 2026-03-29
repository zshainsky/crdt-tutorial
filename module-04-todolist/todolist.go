// Add package declaration (package todolist)
package todolist

// TODO: Import "fmt" and the registers package:
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
	// Add a field for the OR-Set of live task IDs
	taskID *registers.ORSet
	// Add a map from task ID to its title register (LWWRegister)
	title map[string]*registers.LWWRegister
	// Add a map from task ID to its completed register (LWWRegister)
	completed map[string]*registers.LWWRegister
}

// NewTodoList creates a new TodoList for the given replica.
func NewTodoList(replicaID string) *TodoList {
	// Initialize the OR-Set and both maps, then return
	return &TodoList{
		replicaID: replicaID,
		taskID:    registers.NewORSet(replicaID),
		title:     make(map[string]*registers.LWWRegister),
		completed: make(map[string]*registers.LWWRegister),
	}
}

// AddTask creates a new task with the given title and returns its generated ID.
// Task IDs must be globally unique — use replicaID and counter together.
func (t *TodoList) AddTask(title string) string {
	// Increment the counter
	// Generate a unique task ID from replicaID and counter
	// Add the task ID to the OR-Set
	// Create a title register and set its initial value
	// Create a completed register and set its initial value to "false"
	// Return the task ID
	t.counter++
	taskID := fmt.Sprintf("%s:%d", t.replicaID, t.counter)
	t.taskID.Add(taskID)

	t.title[taskID] = registers.NewLWWRegister(t.replicaID)
	t.SetTitle(taskID, title)

	t.completed[taskID] = registers.NewLWWRegister(t.replicaID)
	t.SetCompleted(taskID, false)
	return taskID
}

// RemoveTask removes a task by ID.
func (t *TodoList) RemoveTask(id string) {
	// Remove the task ID from the OR-Set
	t.taskID.Remove(id)
}

// SetTitle updates the title of an existing task.
// Silently ignores unknown task IDs.
func (t *TodoList) SetTitle(id, title string) {
	// If the task's title register exists, update it with the new title
	if _, ok := t.title[id]; ok {
		t.title[id].Set(title)
	}
}

// SetTitleAt updates the title with an explicit timestamp.
// Use in tests for deterministic conflict resolution.
func (t *TodoList) SetTitleAt(id, title string, ts int64) {
	// If the task's title register exists, update it using the explicit timestamp
	if _, ok := t.title[id]; ok {
		t.title[id].SetAt(title, ts)
	}
}

// SetCompleted marks a task as completed or not.
// Silently ignores unknown task IDs.
func (t *TodoList) SetCompleted(id string, done bool) {
	// Convert the bool to the string "true" or "false"
	// If the task's completed register exists, update it
	isCompleted := "false"
	if done {
		isCompleted = "true"
	}

	if _, ok := t.completed[id]; ok {
		t.completed[id].Set(isCompleted)
	}
}

// SetCompletedAt marks a task completed with an explicit timestamp.
// Use in tests for deterministic conflict resolution.
func (t *TodoList) SetCompletedAt(id string, done bool, ts int64) {
	// Same as SetCompleted but use the explicit timestamp
	isCompleted := "false"
	if done {
		isCompleted = "true"
	}

	if _, ok := t.completed[id]; ok {
		t.completed[id].SetAt(isCompleted, ts)
	}
}

// Tasks returns all live tasks in task-ID order.
// Tasks whose IDs have been removed are excluded.
func (t *TodoList) Tasks() []Task {
	// Get the list of live task IDs from the OR-Set (already sorted)
	// For each ID, read the title and completed registers and build a Task
	// Return the slice of Tasks
	var res []Task
	for _, id := range t.taskID.Elements() {
		var isCompleted bool
		if t.completed[id].Get() == "true" {
			isCompleted = true
		}

		res = append(res, Task{
			ID:        id,
			Title:     t.title[id].Get(),
			Completed: isCompleted,
		})
	}
	return res
}

// Merge combines another TodoList's state into this one.
func (t *TodoList) Merge(other *TodoList) {
	// Merge the OR-Sets of task IDs
	// For each task ID the other side knows about, merge its title register
	//     (create a new register for this replica if we don't have one yet)
	// Do the same for each task's completed register
	t.taskID.Merge(other.taskID)
	for id, otherTitle := range other.title {
		if t.title[id] == nil {
			t.title[id] = registers.NewLWWRegister(t.replicaID)
		}
		t.title[id].Merge(otherTitle)
	}
	for id, otherCompleted := range other.completed {
		if t.completed[id] == nil {
			t.completed[id] = registers.NewLWWRegister(t.replicaID)
		}
		t.completed[id].Merge(otherCompleted)
	}
}
