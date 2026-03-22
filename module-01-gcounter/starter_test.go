package gcounter

import "testing"

// Should return 0 for a new counter
func TestGCounterInitialValue(t *testing.T) {
	// TODO: Create a GCounter and assert Value() == 0
}

// Should return correct value after increments
func TestGCounterIncrement(t *testing.T) {
	// TODO: Create a GCounter, increment 5 times, assert Value() == 5
}

// Should handle increments from multiple replicas after merge
func TestGCounterMerge(t *testing.T) {
	// TODO: Create two GCounters (different replica IDs)
	// Increment each a few times
	// Merge one into the other
	// Assert Value() equals the sum of both increments
}

// Should be idempotent (merging same state twice has no effect)
func TestGCounterMergeIdempotent(t *testing.T) {
	// TODO: Create two GCounters, increment them
	// Merge B into A
	// Merge B into A again
	// Assert Value() doesn't change after second merge
}

// Should be commutative (order of merge doesn't matter)
func TestGCounterMergeCommutative(t *testing.T) {
	// TODO: Create three GCounters A, B, C
	// Increment each differently
	// Merge in different orders: A←B←C vs A←C←B
	// Assert both result in same Value()
}
