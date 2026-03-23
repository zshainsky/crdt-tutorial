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
	// Create a PNCounter, increment 5 times, assert Value() == 5
	pnc := NewPNCounter("A")

	for range 5 {
		pnc.Increment()
	}
	if pnc.Value() != 5 {
		t.Errorf("Value should be 5, got %d", pnc.Value())
	}
}

// Should return correct value after decrements
func TestPNCounterDecrement(t *testing.T) {
	// Create a PNCounter, decrement 3 times, assert Value() == -3
	pnc := NewPNCounter("A")

	for range 3 {
		pnc.Decrement()
	}
	if pnc.Value() != -3 {
		t.Errorf("Value should be -3, got %d", pnc.Value())
	}
}

// Should handle both increments and decrements
func TestPNCounterIncrementAndDecrement(t *testing.T) {
	// Create a PNCounter
	// Increment 10 times
	// Decrement 3 times
	// Assert Value() == 7
	pnc := NewPNCounter("A")

	for range 10 {
		pnc.Increment()
	}
	for range 3 {
		pnc.Decrement()
	}
	if pnc.Value() != 7 {
		t.Errorf("Value should be 7, got %d", pnc.Value())
	}
}

// Should handle increments and decrements from multiple replicas after merge
func TestPNCounterMerge(t *testing.T) {
	// Create two PNCounters (different replica IDs)
	// Replica A: increment 5 times
	// Replica B: decrement 2 times
	// Merge B into A
	// Assert A.Value() == 3
	pncA := NewPNCounter("A")
	for range 5 {
		pncA.Increment()
	}
	pncB := NewPNCounter("B")
	for range 2 {
		pncB.Decrement()
	}
	pncA.Merge(pncB)
	if pncA.Value() != 3 {
		t.Errorf("Value should be 3, got %d", pncA.Value())
	}

}

// Should be idempotent (merging same state twice has no effect)
func TestPNCounterMergeIdempotent(t *testing.T) {
	// Create two PNCounters
	// Increment/decrement them
	// Merge B into A
	// Save A's value
	// Merge B into A again
	// Assert value didn't change
	pncA := NewPNCounter("A")
	for range 5 {
		pncA.Increment()
	}
	pncB := NewPNCounter("B")
	for range 2 {
		pncB.Decrement()
	}

	pncA.Merge(pncB)
	pncAVal := pncA.Value()

	pncA.Merge(pncB)
	pncANewVal := pncA.Value()

	if pncAVal != pncANewVal {
		t.Errorf("Values should be %d, got %d", pncAVal, pncANewVal)
	}

}

// Should be commutative (order of merge doesn't matter)
func TestPNCounterMergeCommutative(t *testing.T) {
	// TODO: Create three PNCounters A, B, C
	// A: increment 2 times
	// B: decrement 1 time
	// C: increment 3 times
	// Merge in different orders: A←B←C vs A←C←B
	// Assert both result in same Value()

	pncA1 := NewPNCounter("A")
	pncA1.Increment()
	pncA1.Increment()

	pncB1 := NewPNCounter("B")
	pncB1.Decrement()

	pncC1 := NewPNCounter("C")
	pncC1.Increment()
	pncC1.Increment()
	pncC1.Increment()

	pncB1.Merge(pncC1)
	pncA1.Merge(pncB1)
	value1 := pncA1.Value()

	pncA2 := NewPNCounter("A")
	pncA2.Increment()
	pncA2.Increment()

	pncB2 := NewPNCounter("B")
	pncB2.Decrement()

	pncC2 := NewPNCounter("C")
	pncC2.Increment()
	pncC2.Increment()
	pncC2.Increment()

	pncC2.Merge(pncB2)
	pncA2.Merge(pncC2)
	value2 := pncA2.Value()

	if value1 != value2 {
		t.Errorf("Values should be %d, got %d", value1, value2)
	}
}

// Should handle negative values
func TestPNCounterNegativeValue(t *testing.T) {
	// TODO: Create a PNCounter
	// Decrement 5 times (or increment less than decrement)
	// Assert Value() is negative
	pnc := NewPNCounter("A")
	for range 5 {
		pnc.Decrement()
	}
	if pnc.Value() != -5 {
		t.Errorf("Value should be negative, got %d", pnc.Value())
	}
}
