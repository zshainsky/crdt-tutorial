package gcounter

// GCounter is a grow-only counter CRDT.
// Each replica tracks increments in a vector; merge takes max per replica.
type GCounter struct {
	replicaID string
	counts    map[string]int
}

// NewGCounter creates a new G-Counter for the given replica.
func NewGCounter(replicaID string) *GCounter {
	return &GCounter{
		replicaID: replicaID,
		counts:    make(map[string]int),
	}
}

// Increment increases this replica's count by 1.
func (g *GCounter) Increment() {
	g.counts[g.replicaID]++
}

// Value returns the sum of all replica counts.
func (g *GCounter) Value() int {
	total := 0
	for _, count := range g.counts {
		total += count
	}
	return total
}

// Merge combines another G-Counter's state into this one.
// Takes the max of each replica's count to ensure convergence.
func (g *GCounter) Merge(other *GCounter) {
	for replicaID, count := range other.counts {
		if count > g.counts[replicaID] {
			g.counts[replicaID] = count
		}
	}
}
