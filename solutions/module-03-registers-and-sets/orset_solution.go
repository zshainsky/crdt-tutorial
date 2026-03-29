package registerssolution

import (
	"fmt"
	"sort"
)

// ORSet is an Observed-Remove Set CRDT.
// Each Add generates a unique tag. Remove only tombstones the tags
// observed at remove time, so a concurrent Add (with a new tag) survives.
type ORSet struct {
	replicaID  string
	counter    int
	adds       map[string]map[string]struct{} // element -> set of live tags
	tombstones map[string]struct{}            // tags that have been removed
}

// NewORSet creates a new OR-Set for the given replica.
func NewORSet(replicaID string) *ORSet {
	return &ORSet{
		replicaID:  replicaID,
		adds:       make(map[string]map[string]struct{}),
		tombstones: make(map[string]struct{}),
	}
}

// Add inserts an element with a fresh unique tag.
func (s *ORSet) Add(element string) {
	s.counter++
	tag := fmt.Sprintf("%s:%d", s.replicaID, s.counter)
	if s.adds[element] == nil {
		s.adds[element] = make(map[string]struct{})
	}
	s.adds[element][tag] = struct{}{}
}

// Remove tombstones all currently-known tags for this element.
// Tags added concurrently by other replicas (not yet merged) are unaffected.
func (s *ORSet) Remove(element string) {
	for tag := range s.adds[element] {
		s.tombstones[tag] = struct{}{}
	}
}

// Contains returns true if the element has at least one live (non-tombstoned) tag.
func (s *ORSet) Contains(element string) bool {
	for tag := range s.adds[element] {
		if _, removed := s.tombstones[tag]; !removed {
			return true
		}
	}
	return false
}

// Elements returns a sorted slice of all elements currently in the set.
func (s *ORSet) Elements() []string {
	var result []string
	for elem := range s.adds {
		if s.Contains(elem) {
			result = append(result, elem)
		}
	}
	sort.Strings(result)
	return result
}

// Merge combines another OR-Set's state into this one.
// Takes the union of both add-sets and both tombstone sets.
func (s *ORSet) Merge(other *ORSet) {
	for elem, tags := range other.adds {
		if s.adds[elem] == nil {
			s.adds[elem] = make(map[string]struct{})
		}
		for tag := range tags {
			s.adds[elem][tag] = struct{}{}
		}
	}
	for tag := range other.tombstones {
		s.tombstones[tag] = struct{}{}
	}
}
