// TODO: Add package declaration (package registers)
package registers

// TODO: Import "fmt" and "sort"

// ORSet is an Observed-Remove Set CRDT.
// It supports Add and Remove with add-wins semantics:
// a concurrent Add always survives a concurrent Remove.
//
// Each Add generates a unique tag. Remove only tombstones the tags
// the removing replica has already observed.
type ORSet struct {
	replicaID string
	counter   int
	// TODO: Add a map from element to its set of live tags
	//       (each element maps to a set of tag strings)
	// TODO: Add a set of tombstoned tags (removed tags)
}

// NewORSet creates a new OR-Set for the given replica.
func NewORSet(replicaID string) *ORSet {
	// TODO: Return an initialized ORSet with empty maps
	return nil
}

// Add inserts an element with a fresh unique tag.
func (s *ORSet) Add(element string) {
	// TODO: Increment the counter
	// TODO: Generate a unique tag by combining replicaID and counter
	// TODO: Record this tag as a live tag for the element
}

// Remove tombstones all currently-known tags for this element.
// Tags added concurrently by other replicas (not yet merged) are unaffected.
func (s *ORSet) Remove(element string) {
	// TODO: Move all known tags for this element into the tombstone set
}

// Contains returns true if the element has at least one live (non-tombstoned) tag.
func (s *ORSet) Contains(element string) bool {
	// TODO: Return true if any tag for this element is not in the tombstone set
	return false
}

// Elements returns a sorted slice of all elements currently in the set.
func (s *ORSet) Elements() []string {
	// TODO: Collect all elements where Contains returns true
	// TODO: Sort the result before returning
	return nil
}

// Merge combines another OR-Set's state into this one.
func (s *ORSet) Merge(other *ORSet) {
	// TODO: For each element in the other set, union its tags into this set's tags
	// TODO: Union the tombstone sets
}
