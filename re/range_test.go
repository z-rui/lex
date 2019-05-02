package re

import (
	"reflect"
	"testing"
)

func TestRangeString(t *testing.T) {
	test := func(first, last rune, expect string) {
		if s := (Range{first, last}).String(); s != expect {
			t.Errorf("expect %q, got %q", expect, s)
		}
	}
	test('a', 'a', "a")
	test('\t', '\t', "\\t")
	test(']', ']', "\\]")
	test('a', 'b', "ab")
	test('a', 'z', "a-z")
}

func TestRangeSetCanon(t *testing.T) {
	s0 := RangeSet{
		{'_', '_'},
		// adjacent ranges
		{'a', 'm'},
		{'n', 'z'},
		// overlapping ranges
		{'A', 'Y'},
		{'B', 'Z'},
	}

	expect := RangeSet{
		// assume Unicode ordering
		{'A', 'Z'},
		{'_', '_'},
		{'a', 'z'},
	}

	s := s0.Canon(false)
	if !reflect.DeepEqual(s, expect) {
		t.Errorf("s = %v", s)
	}

	// inverting twice should give the same result
	s = s0.Canon(true).Canon(true)
	if !reflect.DeepEqual(s, expect) {
		t.Errorf("s = %v", s)
	}
}

func TestRangeSetFlatten(t *testing.T) {
	/* ABCDEFGHIJKLMNOPQR
	 * ---------------    1
	 *    ----  ---  - -- 2
	 *    --------  --- - 4
	 * 111777755773157426 */
	s0 := []RangeSet{
		{{'A', 'O'}},
		{{'D', 'G'}, {'J', 'L'}, {'O', 'O'}, {'Q', 'R'}},
		{{'D', 'K'}, {'N', 'P'}, {'R', 'R'}},
	}
	s1, mapping := flattenRangeSets(s0)

	expect := RangeSet{
		{'A', 'C'}, {'D', 'G'}, {'H', 'I'}, {'J', 'K'},
		{'L', 'L'}, {'M', 'M'}, {'N', 'N'}, {'O', 'O'},
		{'P', 'P'}, {'Q', 'Q'}, {'R', 'R'},
	}

	if !reflect.DeepEqual(s1, expect) {
		t.Errorf("s1 = %v, expect %v", s1, expect)
	}

	expect_cksum := []int{1, 7, 5, 7, 3, 1, 5, 7, 4, 2, 6}
	for i, m := range mapping {
		cksum := 0
		for _, idx := range m {
			cksum |= 1 << uint(idx)
		}
		if cksum != expect_cksum[i] {
			t.Errorf("cksum[%d] = %d, expect %d",
				i, cksum, expect_cksum[i])
		}
	}

	s0 = []RangeSet{{{'a', 'a'}, {'c', 'c'}}}
	s1, mapping = flattenRangeSets(s0)

	expect = RangeSet{{'a', 'a'}, {'c', 'c'}}
	if !reflect.DeepEqual(s1, expect) {
		t.Errorf("s1 = %v, expect %v", s1, expect)
	}
}
