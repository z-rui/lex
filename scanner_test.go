package lex

import (
	"reflect"
	"strings"
	"testing"
)

func TestScanner(t *testing.T) {
	const str = `the quick brown fox
jumps over the lazy dog`
	s := &Scanner{
		RuneReader: strings.NewReader(str),
	}
	words := []string{}
	locs := [][2]int{}
L:
	for {
		c := s.Input()
		for !('a' <= c && c <= 'z') {
			if c == -1 {
				break L
			}
			c = s.Input()
		}
		s.Back(1)
		s.Flush()
		if c1 := s.Buffer()[0]; c1 != c {
			t.Errorf("s.Buffer[0] = %q, expect %q", c1, c)
		}
		for 'a' <= c && c <= 'z' {
			c = s.Input()
		}
		if c != -1 {
			s.Back(1)
		}
		words = append(words, string(s.Token))
		lin, col := s.LineCol(s.Position)
		locs = append(locs, [2]int{lin, col})
		s.Flush()
	}
	expect_words := strings.Fields(str)
	expect_locs := [][2]int{
		{0, 0}, {0, 4}, {0, 10}, {0, 16}, {1, 0},
		{1, 6}, {1, 11}, {1, 15}, {1, 20},
	}

	if !reflect.DeepEqual(words, expect_words) {
		t.Errorf("words: expect %q, got %q", expect_words, words)
	}
	if !reflect.DeepEqual(locs, expect_locs) {
		t.Errorf("locs: expect %v, got %v", expect_locs, locs)
	}
}
