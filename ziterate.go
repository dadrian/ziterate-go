// Package ziterate provides ZMap-style cyclic group iterators for visiting a
// finite target space in a randomized order without tracking visited elements.
package ziterate

// Iterator is the common interface for cyclic group iterators.
//
// Next returns the next element, or nil after the iterator has completed one
// full cycle.
type Iterator interface {
	Next() interface{}
}
