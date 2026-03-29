// Add package declaration (package todolist)
package todolist

import (
	"slices"
	"testing"
	"time"
)

// AddTask should return a non-empty ID; Tasks() should show the task with correct title and not completed
func TestTodoListAddTask(t *testing.T) {
	// Create a TodoList, add a task, assert the returned ID is non-empty
	// Assert Tasks() has exactly 1 task with the right title and Completed == false
	tl := NewTodoList("A")
	if tl.AddTask("my new task") == "" {
		t.Errorf("AddTask returned empty id")
	}

}

// Tasks() should be empty after a task is removed
func TestTodoListRemoveTask(t *testing.T) {
	// Add a task, remove it by ID, assert Tasks() returns an empty slice
	tl := NewTodoList("A")
	taskID := tl.AddTask("my new task")
	if tl.Tasks() == nil {
		t.Errorf("Tasks list should not be nil")
	}
	if len(tl.Tasks()) != 1 {
		t.Errorf("Tasks list should have length 1 but got: %d", len(tl.Tasks()))
	}

	tl.RemoveTask(taskID)
	if len(tl.Tasks()) > 0 {
		t.Errorf("Task with id: %s was not removed from the Tasks list. Tasks should be empty but was: %v", taskID, tl.Tasks())
	}

}

// SetTitle should update the task's title
func TestTodoListSetTitle(t *testing.T) {
	// Add a task, call SetTitle with a new title, assert Tasks()[0].Title matches
	tl := NewTodoList("A")

	wantTitle := "my new task"

	taskID := tl.AddTask("")
	tl.SetTitle(taskID, wantTitle)

	gotTitle := tl.Tasks()[0].Title
	if gotTitle != wantTitle {
		t.Errorf("Got title: %s does not match expected title: %s", gotTitle, wantTitle)
	}
}

// SetCompleted should toggle the task's completed state
func TestTodoListSetCompleted(t *testing.T) {
	// Add a task, set completed to true, assert Tasks()[0].Completed == true
	// Then set it back to false and assert again
	tl := NewTodoList("A")

	taskID := tl.AddTask("my new task")

	gotCompleted := tl.Tasks()[0].Completed
	if gotCompleted != false {
		t.Errorf("Got completed: %t, want: %t", gotCompleted, false)
	}

	tl.SetCompleted(taskID, true)
	gotCompleted = tl.Tasks()[0].Completed
	if gotCompleted != true {
		t.Errorf("Got completed: %t does not match expected completed: %t", gotCompleted, true)
	}
}

// After merging two lists, tasks from both replicas should be present
func TestTodoListMerge(t *testing.T) {
	// Create two TodoLists (different replica IDs), each adds a different task
	// Merge one into the other — assert both tasks are present
	A := NewTodoList("A")
	B := NewTodoList("B")

	taskIdA := A.AddTask("taskA")
	taskIdB := B.AddTask("taskB")

	A.Merge(B)

	tasks := A.Tasks()
	if len(tasks) != 2 {
		t.Errorf("expected 2 tasks after merge, got %d", len(tasks))
	}
	if !slices.ContainsFunc(tasks, func(f Task) bool {
		return f.ID == taskIdA && f.Title == "taskA"
	}) {
		t.Errorf("task id: %s is expected but not found", taskIdA)
	}
	if !slices.ContainsFunc(tasks, func(f Task) bool {
		return f.ID == taskIdB && f.Title == "taskB"
	}) {
		t.Errorf("task id: %s is expected but not found", taskIdB)

	}
}

// Merging the same state twice should not duplicate tasks
func TestTodoListMergeIdempotent(t *testing.T) {
	// Merge B into A, record the task count, merge B into A again
	// Assert the count did not change

	A := NewTodoList("A")
	B := NewTodoList("B")

	A.AddTask("taskA")
	B.AddTask("taskB")

	A.Merge(B)
	taskCount1 := len(A.Tasks())
	A.Merge(B)
	taskCount2 := len(A.Tasks())
	if taskCount1 != taskCount2 {
		t.Errorf("first merge task count: %d does not equal second merge task count: %d", taskCount1, taskCount2)
	}
}

// Merging in different orders should yield the same result
func TestTodoListMergeCommutative(t *testing.T) {
	// A←B and B←A should both result in the same number of tasks
	A := NewTodoList("A")
	B := NewTodoList("B")

	A.AddTask("taskA")
	B.AddTask("taskB")

	A.Merge(B)
	taskCount1 := len(A.Tasks())
	B.Merge(A)
	taskCount2 := len(B.Tasks())
	if taskCount1 != taskCount2 {
		t.Errorf("first merge task count: %d does not equal second merge task count: %d", taskCount1, taskCount2)
	}
	if !slices.Equal(A.Tasks(), B.Tasks()) {
		t.Errorf("A Tasks: %v does not equal B Tasks: %v", A.Tasks(), B.Tasks())
	}
}

// When two replicas concurrently edit the same task's title, the higher timestamp wins
func TestTodoListConcurrentTitleEdit(t *testing.T) {
	// Replica A adds a task; B merges to get a copy
	// Both edit the title using SetTitleAt with different timestamps
	// After merge, assert the edit with the higher timestamp survived
	now := time.Now()
	A := NewTodoList("A")
	B := NewTodoList("B")

	taskID := A.AddTask("taskA")
	B.Merge(A)

	nowPlusOneDay := now.AddDate(0, 0, 1)
	A.SetTitleAt(taskID, "A's Edit", nowPlusOneDay.UnixNano())
	nowPlusTwoDay := now.AddDate(0, 0, 2)
	B.SetTitleAt(taskID, "B's Edit", nowPlusTwoDay.UnixNano())

	A.Merge(B)
	if A.Tasks()[0].Title != "B's Edit" {
		t.Errorf("the title: %s for taskID: %s should be %s", A.Tasks()[0].Title, A.Tasks()[0].ID, "B's Edit")
	}

}

// A task added concurrently with a remove should survive (OR-Set add-wins)
func TestTodoListConcurrentAddRemove(t *testing.T) {
	// A adds task one; B merges and then removes task one
	// A adds task two concurrently (before any merge)
	// Merge A and B — assert task one is gone but task two is present

	A := NewTodoList("A")
	B := NewTodoList("B")

	taskID1 := A.AddTask("taskA1")
	B.Merge(A)

	B.RemoveTask(taskID1)

	taskID2 := A.AddTask("taskA2")

	A.Merge(B)
	// Should not contain TaskID1
	if slices.ContainsFunc(A.Tasks(), func(f Task) bool {
		return f.ID == taskID1
	}) {
		t.Errorf("Tasks %v should not contain deleted task id: %s", A.Tasks(), taskID1)
	}
	// Should contain TaskID2
	if !slices.ContainsFunc(A.Tasks(), func(f Task) bool {
		return f.ID == taskID2
	}) {
		t.Errorf("Tasks %v should not contain deleted task id: %s", A.Tasks(), taskID2)
	}
}

// Merge should transfer registers even for tasks that have been removed
// (orphan registers must sync so no data is lost in multi-hop propagation)
func TestTodoListMergeOrphanRegisters(t *testing.T) {
	// Create replica A and add a task
	// Create replica B and merge A's state so B has the same task

	// B removes the task
	// B also updates the removed task's title to a new value
	// (SetTitle still works on a removed task — the register exists even after removal)

	// A merges B
	// Assert: the task is gone from A (OR-Set remove propagated correctly)

	// Assert: A's internal title register for that task ID matches B's updated value
	// Hint: Tasks() won't show it (task is removed), so read the register directly
	//       using t.title[id] — you're in the same package so you can access it

	A := NewTodoList("A")
	B := NewTodoList("B")

	taskID1 := A.AddTask("taskA1")

	B.Merge(A)
	B.RemoveTask(taskID1)
	B.SetTitle(taskID1, "B Updated Title")

	A.Merge(B)
	if slices.ContainsFunc(A.Tasks(), func(f Task) bool {
		return f.ID == taskID1
	}) {
		t.Errorf("Tasks %v should not contain deleted task id: %s", A.Tasks(), taskID1)
	}

	if A.title[taskID1].Get() != "B Updated Title" {
		t.Errorf("Orphan register not synced")
	}

}
