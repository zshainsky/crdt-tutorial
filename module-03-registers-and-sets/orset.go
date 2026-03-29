// Add package declaration (package registers)
package registers

import (
	"fmt"
	"slices"
)

// ORSet is an Observed-Remove Set CRDT.
// It supports Add and Remove with add-wins semantics:
// a concurrent Add always survives a concurrent Remove.
//
// Each Add generates a unique tag. Remove only tombstones the tags
// the removing replica has already observed.
type ORSet struct {
	replicaID string
	counter   int
	// Add a map from element to its set of live tags
	//       (each element maps to a set of tag strings)
	adds map[string]map[string]struct{}
	// Add a set of tombstoned tags (removed tags)
	tombstones map[string]struct{}
}

// NewORSet creates a new OR-Set for the given replica.
func NewORSet(replicaID string) *ORSet {
	// Return an initialized ORSet with empty maps
	return &ORSet{
		replicaID:  replicaID,
		counter:    0,
		adds:       make(map[string]map[string]struct{}),
		tombstones: make(map[string]struct{}),
	}
}

// Add inserts an element with a fresh unique tag.
func (s *ORSet) Add(element string) {
	// Increment the counter
	s.counter++
	// Generate a unique tag by combining replicaID and counter
	tag := fmt.Sprintf("%s:%d", s.replicaID, s.counter)
	// Record this tag as a live tag for the element
	if s.adds[element] == nil {
		s.adds[element] = make(map[string]struct{})
	}
	s.adds[element][tag] = struct{}{}
}

// Remove tombstones all currently-known tags for this element.
// Tags added concurrently by other replicas (not yet merged) are unaffected.
func (s *ORSet) Remove(element string) {
	// Move all known tags for this element into the tombstone set
	for tag := range s.adds[element] {
		s.tombstones[tag] = struct{}{}
	}
}

// Contains returns true if the element has at least one live (non-tombstoned) tag.
func (s *ORSet) Contains(element string) bool {
	// Return true if any tag for this element is not in the tombstone set
	for tag := range s.adds[element] {
		if _, ok := s.tombstones[tag]; !ok {
			return true
		}
	}
	return false
}

// Elements returns a sorted slice of all elements currently in the set.
func (s *ORSet) Elements() []string {
	// Collect all elements where Contains returns true
	// Sort the result before returning
	var res []string
	for elem := range s.adds {
		if s.Contains(elem) {
			res = append(res, elem)
		}
	}
	slices.Sort(res)
	return res
}

// Merge combines another OR-Set's state into this one.
func (s *ORSet) Merge(other *ORSet) {
	// For each element in the other set, union its tags into this set's tags
	// Union the tombstone sets
	for elem, tags := range other.adds {
		for tag := range tags {
			if s.adds[elem] == nil {
				s.adds[elem] = make(map[string]struct{})
			}
			s.adds[elem][tag] = struct{}{}
		}
	}

	for tag := range other.tombstones {
		s.tombstones[tag] = struct{}{}
	}
}
