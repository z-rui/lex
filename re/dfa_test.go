package re

import (
	"testing"
)

func TestNFA2DFA(t *testing.T) {
	acc := &NFAState{}
	s1 := &NFAState{
		In:   RangeSet{{'a', 'z'}},
		Succ: []*NFAState{acc},
	}
	s2 := &NFAState{
		In:   RangeSet{{'b', 'g'}},
		Succ: []*NFAState{acc},
	}
	start := &NFAState{
		Succ: []*NFAState{s1, s2},
	}

	nfa := (&Frag{
		Start:  start,
		Accept: acc,
	}).Canon()

	dfa := NFA2DFA(nfa)

	if len(dfa.States) == 0 {
		t.Error("DFA has no state")
	} else {
		start := dfa.States[0]
		n := len(start.Succ)
		if n != 3 { // [a], [b-g] and [h-z]
			t.Errorf("start.Succ = %d", n)
		}
	}
}
