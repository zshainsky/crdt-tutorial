package syncservicesolution

import (
	"encoding/json"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/websocket"
)

// dial creates a test WebSocket client connected to srv.
func dial(t *testing.T, srv *httptest.Server) *websocket.Conn {
	t.Helper()
	u := "ws" + strings.TrimPrefix(srv.URL, "http")
	conn, _, err := websocket.DefaultDialer.Dial(u, nil)
	if err != nil {
		t.Fatalf("dial: %v", err)
	}
	t.Cleanup(func() { conn.Close() })
	return conn
}

// readMsg reads and decodes the next message from conn (2-second deadline).
func readMsg(t *testing.T, conn *websocket.Conn) Message {
	t.Helper()
	conn.SetReadDeadline(time.Now().Add(2 * time.Second))
	_, data, err := conn.ReadMessage()
	if err != nil {
		t.Fatalf("read: %v", err)
	}
	var m Message
	if err := json.Unmarshal(data, &m); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	return m
}

// sendMsg encodes and sends msg to conn.
func sendMsg(t *testing.T, conn *websocket.Conn, m Message) {
	t.Helper()
	data, _ := json.Marshal(m)
	if err := conn.WriteMessage(websocket.TextMessage, data); err != nil {
		t.Fatalf("write: %v", err)
	}
}

func TestWebSocketUpgrade(t *testing.T) {
	hub := NewHub()
	srv := httptest.NewServer(hub)
	defer srv.Close()

	conn := dial(t, srv)
	_ = conn // connection succeeded — no error from dial
}

func TestStateSync(t *testing.T) {
	hub := NewHub()
	srv := httptest.NewServer(hub)
	defer srv.Close()

	conn := dial(t, srv)
	msg := readMsg(t, conn)
	if msg.Type != TypeState {
		t.Errorf("expected %q message on connect, got %q", TypeState, msg.Type)
	}
}

func TestStateSyncAfterJoin(t *testing.T) {
	hub := NewHub()
	srv := httptest.NewServer(hub)
	defer srv.Close()

	conn := dial(t, srv)
	readMsg(t, conn) // drain initial state

	sendMsg(t, conn, Message{Type: TypeJoin})
	msg := readMsg(t, conn)
	if msg.Type != TypeState {
		t.Errorf("expected state after join, got %q", msg.Type)
	}
}

func TestOperationBroadcast(t *testing.T) {
	hub := NewHub()
	srv := httptest.NewServer(hub)
	defer srv.Close()

	connA := dial(t, srv)
	connB := dial(t, srv)

	readMsg(t, connA) // drain initial state for A
	readMsg(t, connB) // drain initial state for B

	// A adds a task
	sendMsg(t, connA, Message{Type: TypeOp, Action: ActionAdd, Title: "buy milk"})

	// Both A and B should receive the updated state
	msgA := readMsg(t, connA)
	msgB := readMsg(t, connB)

	for _, msg := range []Message{msgA, msgB} {
		if msg.Type != TypeState {
			t.Errorf("expected state broadcast, got %q", msg.Type)
		}
		if len(msg.Tasks) != 1 || msg.Tasks[0].Title != "buy milk" {
			t.Errorf("expected task 'buy milk', got %v", msg.Tasks)
		}
	}
}

func TestConcurrentOperations(t *testing.T) {
	hub := NewHub()
	srv := httptest.NewServer(hub)
	defer srv.Close()

	connA := dial(t, srv)
	connB := dial(t, srv)

	readMsg(t, connA)
	readMsg(t, connB)

	// A adds a task — both clients receive the broadcast
	sendMsg(t, connA, Message{Type: TypeOp, Action: ActionAdd, Title: "task A"})
	readMsg(t, connA)
	readMsg(t, connB)

	// B adds a task — both clients receive the broadcast
	sendMsg(t, connB, Message{Type: TypeOp, Action: ActionAdd, Title: "task B"})
	readMsg(t, connA)
	msgB := readMsg(t, connB)

	// Final state seen by B should have both tasks
	if len(msgB.Tasks) != 2 {
		t.Errorf("expected 2 tasks after concurrent ops, got %d: %v", len(msgB.Tasks), msgB.Tasks)
	}
}

func TestReconnection(t *testing.T) {
	hub := NewHub()
	srv := httptest.NewServer(hub)
	defer srv.Close()

	// First connection: add a task, wait for broadcast confirmation, then disconnect
	connA := dial(t, srv)
	readMsg(t, connA) // drain initial state
	sendMsg(t, connA, Message{Type: TypeOp, Action: ActionAdd, Title: "persisted task"})
	readMsg(t, connA) // wait for server to confirm the operation was applied
	connA.Close()

	// Reconnect: new client should receive state with the task
	connB := dial(t, srv)
	state := readMsg(t, connB)
	if state.Type != TypeState {
		t.Fatalf("expected state on reconnect, got %q", state.Type)
	}
	if len(state.Tasks) != 1 || state.Tasks[0].Title != "persisted task" {
		t.Errorf("expected persisted task on reconnect, got %v", state.Tasks)
	}
}
