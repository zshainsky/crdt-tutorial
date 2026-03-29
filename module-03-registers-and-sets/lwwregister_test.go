// Add package declaration (package registers)
package registers

import (
	"testing"
	"time"
)

// Should return empty string for a new register
func TestLWWRegisterInitialGet(t *testing.T) {
	// Create a register and assert Get() returns ""
	r := NewLWWRegister("A")
	if r.Get() != "" {
		t.Errorf("empty register should have empty value")
	}
}

// Should return the value after Set
func TestLWWRegisterSetGet(t *testing.T) {
	// Use SetAt to assign a value with a known timestamp, then assert Get() returns it
	r := NewLWWRegister("A")
	tNow := time.Now().Unix()

	want := "now"
	r.SetAt(want, tNow)

	if r.Get() != want {
		t.Errorf("got: %s, want: %s", r.Get(), want)
	}
}

// Should keep the higher-timestamp value after merge
func TestLWWRegisterMergeHigherTimestampWins(t *testing.T) {
	// Give register A an older timestamp, register B a newer one
	// Merge B into A — assert A now holds B's value
	A := NewLWWRegister("A")
	B := NewLWWRegister("B")

	tA := time.Date(2026, time.January, 1, 0, 0, 0, 0, &time.Location{}).Unix()
	A.SetAt("A Value", tA)
	if A.Get() != "A Value" {
		t.Errorf("got: %s, want: %s", A.Get(), "A Value")
	}

	tB := time.Date(2026, time.December, 1, 0, 0, 0, 0, &time.Location{}).Unix()
	B.SetAt("B Value", tB)
	if B.Get() != "B Value" {
		t.Errorf("got: %s, want: %s", B.Get(), "B Value")
	}

	// Expect A now holds B's value
	A.Merge(B)
	if A.Get() != "B Value" {
		t.Errorf("got: %s, want: %s", A.Get(), "B Value")
	}

}

// Should ignore a lower-timestamp value during merge
func TestLWWRegisterMergeLowerTimestampIgnored(t *testing.T) {
	// Give register A a newer timestamp, register B an older one
	// Merge B into A — assert A keeps its own value
	A := NewLWWRegister("A")
	B := NewLWWRegister("B")

	tA := time.Date(2026, time.December, 1, 0, 0, 0, 0, &time.Location{}).Unix()
	A.SetAt("A Value", tA)
	if A.Get() != "A Value" {
		t.Errorf("got: %s, want: %s", A.Get(), "A Value")
	}

	tB := time.Date(2026, time.January, 1, 0, 0, 0, 0, &time.Location{}).Unix()
	B.SetAt("B Value", tB)
	if B.Get() != "B Value" {
		t.Errorf("got: %s, want: %s", B.Get(), "B Value")
	}

	// Expect A now holds B's value
	A.Merge(B)
	if A.Get() != "A Value" {
		t.Errorf("got: %s, want: %s", A.Get(), "A Value")
	}
}

// Should be idempotent (merging same state twice has no effect)
func TestLWWRegisterMergeIdempotent(t *testing.T) {
	// Merge B into A twice — assert the value is unchanged after the second merge
	A := NewLWWRegister("A")
	B := NewLWWRegister("B")

	tA := time.Date(2026, time.January, 1, 0, 0, 0, 0, &time.Location{}).Unix()
	A.SetAt("A Value", tA)
	if A.Get() != "A Value" {
		t.Errorf("got: %s, want: %s", A.Get(), "A Value")
	}

	tB := time.Date(2026, time.December, 1, 0, 0, 0, 0, &time.Location{}).Unix()
	B.SetAt("B Value", tB)
	if B.Get() != "B Value" {
		t.Errorf("got: %s, want: %s", B.Get(), "B Value")
	}

	A.Merge(B)
	A.Merge(B)
	if A.Get() != "B Value" {
		t.Errorf("got: %s, want: %s", A.Get(), "B Value")
	}
}

// Should be commutative (A absorbs B and B absorbs A both reach the same value)
func TestLWWRegisterMergeCommutative(t *testing.T) {
	// TODO: Set up two registers with different timestamps
	// Merge in both directions (A←B and B←A separately)
	// Assert both results are the same value
	A := NewLWWRegister("A")
	B := NewLWWRegister("B")

	tA := time.Date(2026, time.January, 1, 0, 0, 0, 0, &time.Location{}).Unix()
	A.SetAt("A Value", tA)
	if A.Get() != "A Value" {
		t.Errorf("got: %s, want: %s", A.Get(), "A Value")
	}

	tB := time.Date(2026, time.December, 1, 0, 0, 0, 0, &time.Location{}).Unix()
	B.SetAt("B Value", tB)
	if B.Get() != "B Value" {
		t.Errorf("got: %s, want: %s", B.Get(), "B Value")
	}

	A.Merge(B)
	if A.Get() != "B Value" {
		t.Errorf("got: %s, want: %s", A.Get(), "B Value")
	}

	B.Merge(A)
	if B.Get() != "B Value" {
		t.Errorf("got: %s, want: %s", A.Get(), "B Value")
	}
}

// Should keep current value when timestamps are equal (tie-break)
func TestLWWRegisterTieBreak(t *testing.T) {
	// Give both registers the same timestamp but different values
	// Merge B into A — assert A keeps its own value
	A := NewLWWRegister("A")
	B := NewLWWRegister("B")

	tA := time.Date(2026, time.January, 1, 0, 0, 0, 0, &time.Location{}).Unix()
	A.SetAt("A Value", tA)
	if A.Get() != "A Value" {
		t.Errorf("got: %s, want: %s", A.Get(), "A Value")
	}

	tB := time.Date(2026, time.January, 1, 0, 0, 0, 0, &time.Location{}).Unix()
	B.SetAt("B Value", tB)
	if B.Get() != "B Value" {
		t.Errorf("got: %s, want: %s", B.Get(), "B Value")
	}

	// A keeps its own value
	A.Merge(B)
	if A.Get() != "A Value" {
		t.Errorf("got: %s, want: %s", A.Get(), "A Value")
	}
}
