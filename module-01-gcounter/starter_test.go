package gcounter
package gcounter

import "testing"

// Should return 0 for a new counter
func TestGCounterInitialValue(t *testing.T) {
	// TODO: Create a GCounter and assert Value() == 0
}

// Should return correct value after increments
func TestGCounterIncrement(t *testing.T) {
	// TODO: Create a GCounter, increment 5 times, assert Value() == 5

























}	// Assert both result in same Value()	// Merge in different orders: A←B←C vs A←C←B	// Increment each differently	// TODO: Create three GCounters A, B, Cfunc TestGCounterMergeCommutative(t *testing.T) {// Should be commutative (order of merge doesn't matter)}	// Assert Value() doesn't change after second merge	// Merge B into A again	// Merge B into A	// TODO: Create two GCounters, increment themfunc TestGCounterMergeIdempotent(t *testing.T) {// Should be idempotent (merging same state twice has no effect)}	// Assert Value() equals the sum of both increments	// Merge one into the other	// Increment each a few times	// TODO: Create two GCounters (different replica IDs)func TestGCounterMerge(t *testing.T) {// Should handle increments from multiple replicas after merge}