package gcountersolution

import (
	"testing"
)

func TestGCounterInitialValue(t *testing.T) {
	gc := NewGCounter("A")
	if gc.Value() != 0 {
		t.Errorf("Expected initial value 0, got %d", gc.Value())
	}
}

func TestGCounterSingleIncrement(t *testing.T) {
	gc := NewGCounter("A")
	gc.Increment()
	if gc.Value() != 1 {
		t.Errorf("Expected value 1 after increment, got %d", gc.Value())
	}
}

func TestGCounterMultipleIncrements(t *testing.T) {
	gc := NewGCounter("A")
	for i := 0; i < 5; i++ {
		gc.Increment()
	}
	if gc.Value() != 5 {
		t.Errorf("Expected value 5 after 5 increments, got %d", gc.Value())
	}
}

func TestGCounterMergeTwoReplicas(t *testing.T) {
	gcA := NewGCounter("A")
	gcB := NewGCounter("B")

	// A increments 3 times
	for i := 0; i < 3; i++ {
		gcA.Increment()
	}

	// B increments 2 times
	for i := 0; i < 2; i++ {
		gcB.Increment()
	}

	// Merge B into A
	gcA.Merge(gcB)

	if gcA.Value() != 5 {
		t.Errorf("Expected value 5 after merge, got %d", gcA.Value())
	}
}

func TestGCounterMergeIdempotent(t *testing.T) {
	gcA := NewGCounter("A")
	gcB := NewGCounter("B")

	gcA.Increment()
	gcA.Increment()
	gcA.Increment()

	gcB.Increment()
	gcB.Increment()

	// First merge
	gcA.Merge(gcB)
	firstValue := gcA.Value()

	// Second merge (should have no effect)
	gcA.Merge(gcB)
	secondValue := gcA.Value()

	if firstValue != secondValue {
		t.Errorf("Merge is not idempotent: first=%d, second=%d", firstValue, secondValue)
	}

	if secondValue != 5 {
		t.Errorf("Expected value 5, got %d", secondValue)
	}
}

func TestGCounterMergeCommutative(t *testing.T) {
	// Scenario 1: A ← B ← C
	gcA1 := NewGCounter("A")
	gcB1 := NewGCounter("B")
	gcC1 := NewGCounter("C")

	for i := 0; i < 3; i++ {
		gcA1.Increment()
	}
	for i := 0; i < 2; i++ {
		gcB1.Increment()
	}
	for i := 0; i < 5; i++ {
		gcC1.Increment()
	}

	gcA1.Merge(gcB1)
	gcA1.Merge(gcC1)
	value1 := gcA1.Value()

	// Scenario 2: A ← C ← B
	gcA2 := NewGCounter("A")
	gcB2 := NewGCounter("B")
	gcC2 := NewGCounter("C")

	for i := 0; i < 3; i++ {
		gcA2.Increment()
	}
	for i := 0; i < 2; i++ {
		gcB2.Increment()
	}
	for i := 0; i < 5; i++ {
		gcC2.Increment()
	}

	gcA2.Merge(gcC2)
	gcA2.Merge(gcB2)
	value2 := gcA2.Value()

	if value1 != value2 {
		t.Errorf("Merge is not commutative: A←B←C=%d, A←C←B=%d", value1, value2)
	}

	if value1 != 10 {
		t.Errorf("Expected value 10, got %d", value1)
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

func TestGCounterMergeWithSelf(t *testing.T) {
	gc := NewGCounter("A")
	gc.Increment()
	gc.Increment()

	// Merge with self (should be idempotent)
	gc.Merge(gc)

	if gc.Value() != 2 {
		t.Errorf("Expected value 2 after self-merge, got %d", gc.Value())
	}
}

func TestGCounterThreeReplicas(t *testing.T) {
	gcA := NewGCounter("A")
	gcB := NewGCounter("B")
	gcC := NewGCounter("C")

	gcA.Increment()
	gcB.Increment()
	gcB.Increment()
	gcC.Increment()
	gcC.Increment()
	gcC.Increment()

	// Merge all into A
	gcA.Merge(gcB)
	gcA.Merge(gcC)

	if gcA.Value() != 6 {
		t.Errorf("Expected value 6, got %d", gcA.Value())
	}
}

func TestGCounterPartialMerge(t *testing.T) {
	// A has seen updates from A and B
	gcA := NewGCounter("A")
	gcB := NewGCounter("B")
	gcA.Increment()
	gcB.Increment()
	gcB.Increment()
	gcA.Merge(gcB)

	// C only has its own updates
	gcC := NewGCounter("C")
	gcC.Increment()
	gcC.Increment()
	gcC.Increment()

	// When A merges C, A should have all three replicas
	gcA.Merge(gcC)

	if gcA.Value() != 6 {
		t.Errorf("Expected value 6 (1 + 2 + 3), got %d", gcA.Value())
	}
}

func TestGCounterZeroValueReplica(t *testing.T) {
	gcA := NewGCounter("A")
	gcB := NewGCounter("B")

	// A increments
	gcA.Increment()
	gcA.Increment()

	// B doesn't increment (stays at 0)

	// Merge B into A (should have no effect since B=0)
	gcA.Merge(gcB)

	if gcA.Value() != 2 {
		t.Errorf("Expected value 2 after merging zero-value replica, got %d", gcA.Value())
	}
}
