# Bloom Filter
A simple data structure that can be used to test set membership. It's main
advantage is that it is space efficient compared to structures like hash tables
or basic arrays. The downside however is that it can result in false positives,
though false negatives are not possible. It's time complexity is dependent on
the number of hash functions used.

The basic structure is an array of _m_ bits with _k_ different hash functions.
When an element is added or checked, it is run through the hash functions. The
results of the hash identify positions in the bit array. The structure can then
set the bit to add an element or check the bit for existence of the element. 


### Usage
Call `go get github.com/asadmshah/bloomfilter` to install the package. Using it
is just as straightforward: 
~~~
package main

import (
	"fmt"
	"github.com/asadmshah/bloomfilter"
)

func main() {
	n := 500000
	p := 0.01
	e := []byte("Hello World")
	filter := NewBloomFilter(n, p)
	filter.Add(e)
	fmt.Println(filter.Has(e))
}
~~~
The `NewBloomFilter` constructor requires two variables: `n` the expected number
of elements and `p` the acceptable false positive rate. `Add` and `Has`, both
require `[]byte` type as their argument. This package does provide a `ToBytes`
helper function for converting anything to bytes but it's probably very
inefficient compared to using something like `encoding/binary` when converting
numbers. 

By default, the constructor will evaluate the optimum number of bits and hash
functions based on the arguments provided. If you'd like more control, the
`BloomFilter` type provides methods `SetK` and `SetM` for manually setting them.
Calling these methods will automatically clear out any existing data.


### Hashing
You're free to change the default Hashing algorithm to any that satisfies the
`hash.Hash64` interface. The FNV1a hash provided by the standard library is the
default one in use. Based on [this
paper](http://www.eecs.harvard.edu/~kirsch/pubs/bbbf/rsa.pdf), only two hash
functions are required to create a Bloom filter as long as _k_ hashes are
generated using `gi(x) = h1(x) + (h2(x) * i)` with `i` running from `0:k-1`.
We can use a single 64-bit hash function to mock two different hash functions by
splitting the bytes in half, hence the need for `filter.H` to satisfy the
interface for `hash.Hash64` instead of `hash.Hash32`, [from
here](http://willwhim.wpengine.com/2011/09/03/producing-n-hash-functions-by-hashing-only-once/).

