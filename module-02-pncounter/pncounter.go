// Add package declaration here (package pncounter)
package pncounter

// Add import statement:
import gcounter "github.com/zshainsky/crdt-tutorial/module-01-gcounter"

// PNCounter is a positive-negative counter that can increment and decrement.
// It uses two G-Counters: P (positive/increments) and N (negative/decrements).
// Value = P - N
type PNCounter struct {
	replicaID string
	// Add two fields:
	p *gcounter.GCounter // Positive counter (increments)
	n *gcounter.GCounter // Negative counter (decrements)
}

// NewPNCounter creates a new PN-Counter for the given replica.
func NewPNCounter(replicaID string) *PNCounter {
	// Initialize both P and N counters using gcounter.NewGCounter
	// Remember: both should use the same replicaID
	return &PNCounter{
		replicaID: replicaID,
		p:         gcounter.NewGCounter(replicaID),
		n:         gcounter.NewGCounter(replicaID),
	}
}

// Increment increases the counter value by 1.
func (pn *PNCounter) Increment() {
	// Increment the P counter
	pn.p.Increment()
}

// Decrement decreases the counter value by 1.
func (pn *PNCounter) Decrement() {
	// Increment the N counter (yes, increment!)
	// We're counting "how many times we decremented"
	pn.n.Increment()
}

// Value returns the current value of the counter (P - N).
func (pn *PNCounter) Value() int {
	// Return P.Value() - N.Value()
	return pn.p.Value() - pn.n.Value()
}

// Merge combines another PN-Counter's state into this one.
// Merges both the P and N counters.
func (pn *PNCounter) Merge(other *PNCounter) {
	// Merge both P and N counters
	// pn.p.Merge(other.p)
	// pn.n.Merge(other.n)
	pn.p.Merge(other.p)
	pn.n.Merge(other.n)
}
