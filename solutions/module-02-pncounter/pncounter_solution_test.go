// Add package declaration here (package pncounter)
package pncountersolution

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
	pnc := NewPNCounter("A")
	for i := 0; i < 5; i++ {
		pnc.Increment()
	}
	if pnc.Value() != 5 {
		t.Errorf("Expected value 5, got %d", pnc.Value())
	}
}

// Should return correct value after decrements
func TestPNCounterDecrement(t *testing.T) {
	pnc := NewPNCounter("A")
	for i := 0; i < 3; i++ {
		pnc.Decrement()
	}
	if pnc.Value() != -3 {
		t.Errorf("Expected value -3, got %d", pnc.Value())
	}
}

// Should handle both increments and decrements
func TestPNCounterIncrementAndDecrement(t *testing.T) {
	pnc := NewPNCounter("A")
	for i := 0; i < 10; i++ {
		pnc.Increment()
	}
	for i := 0; i < 3; i++ {
		pnc.Decrement()
	}
	if pnc.Value() != 7 {
		t.Errorf("Expected value 7, got %d", pnc.Value())
	}
}

// Should handle increments and decrements from multiple replicas after merge
func TestPNCounterMerge(t *testing.T) {
	pncA := NewPNCounter("A")
	pncB := NewPNCounter("B")

	// Replica A: increment 5 times
	for i := 0; i < 5; i++ {
		pncA.Increment()
	}

	// Replica B: decrement 2 times
	for i := 0; i < 2; i++ {
		pncB.Decrement()
	}

	// Merge B into A
	pncA.Merge(pncB)

	if pncA.Value() != 3 {
		t.Errorf("Expected value 3, got %d", pncA.Value())
	}
}

// Should be idempotent (merging same state twice has no effect)
func TestPNCounterMergeIdempotent(t *testing.T) {
	pncA := NewPNCounter("A")
	pncB := NewPNCounter("B")

	pncA.Increment()
	pncA.Increment()
	pncB.Decrement()

	// First merge
	pncA.Merge(pncB)
	value1 := pncA.Value()

	// Second merge (should have no effect)
	pncA.Merge(pncB)
	value2 := pncA.Value()

	if value1 != value2 {
		t.Errorf("Merge not idempotent: first=%d, second=%d", value1, value2)
	}

	if value1 != 1 {
		t.Errorf("Expected value 1, got %d", value1)
	}
}

// Should be commutative (order of merge doesn't matter)
func TestPNCounterMergeCommutative(t *testing.T) {
	// Scenario 1: A ← B ← C
	pncA1 := NewPNCounter("A")
	pncB1 := NewPNCounter("B")
	pncC1 := NewPNCounter("C")

	pncA1.Increment()
	pncA1.Increment()
	pncB1.Decrement()
	pncC1.Increment()
	pncC1.Increment()
	pncC1.Increment()

	pncA1.Merge(pncB1)
	pncA1.Merge(pncC1)
	value1 := pncA1.Value()

	// Scenario 2: A ← C ← B (different order)
	pncA2 := NewPNCounter("A")
	pncB2 := NewPNCounter("B")
	pncC2 := NewPNCounter("C")

	pncA2.Increment()
	pncA2.Increment()
	pncB2.Decrement()
	pncC2.Increment()
	pncC2.Increment()
	pncC2.Increment()

	pncA2.Merge(pncC2)
	pncA2.Merge(pncB2)
	value2 := pncA2.Value()

	if value1 != value2 {
		t.Errorf("Merge not commutative: order1=%d, order2=%d", value1, value2)
	}

	if value1 != 4 {
		t.Errorf("Expected value 4 (2 + 3 - 1), got %d", value1)
	}
}

// Should handle negative values
func TestPNCounterNegativeValue(t *testing.T) {
	pnc := NewPNCounter("A")
	for i := 0; i < 5; i++ {
		pnc.Decrement()
	}
	if pnc.Value() >= 0 {
		t.Errorf("Expected negative value, got %d", pnc.Value())
	}
	if pnc.Value() != -5 {
		t.Errorf("Expected value -5, got %d", pnc.Value())
	}
}
