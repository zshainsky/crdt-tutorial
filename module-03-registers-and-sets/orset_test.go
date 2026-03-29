// TODO: Add package declaration (package registers)
package registers

import (
	"slices"
	"testing"
)

// Should contain element after Add
func TestORSetContainsAfterAdd(t *testing.T) {
	// Create an OR-Set, add an element, assert Contains returns true
	A := NewORSet("A")
	A.Add("elem1")
	if !A.Contains("elem1") {
		t.Errorf("OR-Set does not contain expectd value: %s", "elem1")
	}
}

// Should not contain element after Remove
func TestORSetNotContainsAfterRemove(t *testing.T) {
	// Add an element, then remove it, assert Contains returns false
	A := NewORSet("A")
	elem1 := "elem1"
	A.Add(elem1)
	if !A.Contains(elem1) {
		t.Errorf("OR-Set does not contain expectd value: %s", elem1)
	}

	A.Remove(elem1)
	if A.Contains(elem1) {
		t.Errorf("OR-Set contains value: %s but should not", elem1)
	}

}

// Elements should return a sorted list of all live elements
func TestORSetElementsSorted(t *testing.T) {
	// Add several elements in non-alphabetical order
	// Assert Elements() returns them in sorted order
	A := NewORSet("A")

	if A.Elements() != nil {
		t.Errorf("elements should be nil")
	}

	A.Add("a")
	A.Add("c")
	A.Add("b")

	if A.Elements() == nil {
		t.Errorf("elements should not be nil")
	}

	if !slices.IsSorted(A.Elements()) {
		t.Errorf("elements does not return a sorted list")
	}
}

// Should contain elements from both replicas after merge
func TestORSetMergeCombinesElements(t *testing.T) {
	// Create two OR-Sets with different elements, merge one into the other
	// Assert both elements are present in the merged set
	A := NewORSet("A")
	B := NewORSet("B")

	A.Add("a")
	B.Add("b")

	A.Merge(B)
	if !A.Contains("a") || !A.Contains("b") {
		t.Errorf("A did not recieve B's merged elements")
	}
	B.Merge(A)
	if !B.Contains("a") || !B.Contains("b") {
		t.Errorf("B did not recieve A's merged elements")
	}
}

// Should be idempotent (merging same state twice has no effect)
func TestORSetMergeIdempotent(t *testing.T) {
	// Merge B into A twice — assert Elements() is identical after both merges
	A := NewORSet("A")
	B := NewORSet("B")

	A.Add("a")
	B.Add("b")

	A.Merge(B)
	elem1 := A.Elements()
	A.Merge(B)
	elem2 := A.Elements()

	if !slices.Equal(elem1, elem2) {
		t.Errorf("Merge must be idempotent. first merge result %v must equal second merge result %v", elem1, elem2)
	}
}

// Should be commutative (A absorbs B and B absorbs A reach the same result)
func TestORSetMergeCommutative(t *testing.T) {
	// Merge in both directions — assert both results contain the same elements
	A := NewORSet("A")
	B := NewORSet("B")

	A.Add("a")
	B.Add("b")

	A.Merge(B)
	elemA := A.Elements()
	B.Merge(A)
	elemB := B.Elements()

	if !slices.Equal(elemA, elemB) {
		t.Errorf("Merge must work in both directions. A Merge B result %v must equal B Merge A result %v", elemA, elemB)
	}

}

// Concurrent Add should survive a concurrent Remove (add-wins)
func TestORSetConcurrentAddRemoveAddWins(t *testing.T) {
	// Replica A adds an element (gets tag A:1)
	// Replica B also adds the same element (gets tag B:1) — before any merge
	// Replica A removes the element — this only tombstones A:1
	// Merge A and B — assert the element is still in the set (B:1 is still live)

	A := NewORSet("A")
	B := NewORSet("B")

	A.Add("a")
	B.Add("a")
	A.Remove("a")

	A.Merge(B)

	if !slices.Contains(A.Elements(), "a") {
		t.Errorf("Removed value: %s should still exist in A.Elements but got: %v", "a", A.Elements())
	}
}

// Re-adding a removed element should make it visible again
func TestORSetReAddAfterRemove(t *testing.T) {
	// Add, remove, then add again
	// Assert the element is in the set (the new Add creates a fresh tag)
	A := NewORSet("A")

	A.Add("a")
	A.Remove("a")
	if slices.Contains(A.Elements(), "a") {
		t.Errorf("Removed value: %s should NOT exist in A.Elements but got: %v", "a", A.Elements())
	}
	A.Add("a")

	if !slices.Contains(A.Elements(), "a") {
		t.Errorf("Removed value: %s was added back into the list and should still exist in A.Elements but got: %v", "a", A.Elements())
	}
}
