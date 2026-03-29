package registerssolution

import "time"

// LWWRegister is a Last-Write-Wins Register CRDT.
// It stores a single string value. On merge, the replica with the
// higher timestamp wins. Equal timestamps keep the current value.
type LWWRegister struct {
	replicaID string
	value     string
	timestamp int64 // Unix nanoseconds
}

// NewLWWRegister creates a new LWW-Register for the given replica.
func NewLWWRegister(replicaID string) *LWWRegister {
	return &LWWRegister{replicaID: replicaID}
}

// Set updates the register's value, stamped with the current time.
func (r *LWWRegister) Set(value string) {
	r.SetAt(value, time.Now().UnixNano())
}

// SetAt updates the register's value with an explicit timestamp.
// Useful in tests for deterministic behavior.
func (r *LWWRegister) SetAt(value string, ts int64) {
	r.value = value
	r.timestamp = ts
}

// Get returns the current value of the register.
func (r *LWWRegister) Get() string {
	return r.value
}

// Merge combines another LWW-Register's state into this one.
// The higher timestamp wins; equal timestamps keep the current value.
func (r *LWWRegister) Merge(other *LWWRegister) {
	if other.timestamp > r.timestamp {
		r.value = other.value
		r.timestamp = other.timestamp
	}
}
