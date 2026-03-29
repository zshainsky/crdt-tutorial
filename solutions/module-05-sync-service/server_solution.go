package syncservicesolution

import (
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
	todolist "github.com/zshainsky/crdt-tutorial/module-04-todolist"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

// Hub maintains the shared TodoList and all connected WebSocket clients.
type Hub struct {
	todos   *todolist.TodoList
	clients map[*websocket.Conn]bool
	mu      sync.Mutex
}

// NewHub creates a Hub with an empty TodoList.
func NewHub() *Hub {
	return &Hub{
		todos:   todolist.NewTodoList("server"),
		clients: make(map[*websocket.Conn]bool),
	}
}

// ServeHTTP upgrades the HTTP connection to WebSocket and runs the message loop.
// Hub implements http.Handler so it can be passed directly to http.Handle.
func (h *Hub) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("upgrade:", err)
		return
	}

	h.mu.Lock()
	h.clients[conn] = true
	h.mu.Unlock()

	defer func() {
		h.mu.Lock()
		delete(h.clients, conn)
		h.mu.Unlock()
		conn.Close()
	}()

	// Send current state to this new client immediately.
	h.sendStateTo(conn)

	for {
		_, data, err := conn.ReadMessage()
		if err != nil {
			break
		}
		msg, err := Decode(data)
		if err != nil {
			log.Println("decode:", err)
			continue
		}
		switch msg.Type {
		case TypeJoin:
			// Client requested a fresh state snapshot.
			h.sendStateTo(conn)
		case TypeOp:
			// Apply the operation and broadcast updated state to all clients.
			h.mu.Lock()
			msg.Apply(h.todos)
			state := StateMessage(h.todos)
			// Snapshot the client list before releasing the lock.
			conns := make([]*websocket.Conn, 0, len(h.clients))
			for c := range h.clients {
				conns = append(conns, c)
			}
			h.mu.Unlock()
			h.broadcastTo(state, conns)
		}
	}
}

// sendStateTo sends the current state snapshot to a single connection.
func (h *Hub) sendStateTo(conn *websocket.Conn) {
	h.mu.Lock()
	state := StateMessage(h.todos)
	h.mu.Unlock()
	data, err := Encode(state)
	if err != nil {
		return
	}
	if err := conn.WriteMessage(websocket.TextMessage, data); err != nil {
		log.Println("write:", err)
	}
}

// broadcastTo sends a message to a specific set of connections.
func (h *Hub) broadcastTo(msg Message, conns []*websocket.Conn) {
	data, err := Encode(msg)
	if err != nil {
		return
	}
	for _, conn := range conns {
		if err := conn.WriteMessage(websocket.TextMessage, data); err != nil {
			log.Println("broadcast write:", err)
		}
	}
}
