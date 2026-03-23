package pncountersolution

import "github.com/zshainsky/crdt-tutorial/module-01-gcounter"

// PNCounter is a positive-negative counter CRDT.
// It uses two G-Counters: P for increments and N for decrements.
// Value = P - N
type PNCounter struct {
	replicaID string
	p         *gcounter.GCounter
	n         *gcounter.GCounter
}

// NewPNCounter creates a new PN-Counter for the given replica.
func NewPNCounter(replicaID string) *PNCounter {
	return &PNCounter{
		replicaID: replicaID,
		p:         gcounter.NewGCounter(replicaID),
		n:         gcounter.NewGCounter(replicaID),
	}
}

// Increment increases the counter value by 1.
func (pn *PNCounter) Increment() {
	pn.p.Increment()
}

// Decrement decreases the counter value by 1.
func (pn *PNCounter) Decrement() {
	pn.n.Increment()
}

// Value returns the current value of the counter (P - N).
func (pn *PNCounter) Value() int {
	return pn.p.Value() - pn.n.Value()
}

// Merge combines another PN-Counter's state into this one.
func (pn *PNCounter) Merge(other *PNCounter) {
	pn.p.Merge(other.p)
	pn.n.Merge(other.n)
}
