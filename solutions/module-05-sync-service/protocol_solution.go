package syncservicesolution

import (
	"encoding/json"

	todolist "github.com/zshainsky/crdt-tutorial/module-04-todolist"
)

// Message type constants.
const (
	TypeJoin  = "join"
	TypeState = "state"
	TypeOp    = "op"
)

// Operation action constants.
const (
	ActionAdd          = "add"
	ActionRemove       = "remove"
	ActionSetTitle     = "setTitle"
	ActionSetCompleted = "setCompleted"
)

// Message is the JSON envelope for all WebSocket communication.
type Message struct {
	Type      string `json:"type"`
	Action    string `json:"action,omitempty"`
	ID        string `json:"id,omitempty"`
	Title     string `json:"title,omitempty"`
	Completed bool   `json:"completed,omitempty"`
	Tasks     []Task `json:"tasks,omitempty"`
}

// Task is the wire representation of a single todo item.
type Task struct {
	ID        string `json:"id"`
	Title     string `json:"title"`
	Completed bool   `json:"completed"`
}

// Apply executes the operation described by m against todos.
// Returns the new task ID for "add" actions, empty string otherwise.
func (m *Message) Apply(todos *todolist.TodoList) string {
	switch m.Action {
	case ActionAdd:
		return todos.AddTask(m.Title)
	case ActionRemove:
		todos.RemoveTask(m.ID)
	case ActionSetTitle:
		todos.SetTitle(m.ID, m.Title)
	case ActionSetCompleted:
		todos.SetCompleted(m.ID, m.Completed)
	}
	return ""
}

// StateMessage builds a TypeState message from the current TodoList.
func StateMessage(todos *todolist.TodoList) Message {
	raw := todos.Tasks()
	tasks := make([]Task, len(raw))
	for i, t := range raw {
		tasks[i] = Task{ID: t.ID, Title: t.Title, Completed: t.Completed}
	}
	return Message{Type: TypeState, Tasks: tasks}
}

// Encode serializes a Message to JSON bytes.
func Encode(m Message) ([]byte, error) {
	return json.Marshal(m)
}

// Decode parses JSON bytes into a Message.
func Decode(data []byte) (Message, error) {
	var m Message
	err := json.Unmarshal(data, &m)
	return m, err
}
