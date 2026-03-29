package registerssolution

import (
	"reflect"
	"testing"
)

func TestORSetContainsAfterAdd(t *testing.T) {
	s := NewORSet("A")
	s.Add("apple")
	if !s.Contains("apple") {
		t.Error("expected apple to be in set after Add")
	}
}

func TestORSetNotContainsAfterRemove(t *testing.T) {
	s := NewORSet("A")
	s.Add("apple")
	s.Remove("apple")
	if s.Contains("apple") {
		t.Error("expected apple to not be in set after Remove")
	}
}

func TestORSetElementsEmpty(t *testing.T) {
	s := NewORSet("A")
	if len(s.Elements()) != 0 {
		t.Errorf("expected empty elements, got %v", s.Elements())
	}
}

func TestORSetElementsSorted(t *testing.T) {
	s := NewORSet("A")
	s.Add("banana")
	s.Add("apple")
	s.Add("cherry")
	expected := []string{"apple", "banana", "cherry"}
	if !reflect.DeepEqual(s.Elements(), expected) {
		t.Errorf("expected %v, got %v", expected, s.Elements())
	}
}

func TestORSetMergeCombinesElements(t *testing.T) {
	a := NewORSet("A")
	b := NewORSet("B")
	a.Add("apple")
	b.Add("banana")
	a.Merge(b)
	if !a.Contains("apple") || !a.Contains("banana") {
		t.Errorf("expected both elements after merge, got %v", a.Elements())
	}
}

func TestORSetMergeIdempotent(t *testing.T) {
	a := NewORSet("A")
	b := NewORSet("B")
	a.Add("apple")
	b.Add("banana")
	a.Merge(b)
	elems1 := a.Elements()
	a.Merge(b)
	elems2 := a.Elements()
	if !reflect.DeepEqual(elems1, elems2) {
		t.Errorf("merge not idempotent: %v vs %v", elems1, elems2)
	}
}

func TestORSetMergeCommutative(t *testing.T) {
	// Scenario 1: A absorbs B
	a1 := NewORSet("A")
	b1 := NewORSet("B")
	a1.Add("apple")
	b1.Add("banana")
	a1.Merge(b1)

	// Scenario 2: B absorbs A
	a2 := NewORSet("A")
	b2 := NewORSet("B")
	a2.Add("apple")
	b2.Add("banana")
	b2.Merge(a2)

	if !reflect.DeepEqual(a1.Elements(), b2.Elements()) {
		t.Errorf("merge not commutative: %v vs %v", a1.Elements(), b2.Elements())
	}
}

func TestORSetConcurrentAddRemoveAddWins(t *testing.T) {
	// A adds "apple" — tag A:1
	// B adds "apple" — tag B:1 (concurrent, A doesn't know yet)
	// A removes "apple" — tombstones A:1 only (B:1 not yet observed by A)
	// After merge: B:1 is still live, so apple remains in the set
	a := NewORSet("A")
	b := NewORSet("B")
	a.Add("apple") // tag A:1
	b.Add("apple") // tag B:1
	a.Remove("apple")
	a.Merge(b)
	if !a.Contains("apple") {
		t.Error("expected apple in set: concurrent add should survive concurrent remove (add-wins)")
	}
}

func TestORSetReAddAfterRemove(t *testing.T) {
	s := NewORSet("A")
	s.Add("apple")
	s.Remove("apple")
	s.Add("apple") // generates a new tag — element lives again
	if !s.Contains("apple") {
		t.Error("expected apple to be in set after re-add")
	}
}

func TestORSetRemoveNonExistent(t *testing.T) {
	s := NewORSet("A")
	s.Remove("ghost") // should not panic or error
	if s.Contains("ghost") {
		t.Error("expected ghost to not be in set")
	}
}
