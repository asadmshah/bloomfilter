package bloomfilter

import (
	"fmt"
)

// BitSet is an array of 8-bit integers.
type BitSet []uint8

// NewBitSet creates a new BitSet with size n.
func NewBitSet(n uint) BitSet {
	return BitSet(make([]uint8, n))
}

// Set sets the bit at index x to 1.
func (self BitSet) Set(x uint) {
	self[x / 8] |= (1 << uint8(x % 8))
}

// Get returns true if the bit at index x is set to 1.
func (self BitSet) Get(x uint) bool {
	var m uint8 = 1 << uint8(x % 8)
	return self[x / 8] & m == m
}

// ToBytes is a helper function to convert any type to an array of bytes.
// This can be used to convert numbers, but using encoding/binary is much more
// efficient.
func ToBytes(x interface{}) []byte {
	return []byte(fmt.Sprintf("%v", x))
}

