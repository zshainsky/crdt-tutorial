package gcounter

import (
	"testing"
)

// Should return 0 for a new counter
func TestGCounterInitialValue(t *testing.T) {
	// Create a GCounter and assert Value() == 0
	gc := NewGCounter("A")
	if gc == nil {
		t.Errorf("GCounter should not be nil")
	}
	if gc.Value() != 0 {
		t.Errorf("Expected value 0, got %d", gc.Value())
	}
}

// Should return correct value after increments
func TestGCounterIncrement(t *testing.T) {
	// TODO: Create a GCounter, increment 5 times, assert Value() == 5
	gc := NewGCounter("A")
	if gc == nil {
		t.Errorf("GCounter should not be nil")
	}

	want := 5
	for range want {
		gc.Increment()
	}
	if gc.Value() != want {
		t.Errorf("Expected value %d, got %d", want, gc.Value())
	}

}

// Should handle increments from multiple replicas after merge
func TestGCounterMerge(t *testing.T) {
	// Create two GCounters (different replica IDs)
	// Increment each a few times
	// Merge one into the other
	// Assert Value() equals the sum of both increments

	A := NewGCounter("A")
	aCount := 5
	for range aCount {
		A.Increment()
	}

	B := NewGCounter("B")
	bCount := 100
	for range bCount {
		B.Increment()
	}

	// Merge B into A
	A.Merge(B)
	want := aCount + bCount
	if A.Value() != want {
		t.Errorf("Expected value %d, got %d", want, A.Value())
	}
}

// Should be idempotent (merging same state twice has no effect)
func TestGCounterMergeIdempotent(t *testing.T) {
	// Create two GCounters, increment them
	// Merge B into A
	// Merge B into A again
	// Assert Value() doesn't change after second merge
	A := NewGCounter("A")
	B := NewGCounter("B")
	A.Increment()
	B.Increment()

	A.Merge(B)
	got := A.Value()

	A.Merge(B)
	want := A.Value()

	if got != want {
		t.Errorf("Expectd value %d does not equal expected value %d", got, want)
	}

}

// Should be commutative (order of merge doesn't matter)
func TestGCounterMergeCommutative(t *testing.T) {
	// create three GCounters A, B, C
	// Increment each differently
	// Merge in different orders: A←B←C vs A←C←B
	// Assert both result in same Value()

	// A ← B, then A ← C
	A1 := NewGCounter("A")
	B1 := NewGCounter("B")
	C1 := NewGCounter("C")

	// Increment each
	for range 1 {
		A1.Increment()
	}
	for range 5 {
		B1.Increment()
	}
	for range 3 {
		C1.Increment()
	}

	A1.Merge(B1) // A ← B first
	A1.Merge(C1) // then A ← C
	got := A1.Value()

	// A ← C, then A ← B (FRESH counters!)
	A2 := NewGCounter("A")
	B2 := NewGCounter("B")
	C2 := NewGCounter("C")

	// Increment the same way
	for range 1 {
		A2.Increment()
	}
	for range 5 {
		B2.Increment()
	}
	for range 3 {
		C2.Increment()
	}

	A2.Merge(C2) // A ← C first (different order!)
	A2.Merge(B2) // then A ← B
	want := A2.Value()

	// Now compare
	if got != want {
		t.Errorf("Not commutative: order1=%d, order2=%d", got, want)
	}

	if got != 1+5+3 {
		t.Errorf("Value %d does not equal 1+5+3", got)
	}
}

func TestGCounterMergeAssociative(t *testing.T) {
	// Scenario 1: (A ← B) ← C
	gcA1 := NewGCounter("A")
	gcB1 := NewGCounter("B")
	gcC1 := NewGCounter("C")

	gcA1.Increment()
	gcB1.Increment()
	gcB1.Increment()
	gcC1.Increment()
	gcC1.Increment()
	gcC1.Increment()

	gcA1.Merge(gcB1)
	gcA1.Merge(gcC1)
	value1 := gcA1.Value()

	// Scenario 2: A ← (B ← C)
	gcA2 := NewGCounter("A")
	gcB2 := NewGCounter("B")
	gcC2 := NewGCounter("C")

	gcA2.Increment()
	gcB2.Increment()
	gcB2.Increment()
	gcC2.Increment()
	gcC2.Increment()
	gcC2.Increment()

	gcB2.Merge(gcC2)
	gcA2.Merge(gcB2)
	value2 := gcA2.Value()

	if value1 != value2 {
		t.Errorf("Merge is not associative: (A←B)←C=%d, A←(B←C)=%d", value1, value2)
	}

	if value1 != 6 {
		t.Errorf("Expected value 6, got %d", value1)
	}
}
