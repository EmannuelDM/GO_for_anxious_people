package bitset

import (
	"bytes"
	"fmt"
)

// IntSet is a set of small non-negative integers represented as a bit vector.
//
// Key Encapsulation Lesson:
// The client experiences this type as a logical mathematical set of integers.
// The internal implementation (a slice of 64-bit word bitmasks) is completely
// hidden (unexported 'words').
//
// The package author is free to change this implementation tomorrow (e.g. to a
// hash map, a balanced tree, or compressed Roaring Bitmaps) without breaking
// any client code!
type IntSet struct {
	words []uint64 // Unexported bit vector slice
}

// New creates an empty IntSet.
func New() *IntSet {
	return &IntSet{
		words: make([]uint64, 0, 4),
	}
}

// Add inserts a non-negative integer x into the set.
func (s *IntSet) Add(x int) {
	if x < 0 {
		return // Invariant: set only contains non-negative numbers
	}
	word, bit := x/64, uint(x%64)
	for word >= len(s.words) {
		s.words = append(s.words, 0)
	}
	s.words[word] |= 1 << bit
}

// Has reports whether the non-negative integer x is in the set.
func (s *IntSet) Has(x int) bool {
	if x < 0 {
		return false
	}
	word, bit := x/64, uint(x%64)
	return word < len(s.words) && (s.words[word]&(1<<bit)) != 0
}

// Remove deletes x from the set if present.
func (s *IntSet) Remove(x int) {
	if x < 0 {
		return
	}
	word, bit := x/64, uint(x%64)
	if word < len(s.words) {
		s.words[word] &^= 1 << bit
	}
}

// Len returns the total count of elements in the set.
func (s *IntSet) Len() int {
	var count int
	for _, word := range s.words {
		// Count set bits using Brian Kernighan's bit-counting algorithm
		w := word
		for w != 0 {
			w &= w - 1
			count++
		}
	}
	return count
}

// String returns the set representation as "{1 3 5}".
func (s *IntSet) String() string {
	var buf bytes.Buffer
	buf.WriteByte('{')
	var first = true
	for i, word := range s.words {
		if word == 0 {
			continue
		}
		for bit := 0; bit < 64; bit++ {
			if (word & (1 << bit)) != 0 {
				if !first {
					buf.WriteByte(' ')
				}
				first = false
				fmt.Fprintf(&buf, "%d", 64*i+bit)
			}
		}
	}
	buf.WriteByte('}')
	return buf.String()
}
