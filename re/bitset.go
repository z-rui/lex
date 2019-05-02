package re

import (
	"strconv"
)

// Bitset represents a set of integers.
// Its representation is compact, and can be converted
// into a string for using as a key in maps.
// Note: Bitset is a reference type; call Clone if you
// need a copy.
type Bitset []byte

// NewBitset creates a bitset whose members are within
// the range [0,n).
func NewBitset(n int) Bitset {
	return make(Bitset, (n+7)/8)
}

// Set adds i into the set.
func (s Bitset) Set(i int) {
	s[i/8] |= 1 << (uint(i) % 8)
}

// Get tests if i is in the set.
func (s Bitset) Get(i int) bool {
	return (s[i/8] & (1 << (uint(i) % 8))) != 0
}

// UnionWith adds all elements in t into s.
func (s Bitset) UnionWith(t Bitset) {
	n := len(s)
	if len(t) != n {
		panic("bitset width mismatch")
	}
	for i := 0; i < n; i++ {
		s[i] |= t[i]
	}
	return
}

func (s Bitset) InterWith(t Bitset) {
	n := len(s)
	if len(t) != n {
		panic("bitset width mismatch")
	}
	for i := 0; i < n; i++ {
		s[i] &= t[i]
	}
	return
}

func (s Bitset) DiffWith(t Bitset) {
	n := len(s)
	if len(t) != n {
		panic("bitset width mismatch")
	}
	for i := 0; i < n; i++ {
		s[i] &^= t[i]
	}
	return
}

func (s Bitset) Subset(t Bitset) bool {
	n := len(s)
	if len(t) != n {
		panic("bitset width mismatch")
	}
	for i := 0; i < n; i++ {
		if s[i]|t[i] != t[i] {
			return false
		}
	}
	return true
}

func (s Bitset) Equal(t Bitset) bool {
	n := len(s)
	if len(t) != n {
		panic("bitset width mismatch")
	}
	for i := 0; i < n; i++ {
		if s[i] != t[i] {
			return false
		}
	}
	return true
}

// Clone creates a new bitset identical to s.
func (s Bitset) Clone() Bitset {
	n := len(s)
	s1 := make(Bitset, n)
	copy(s1, s)
	return s1
}

var log2 = []uint8{1: 0, 2: 1, 4: 2, 8: 3, 16: 4, 32: 5, 64: 6, 128: 7}

// ForEach calls the callback function for each elements in the bitset.
func (s Bitset) ForEach(callback func(int)) {
	for i, b := range s {
		for b > 0 {
			c := b & -b
			b ^= c
			callback(i*8 + int(log2[c]))
		}
	}
}

func (s Bitset) String() string {
	buf := []byte{'{'}
	s.ForEach(func(i int) {
		if len(buf) > 1 {
			buf = append(buf, ',')
		}
		buf = strconv.AppendInt(buf, int64(i), 10)
	})
	buf = append(buf, '}')
	return string(buf)
}

func (s Bitset) PopCount() (n int) {
	for _, b := range s {
		for b > 0 {
			b &= b - 1
			n++
		}
	}
	return
}

// StringKey converts the bitset to a string which can be used
// as a key in a map.
func (s Bitset) StringKey() string {
	return string(s)
}
