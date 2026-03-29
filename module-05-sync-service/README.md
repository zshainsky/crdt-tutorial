# Module 05 — Real-Time Sync Service

## Networked CRDTs: WebSocket Server + Browser Client

**Concepts:** Client-server sync, WebSocket protocol, state vs operation-based replication
**Challenge:** Expose your TodoList CRDT over the network with a collaborative web UI
**Prerequisites:** Module 04 complete (server imports your TodoList directly)
**Time:** ~3–4 hours

---

## Introduction

You've built a TodoList CRDT that converges locally. Now expose it over the network.

In this module you write a WebSocket server in Go and a minimal browser client in JavaScript. Multiple users can open the same URL and edit the todo list simultaneously — adds, completions, title edits, and deletes all propagate in real time. Conflicts resolve automatically using the merge rules you built in Module 04.

This is the payoff: **real-time collaborative software with zero coordination logic.**

---

## The Problem

Imagine two browsers editing the same list. A naive REST approach fails:

```
Browser A polls GET /tasks every second
Browser B polls GET /tasks every second
Browser A posts "Buy milk"
Browser B posts "Walk dog" (race condition — A's task may be overwritten)
```

Problems: race conditions, duplicates on retry, lost updates, polling lag, wasted requests.

What we need instead:

- Changes push instantly (no polling)
- Duplicate operations are harmless (idempotence)
- Operations can arrive in any order (commutativity)
- Conflicts resolve automatically (CRDT merge)

Sound familiar? You already built the conflict-resolution layer. Now add the transport.

---

## Architecture

```
┌─────────────┐         ┌─────────────┐         ┌─────────────┐
│  Browser A  │         │   Server    │         │  Browser B  │
│             │         │             │         │             │
│  UI Layer   │         │  TodoList   │         │  UI Layer   │
│             │◄───────►│   CRDT      │◄───────►│             │
│             │WebSocket│             │WebSocket│             │
└─────────────┘         └─────────────┘         └─────────────┘
```

The server holds one TodoList CRDT and a registry of connected browser clients. When a client sends an operation, the server applies it to the CRDT and broadcasts the updated state to everyone. Clients are thin: they send operations and render whatever state the server returns.

### Why hub-and-spoke?

CRDTs can run peer-to-peer, but a centralized server is easier to build, debug, and test:

| | Hub-and-spoke (this module) | Peer-to-peer |
|---|---|---|
| Connections | O(N) — each client → server | O(N²) — each client → all others |
| Discovery | Connect to one URL | NAT traversal, WebRTC setup |
| Testing | Simulate clients against localhost | Simulate network partitions |
| CRDT logic | Identical | Identical |

The merge logic you wrote transfers unchanged to P2P; only the transport changes. Libraries like [Automerge](https://automerge.org/) and [Yjs](https://yjs.dev/) use the same CRDT patterns over peer-to-peer transports.

---

## Protocol

Three message types, all JSON over WebSocket:

**1. `join` — client requests current state on connect**
```json
{"type": "join"}
```

**2. `state` — server sends full snapshot**
```json
{
  "type": "state",
  "tasks": [
    {"id": "server:1", "title": "Buy milk", "completed": false},
    {"id": "server:2", "title": "Walk dog", "completed": true}
  ]
}
```

**3. `op` — client or server sends an action**
```json
{"type": "op", "action": "add",          "title": "New task"}
{"type": "op", "action": "remove",       "id": "server:1"}
{"type": "op", "action": "setTitle",     "id": "server:1", "title": "Buy oat milk"}
{"type": "op", "action": "setCompleted", "id": "server:1", "completed": true}
```

### How operations map to your CRDT

Every operation message calls exactly one TodoList method. The server is a thin translation layer:

| Op action | TodoList method called | What happens in the CRDT |
|---|---|---|
| `add` | `AddTask(title)` | New unique ID added to OR-Set; LWW-Registers created for title and completed |
| `remove` | `RemoveTask(id)` | ID removed from OR-Set; registers become orphans |
| `setTitle` | `SetTitle(id, title)` | LWW-Register for title updated with current timestamp |
| `setCompleted` | `SetCompleted(id, done)` | LWW-Register for completed updated with current timestamp |

The server never resolves conflicts. Your Module 04 code does that automatically.

### Sequence diagram

```
Browser A          Server               Browser B
   |                  |                     |
   |──── join ───────►|                     |
   |◄─── state ───────|                     |
   |                  |◄──── join ──────────|
   |                  |───── state ────────►|
   |                  |                     |
   |──── op(add) ────►|                     |
   |◄─── state ───────|─────── state ──────►|  (broadcast to all)
   |                  |◄──── op(complete) ──|
   |◄─── state ───────|─────── state ──────►|  (broadcast to all)
```

---

## Getting Started

```bash
cd module-05-sync-service

# Copy all starter files (Go, JS, static assets, and the provided cmd/main.go)
cp -r ../starters/module-05-sync-service/* .

# Rename .tmpl files to their real extensions (handles subdirectories like cmd/ and static/)
find . -name "*.tmpl" | while read f; do mv "$f" "${f%.tmpl}"; done
```

This gives you:
- `server.go` — Hub struct and WebSocket handler (with TODOs)
- `protocol.go` — Message types and Apply/Encode/Decode (with TODOs)
- `server_test.go` — integration test suite (with TODOs)
- `static/app.js` — browser WebSocket client (with TODOs)
- `static/index.html`, `static/style.css` — complete, no changes needed
- `cmd/main.go` — HTTP server entry point, complete (no changes needed)

**Add the package declaration** (`package syncservice`) to `server.go`, `protocol.go`, and `server_test.go` before you start.

### Need to reset?

```bash
rm -f *.go static/app.js cmd/main.go
cp -r ../starters/module-05-sync-service/* .
find . -name "*.tmpl" | while read f; do mv "$f" "${f%.tmpl}"; done
```

---

## The Challenge

### File 1: `protocol.go`

Implement the message layer:

- `Message` and `Task` structs (already defined — just add the package)
- `Apply(todos *todolist.TodoList) string` — dispatches an op message to the right TodoList method; returns the new task ID for "add", empty string otherwise
- `StateMessage(todos *todolist.TodoList) Message` — builds a TypeState message from the current task list
- `Encode(m Message) ([]byte, error)` — JSON serialization
- `Decode(data []byte) (Message, error)` — JSON deserialization

### File 2: `server.go`

Implement the WebSocket hub:

```
Hub
├── todos     *todolist.TodoList   — the shared CRDT
├── clients   map[*websocket.Conn]bool   — connected browsers
└── mu        sync.Mutex           — protects todos and clients

NewHub() *Hub
ServeHTTP(w, r)   — upgrades to WebSocket, registers client, runs message loop
sendStateTo(conn) — sends a state snapshot to one connection
broadcastTo(msg, conns) — sends a message to a slice of connections
```

**Hub implements `http.Handler`** (it has `ServeHTTP`), so you can pass it directly to `http.Handle`.

**Message loop logic** (in `ServeHTTP`):
1. Upgrade HTTP → WebSocket
2. Register client, defer cleanup (remove from map, close connection)
3. Send current state to the new client
4. Loop: read message → if `join`, send state; if `op`, apply to TodoList and broadcast

**Concurrency note:** multiple browser goroutines share `h.todos` and `h.clients`. Use the mutex when reading or writing either. To avoid holding the lock during `WriteMessage` calls (which can block), snapshot the client list under the lock, then write outside it.

### File 3: `static/app.js`

Implement the browser WebSocket client:

- Open a connection to `ws://<host>/ws`
- On open: send a `join` message
- On message: if `type === "state"`, call `renderTasks`
- `renderTasks(tasks)` — rebuild the task list DOM
- `addTask()`, `removeTask(id)`, `setCompleted(id, done)`, `editTitle(id, span)` — send the matching `op` message

---

## Implementation Phases

Work through these in order. Each phase has tests.

### Phase 1: Protocol layer (~30 min)

Implement `protocol.go` completely before touching `server.go`. Start here because `server.go` depends on it.

Tests: the server tests will exercise protocol indirectly; you can also test `Encode`/`Decode` manually.

### Phase 2: Echo server (~20 min)

In `server.go`: upgrade HTTP to WebSocket, register the client, read messages in a loop, and echo each message back. Don't worry about the TodoList yet.

Test: `TestWebSocketUpgrade`

### Phase 3: State broadcast (~30 min)

Add the TodoList to Hub. On each new connection, send the current state. On `TypeJoin`, send state again.

Test: `TestStateSync`, `TestStateSyncAfterJoin`

### Phase 4: Operation handling (~40 min)

On `TypeOp`: call `msg.Apply(h.todos)`, build a state message, broadcast to all clients.

Test: `TestOperationBroadcast`, `TestConcurrentOperations`

### Phase 5: Reconnection (~20 min)

When a client disconnects, remove it from the map. State persists in `h.todos`, so reconnecting clients get the current state automatically — the earlier phases already handle this.

Test: `TestReconnection`

### Phase 6: Browser client (~60 min)

Implement `static/app.js`. Run the server (`go run ./cmd/`) and test manually with two browser tabs.

---

## Running the Server

```bash
# From module-05-sync-service/
go run ./cmd/

# Open two browser tabs
open http://localhost:8080
```

Add a task in one tab — it should appear in the other within milliseconds.

---

## Running Tests

```bash
# Run a single phase's tests
go test -run TestWebSocketUpgrade
go test -run TestStateSync
go test -run TestOperationBroadcast

# Run everything
go test .
```

---

## Hints

<details>
<summary><strong>Hint 1:</strong> How do I upgrade an HTTP connection to WebSocket?</summary>

Create a `websocket.Upgrader` at the package level, then call its `Upgrade` method inside your handler. The upgrader handles the HTTP handshake. The returned `*websocket.Conn` gives you `ReadMessage` and `WriteMessage`.

Set `CheckOrigin` to a function that returns `true` to allow connections from any origin (fine for local development).

</details>

<details>
<summary><strong>Hint 2:</strong> How do I send to all clients without holding the lock during WriteMessage?</summary>

`WriteMessage` can block if the client's TCP buffer is full. Holding the mutex during a blocking call is a deadlock risk.

Safe pattern: while holding the mutex, copy the client map keys into a slice. Then release the mutex and iterate the slice to call `WriteMessage`. Your `broadcastTo` method takes a slice of connections for exactly this reason.

</details>

<details>
<summary><strong>Hint 3:</strong> My test connects but never receives a message — why?</summary>

The test server URL starts with `http://`. WebSocket URLs start with `ws://`. Replace the scheme when dialing: `"ws" + strings.TrimPrefix(srv.URL, "http")`.

Also check that you have a read deadline set on the test connection — without it, `ReadMessage` blocks forever if no message arrives.

</details>

<details>
<summary><strong>Hint 4:</strong> How do I handle concurrent reads and writes on the same connection?</summary>

gorilla/websocket requires that only one goroutine calls `WriteMessage` at a time. In this server, each connection has exactly one goroutine running its read loop (the `ServeHTTP` goroutine). Other goroutines write via `broadcast`.

For a single-room server with low concurrency this rarely causes issues in practice, but if you see "concurrent write" panics, protect `WriteMessage` with a per-connection mutex.

</details>

<details>
<summary><strong>Hint 5:</strong> My JavaScript client doesn't show updates from other tabs</summary>

Check that `ws.onmessage` is set before the connection opens. Verify the JSON you're parsing has a `tasks` field (not `undefined`). Use `console.log(msg)` to inspect what the server sends.

</details>

---

## Extensions

1. **Multi-room support** — add a `/ws?room=teamA` URL parameter and maintain a separate `Hub` (and TodoList) per room
2. **Persistence** — serialize the TodoList state to a JSON file after each operation; reload on server startup
3. **Presence** — broadcast "N users online" when clients connect or disconnect
4. **Conflict visualization** — when a `setTitle` operation overwrites a local edit, display "Your change was overridden"
5. **Delta sync** — instead of sending the full state on each op, send only the changed fields (version vectors track what each client has seen)
6. **Operation log** — record all operations with timestamps; add an "undo last N ops" feature

---

## Compare with the Solution

```bash
cat ../solutions/module-05-sync-service/server_solution.go
cd ../solutions/module-05-sync-service && go test .
```

---

## Congratulations

You built a real-time collaborative application using CRDTs from first principles:

- **Module 01** G-Counter — monotone counting, merge by max
- **Module 02** PN-Counter — increment and decrement via two G-Counters
- **Module 03** LWW-Register + OR-Set — last-write-wins values, add-wins sets
- **Module 04** TodoList — composing CRDTs into a collaborative data structure
- **Module 05** Sync Service — exposing your CRDT over a real network

The merge logic you wrote in Go is the same logic that powers production collaborative tools. If you want to go further, explore [Automerge](https://automerge.org/), [Yjs](https://yjs.dev/), or academic papers on CRDTs (Shapiro et al., "A Comprehensive Study of Convergent and Commutative Replicated Data Types").
