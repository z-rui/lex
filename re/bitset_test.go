package re

import (
	"testing"
)

func TestBitset(t *testing.T) {
	s := NewBitset(10)
	s.Set(0)
	s.Set(8)

	if s.Get(0) == false {
		t.Error("0 is not in set")
	}
	if s.Get(8) == false {
		t.Error("8 is not in set")
	}
	if s.Get(1) {
		t.Error("1 is in set")
	}
	if s.Get(9) {
		t.Error("1 is in set")
	}

	s1 := NewBitset(10)
	s1.Set(1)
	s1.Set(9)

	if s.Equal(s1) {
		t.Error("s = s1")
	}
	s.UnionWith(s1)
	if !s1.Subset(s) {
		t.Error("s not subset s1")
	}

	result := uint(0)
	s.ForEach(func(i int) {
		result |= 1 << uint(i)
	})
	// result = {0, 1, 8, 9} = 0x0303
	if result != 0x0303 {
		t.Errorf("result = %#x", result)
	}

	if key := s.StringKey(); key != "\x03\x03" {
		t.Errorf("key = %q", key)
	}

	s.DiffWith(s1)
	if n := s.PopCount(); n != 2 {
		t.Errorf("s.PopCount() = %d", n)
	}
	s.InterWith(s1)
	if n := s.PopCount(); n != 0 {
		t.Errorf("s.PopCount() = %d", n)
	}
}
