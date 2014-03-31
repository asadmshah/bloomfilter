// Package bloomfilter provides a space efficient data structure called a bloom
// filter. It's main advantage is that it is space efficient compared to
// structures like hash tables or basic arrays. The downside however is that it
// can result in false positives,/ though false negatives are not possible. It's
// time complexity is dependent on/ the number of hash functions used.
package bloomfilter

import (
	"encoding/binary"
	"hash"
	"hash/fnv"
	"math"
)

// BloomFilter represents a basic bloom filter.
type BloomFilter struct {
	m	float64
	k	float64
	p	float64
	n	float64
	b	BitSet
	H	hash.Hash64
}

// NewBloomFilter returns an empty bloom filter.
// It chooses the optimum number of hash functions, k, and bits, m, to be used
// based on the given expected number of elements, n, and the acceptable false
// positive rate.
func NewBloomFilter(n int, p float64) *BloomFilter {
	filter := new(BloomFilter)
	filter.n = float64(n)
	filter.p = p
	filter.m = filter.CalculateM()
	filter.k = filter.CalculateK()
	filter.H = fnv.New64a()
	filter.Clear()
	return filter
}

// CalculateM returns the number of bits required in order to satisfy the false
// positive rate.
func (self *BloomFilter) CalculateM() float64 {
	return -self.n * math.Log(self.p) / math.Pow(math.Log(2.0), 2.0)
}

// CalculateK returns the number of hashes required in order to satisfy the
// false positive rate.
func (self *BloomFilter) CalculateK() float64 {
	return self.m / self.n * math.Log(2.0)
}

// CalculateP returns the expected false positive rate based on the number of
// hash functions and the size of the bitset.
func (self *BloomFilter) CalculateP() float64 {
	return math.Pow(1.0 - math.Pow(math.E, (-self.k * self.n / self.m)), self.k)
}

// Clear generates and sets a new empty bitset.
func (self *BloomFilter) Clear() {
	self.b = NewBitSet(uint(math.Ceil(self.m / 8.0)))
}

// hash generates the hash of the given byte array and returns the hash split
// into two. This turns a single 64-bit hash function into two different ones.
func (self *BloomFilter) hash(member []byte) (uint, uint) {
	self.H.Reset()
	self.H.Write(member)
	h := self.H.Sum(nil)
	a := binary.BigEndian.Uint32(h[:4])
	b := binary.BigEndian.Uint32(h[4:])
	return uint(a), uint(b)
}

// Add adds a member to the set.
func (self *BloomFilter) Add(member []byte) {
	a, b := self.hash(member)
	m := self.M()
	x := uint(0)
	for i := uint(0); i < self.K(); i++ {
		x = (a + b * i) % m
		self.b.Set(x)
	}
}

// Has returns true if the member exists in the set.
func (self *BloomFilter) Has(member []byte) bool {
	a, b := self.hash(member)
	m := self.M()
	x := uint(0)
	for i := uint(0); i < self.K(); i++ {
		x = (a + b * i) % m
		if !self.b.Get(x) {
			return false
		}
	}
	return true
}

// K returns the number of hash functions to be used in each operation.
func (self *BloomFilter) K() uint {
	return uint(math.Ceil(self.k))
}

// M returns the size of the set.
func (self *BloomFilter) M() uint {
	return uint(math.Ceil(self.m))
}

// P returns the false positive rate.
func (self *BloomFilter) P() float64 {
	return self.p
}

// SetK can be used to manually set the number of hash functions.
// This will empty out the set, so use this before inserting any items.
func (self *BloomFilter) SetK(k uint) {
	self.k = float64(k)
	self.Clear()
}

// SetM can be used to manually set the number of bits.
// This will empty out the set, so use this before inserting any items.
func (self *BloomFilter) SetM(m uint) {
	self.m = float64(m)
	self.Clear()
}
