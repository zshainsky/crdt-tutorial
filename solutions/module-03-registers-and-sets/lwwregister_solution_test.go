package registerssolution

import "testing"

func TestLWWRegisterInitialGet(t *testing.T) {
	r := NewLWWRegister("A")
	if r.Get() != "" {
		t.Errorf("expected empty string, got %q", r.Get())
	}
}

func TestLWWRegisterSetGet(t *testing.T) {
	r := NewLWWRegister("A")
	r.SetAt("hello", 100)
	if r.Get() != "hello" {
		t.Errorf("expected %q, got %q", "hello", r.Get())
	}
}

func TestLWWRegisterMergeHigherTimestampWins(t *testing.T) {
	a := NewLWWRegister("A")
	b := NewLWWRegister("B")
	a.SetAt("old", 100)
	b.SetAt("new", 200)
	a.Merge(b)
	if a.Get() != "new" {
		t.Errorf("expected %q, got %q", "new", a.Get())
	}
}

func TestLWWRegisterMergeLowerTimestampIgnored(t *testing.T) {
	a := NewLWWRegister("A")
	b := NewLWWRegister("B")
	a.SetAt("current", 200)
	b.SetAt("old", 100)
	a.Merge(b)
	if a.Get() != "current" {
		t.Errorf("expected %q, got %q", "current", a.Get())
	}
}

func TestLWWRegisterMergeIdempotent(t *testing.T) {
	a := NewLWWRegister("A")
	b := NewLWWRegister("B")
	a.SetAt("hello", 100)
	b.SetAt("world", 200)
	a.Merge(b)
	v1 := a.Get()
	a.Merge(b)
	v2 := a.Get()
	if v1 != v2 {
		t.Errorf("merge not idempotent: %q vs %q", v1, v2)
	}
	if v1 != "world" {
		t.Errorf("expected %q, got %q", "world", v1)
	}
}

func TestLWWRegisterMergeCommutative(t *testing.T) {
	// Scenario 1: A absorbs B
	a1 := NewLWWRegister("A")
	b1 := NewLWWRegister("B")
	a1.SetAt("from-a", 100)
	b1.SetAt("from-b", 200)
	a1.Merge(b1)

	// Scenario 2: B absorbs A
	a2 := NewLWWRegister("A")
	b2 := NewLWWRegister("B")
	a2.SetAt("from-a", 100)
	b2.SetAt("from-b", 200)
	b2.Merge(a2)

	if a1.Get() != b2.Get() {
		t.Errorf("merge not commutative: a1=%q, b2=%q", a1.Get(), b2.Get())
	}
	if a1.Get() != "from-b" {
		t.Errorf("expected %q (higher timestamp wins), got %q", "from-b", a1.Get())
	}
}

func TestLWWRegisterTieBreak(t *testing.T) {
	// Equal timestamps: current value is preserved (deterministic tie-breaking)
	a := NewLWWRegister("A")
	b := NewLWWRegister("B")
	a.SetAt("from-a", 100)
	b.SetAt("from-b", 100)
	a.Merge(b)
	if a.Get() != "from-a" {
		t.Errorf("expected %q on equal-timestamp tie, got %q", "from-a", a.Get())
	}
}
