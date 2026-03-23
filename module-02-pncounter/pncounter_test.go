// Add package declaration here (package pncounter)
package pncounter

import "testing"

// Should return 0 for a new counter
func TestPNCounterInitialValue(t *testing.T) {
	// Create a PNCounter and assert Value() == 0
	pnc := NewPNCounter("A")
	if pnc == nil {
		t.Errorf("NewPNCounter should not be nil")
	}

	if pnc.Value() != 0 {
		t.Errorf("Value should be 0, got %d", pnc.Value())
	}
}

// Should return correct value after increments
func TestPNCounterIncrement(t *testing.T) {
	// TODO: Create a PNCounter, increment 5 times, assert Value() == 5
}

// Should return correct value after decrements
func TestPNCounterDecrement(t *testing.T) {
	// TODO: Create a PNCounter, decrement 3 times, assert Value() == -3
}

// Should handle both increments and decrements
func TestPNCounterIncrementAndDecrement(t *testing.T) {
	// TODO: Create a PNCounter
	// Increment 10 times
	// Decrement 3 times
	// Assert Value() == 7
}

// Should handle increments and decrements from multiple replicas after merge
func TestPNCounterMerge(t *testing.T) {
	// TODO: Create two PNCounters (different replica IDs)
	// Replica A: increment 5 times
	// Replica B: decrement 2 times
	// Merge B into A
	// Assert A.Value() == 3
}

// Should be idempotent (merging same state twice has no effect)
func TestPNCounterMergeIdempotent(t *testing.T) {
	// TODO: Create two PNCounters
	// Increment/decrement them
	// Merge B into A
	// Save A's value
	// Merge B into A again
	// Assert value didn't change
}

// Should be commutative (order of merge doesn't matter)
func TestPNCounterMergeCommutative(t *testing.T) {
	// TODO: Create three PNCounters A, B, C
	// A: increment 2 times
	// B: decrement 1 time
	// C: increment 3 times
	// Merge in different orders: A←B←C vs A←C←B
	// Assert both result in same Value()
}

// Should handle negative values
func TestPNCounterNegativeValue(t *testing.T) {
	// TODO: Create a PNCounter
	// Decrement 5 times (or increment less than decrement)
	// Assert Value() is negative
}
