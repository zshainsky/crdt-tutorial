package todolistsolution

import "testing"

func TestTodoListAddTask(t *testing.T) {
	tl := NewTodoList("A")
	id := tl.AddTask("buy milk")

	if id == "" {
		t.Fatal("AddTask should return a non-empty task ID")
	}

	tasks := tl.Tasks()
	if len(tasks) != 1 {
		t.Fatalf("expected 1 task, got %d", len(tasks))
	}
	if tasks[0].ID != id {
		t.Errorf("expected task ID %q, got %q", id, tasks[0].ID)
	}
	if tasks[0].Title != "buy milk" {
		t.Errorf("expected title %q, got %q", "buy milk", tasks[0].Title)
	}
	if tasks[0].Completed {
		t.Error("new task should not be completed")
	}
}

func TestTodoListRemoveTask(t *testing.T) {
	tl := NewTodoList("A")
	id := tl.AddTask("buy milk")
	tl.RemoveTask(id)

	if len(tl.Tasks()) != 0 {
		t.Errorf("expected 0 tasks after removal, got %d", len(tl.Tasks()))
	}
}

func TestTodoListMultipleTasks(t *testing.T) {
	tl := NewTodoList("A")
	tl.AddTask("buy milk")
	tl.AddTask("walk dog")
	tl.AddTask("read book")

	if len(tl.Tasks()) != 3 {
		t.Errorf("expected 3 tasks, got %d", len(tl.Tasks()))
	}
}

func TestTodoListSetTitle(t *testing.T) {
	tl := NewTodoList("A")
	id := tl.AddTask("buy milk")
	tl.SetTitle(id, "buy oat milk")

	tasks := tl.Tasks()
	if tasks[0].Title != "buy oat milk" {
		t.Errorf("expected title %q, got %q", "buy oat milk", tasks[0].Title)
	}
}

func TestTodoListSetCompleted(t *testing.T) {
	tl := NewTodoList("A")
	id := tl.AddTask("buy milk")

	tl.SetCompleted(id, true)
	if !tl.Tasks()[0].Completed {
		t.Error("expected task to be completed")
	}

	tl.SetCompleted(id, false)
	if tl.Tasks()[0].Completed {
		t.Error("expected task to be uncompleted after toggling back")
	}
}

func TestTodoListMerge(t *testing.T) {
	a := NewTodoList("A")
	b := NewTodoList("B")
	a.AddTask("buy milk")
	b.AddTask("walk dog")

	a.Merge(b)

	if len(a.Tasks()) != 2 {
		t.Errorf("expected 2 tasks after merge, got %d", len(a.Tasks()))
	}
}

func TestTodoListMergeIdempotent(t *testing.T) {
	a := NewTodoList("A")
	b := NewTodoList("B")
	a.AddTask("buy milk")
	b.AddTask("walk dog")

	a.Merge(b)
	count1 := len(a.Tasks())
	a.Merge(b)
	count2 := len(a.Tasks())

	if count1 != count2 {
		t.Errorf("merge not idempotent: %d tasks then %d tasks", count1, count2)
	}
}

func TestTodoListMergeCommutative(t *testing.T) {
	// Scenario 1: A absorbs B
	a1 := NewTodoList("A")
	b1 := NewTodoList("B")
	a1.AddTask("from A")
	b1.AddTask("from B")
	a1.Merge(b1)

	// Scenario 2: B absorbs A
	a2 := NewTodoList("A")
	b2 := NewTodoList("B")
	a2.AddTask("from A")
	b2.AddTask("from B")
	b2.Merge(a2)

	if len(a1.Tasks()) != len(b2.Tasks()) {
		t.Errorf("merge not commutative: %d vs %d tasks", len(a1.Tasks()), len(b2.Tasks()))
	}
	if len(a1.Tasks()) != 2 {
		t.Errorf("expected 2 tasks, got %d", len(a1.Tasks()))
	}
}

func TestTodoListConcurrentTitleEdit(t *testing.T) {
	// Both replicas start with the same task
	a := NewTodoList("A")
	id := a.AddTask("original title")

	b := NewTodoList("B")
	b.Merge(a)

	// Concurrent edits: A edits at t=100, B edits at t=200
	a.SetTitleAt(id, "A's edit", 100)
	b.SetTitleAt(id, "B's edit", 200)

	// After merge, B's edit should win (higher timestamp)
	a.Merge(b)
	tasks := a.Tasks()
	if tasks[0].Title != "B's edit" {
		t.Errorf("expected B's edit to win (higher timestamp), got %q", tasks[0].Title)
	}
}

func TestTodoListConcurrentAddRemove(t *testing.T) {
	// A adds a task; B gets a copy then removes it; A also adds another task concurrently.
	// After merge: A's second add should survive (OR-Set add-wins semantics for that element).
	a := NewTodoList("A")
	id1 := a.AddTask("task one")

	b := NewTodoList("B")
	b.Merge(a)

	// B removes the task A added
	b.RemoveTask(id1)

	// A adds a second task concurrently
	a.AddTask("task two")

	// After merge: task one removed (B's remove wins for that specific add),
	// task two survives (B never saw it)
	a.Merge(b)
	tasks := a.Tasks()

	for _, task := range tasks {
		if task.ID == id1 {
			t.Errorf("task one should be removed, but it appears in Tasks()")
		}
	}
	if len(tasks) != 1 {
		t.Errorf("expected 1 task (task two only), got %d: %v", len(tasks), tasks)
	}
	if tasks[0].Title != "task two" {
		t.Errorf("expected task two to survive, got %q", tasks[0].Title)
	}
}

func TestTodoListSetTitleIgnoresUnknownID(t *testing.T) {
	tl := NewTodoList("A")
	tl.SetTitle("nonexistent", "should not panic") // must not panic
}

func TestTodoListSetCompletedIgnoresUnknownID(t *testing.T) {
	tl := NewTodoList("A")
	tl.SetCompleted("nonexistent", true) // must not panic
}
