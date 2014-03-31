package bloomfilter

import (
	"testing"
	"fmt"
)

const (
	errExpected = "expected %v, got %v"
	fpRateGreater = "false pos rate (%f) is greater than expected (%f)."
)

func TestCalculations(t *testing.T) {
	filter := NewBloomFilter(50000, 0.01)
	if e, r := uint(479253), filter.M(); e != r {
		t.Errorf(errExpected, e, r)
	}
	if e, r := uint(7), filter.K(); e != r {
		t.Errorf(errExpected, e, r)
	}
}

func TestFalsePositiveRate(t *testing.T) {
	c := 0
	n := 100000
	filter := NewBloomFilter(n, 0.01)
	var entry []byte
	for i := 0; i < n; i++ {
		entry = []byte(fmt.Sprintf("%d", i))
		filter.Add(entry)
		entry = []byte(fmt.Sprintf("%d", i + 1))
		if filter.Has(entry) {
			c++
		}
	}
	p := float64(c) / float64(n)
	if p > filter.P() {
		t.Errorf(fpRateGreater, p, filter.P())
	}
	t.Logf("expected false positive rate: %f", filter.P())
	t.Logf("actual false positive rate: %f", p)
}

func BenchmarkBasic(b *testing.B) {
	var filter *BloomFilter
	var entry []byte
	for i := 0; i < b.N; i++ {
		filter = NewBloomFilter(50000, 0.01)
		for j := 0; j < 50000; j++ {
			b.StopTimer()
			entry = []byte(fmt.Sprintf("%d", i))
			b.StartTimer()
			filter.Add(entry)
			if filter.Has(entry) == false {
				b.Fatal("false negative encountered.")
			}
		}
	}
}

