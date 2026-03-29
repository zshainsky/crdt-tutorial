// Add package declaration (package registers)
package registers

import "time"

// Import the "time" package

// LWWRegister is a Last-Write-Wins Register CRDT.
// It stores a single string value. On merge, the replica with the
// higher timestamp wins. Equal timestamps keep the current value.
type LWWRegister struct {
	replicaID string
	// Add a field for the current value (string)
	value string
	// Add a field for the timestamp of that value (int64, Unix nanoseconds)
	timestamp int64
}

// NewLWWRegister creates a new LWW-Register for the given replica.
func NewLWWRegister(replicaID string) *LWWRegister {
	// Return an initialized LWWRegister
	return &LWWRegister{
		replicaID: replicaID,
	}
}

// Set updates the register's value, stamped with the current time.
func (r *LWWRegister) Set(value string) {
	// Delegate to SetAt, passing the current time in nanoseconds
	r.SetAt(value, time.Now().UnixNano())
}

// SetAt updates the register's value with an explicit timestamp.
// Use this in tests for deterministic behavior.
func (r *LWWRegister) SetAt(value string, ts int64) {
	// Store both the value and the timestamp
	r.value = value
	r.timestamp = ts
}

// Get returns the current value of the register.
func (r *LWWRegister) Get() string {
	// Return the stored value
	return r.value
}

// Merge combines another LWW-Register's state into this one.
// The replica with the higher timestamp wins.
func (r *LWWRegister) Merge(other *LWWRegister) {
	// If the other register has a strictly higher timestamp,
	// adopt its value and timestamp
	if other.timestamp > r.timestamp {
		r.value = other.value
		r.timestamp = other.timestamp
	}
}
