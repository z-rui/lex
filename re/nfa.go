package re

import (
	"fmt"
	"io"
)

// NFAState represents a state in the NFA.
// Because of how the NFA is constructed, the incoming edges
// of each state must have the same charset.
// Thus we store this charset together with the state.
type NFAState struct {
	Id   int
	In   RangeSet // nil if ϵ
	Succ []*NFAState
}

func (s *NFAState) Clone() *NFAState {
	return &NFAState{
		Id:   s.Id,
		In:   s.In,
		Succ: append([]*NFAState(nil), s.Succ...),
	}
}

// Frag represents a fragment of NFA.
// Each fragments has one start state and
// (for our purpose) only one accepting state.
// Fragments can be constructed by construction functions
// (e.g., Empty, Literal, Concat, Alter and Kleene).
type Frag struct {
	Start, Accept *NFAState
}

// Once the final NFA is constructed, the user should call
// Canon to canonicalize the NFA.
type NFA struct {
	Frag
	States []*NFAState
}

// Empty returns an NFA that accepts an empty string.
func Empty() *Frag {
	s := new(NFAState)
	return &Frag{s, s}
}

// Literal returns an NFA that accepts a charset.
func Literal(rs RangeSet, inverted bool) *Frag {
	s := &NFAState{In: rs.Canon(inverted)}
	return &Frag{s, s}
}

// Concat constructs e1 e2 ... en.
func Concat(l []*Frag) *Frag {
	n := len(l)
	if n < 1 {
		panic("number of frags < 1")
	}
	for i := 0; i < n-1; i++ {
		prev := l[i]
		next := l[i+1]
		prev.Accept.Succ = append(prev.Accept.Succ, next.Start)
	}
	return &Frag{
		Start:  l[0].Start,
		Accept: l[n-1].Accept,
	}
}

// Alter constructs e1|e2|...|en.
func Alter(l []*Frag) *Frag {
	n := len(l)
	if n < 1 {
		panic("number of NFAs < 1")
	}
	if n == 1 {
		return l[0]
	}
	in := new(NFAState)
	out := new(NFAState)
	for _, s := range l {
		in.Succ = append(in.Succ, s.Start)
		s.Accept.Succ = append(s.Accept.Succ, out)
	}
	return &Frag{
		Start:  in,
		Accept: out,
	}
}

// Kleene constructs e* (Kleene star).
func Kleene(s *Frag) *Frag {
	out := new(NFAState)
	out.Succ = []*NFAState{s.Start}
	s.Accept.Succ = append(s.Accept.Succ, out)
	return &Frag{
		Start:  out,
		Accept: out,
	}
}

// KleenePlus constructs e+ (= ee*, but with fewer states).
func KleenePlus(s *Frag) *Frag {
	s1 := Kleene(s)
	s1.Accept = s.Accept
	return s1
}

func (a *Frag) numberStates() (states []*NFAState) {
	nextId := 0
	// DFS
	vis := map[*NFAState]bool{}
	stack := []*NFAState{a.Start}
	for len(stack) > 0 {
		s := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		s.Id = nextId
		nextId++
		states = append(states, s)
		vis[s] = true
		for _, succ := range s.Succ {
			if !vis[succ] {
				stack = append(stack, succ)
			}
		}
	}
	return
}

func (a *Frag) Clone() *Frag {
	// Graph copy
	b := &Frag{
		Start: a.Start.Clone(),
	}
	ptr := map[*NFAState]*NFAState{
		a.Start: b.Start,
	}
	stack := []*NFAState{b.Start}
	for len(stack) > 0 {
		s := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		for i, succ := range s.Succ {
			t, ok := ptr[succ]
			if !ok {
				t = succ.Clone()
				ptr[succ] = t
				stack = append(stack, t)
			}
			s.Succ[i] = t
		}
	}
	b.Accept = ptr[a.Accept]
	if b.Accept == nil {
		//a.Canon().WriteDotTo(os.Stderr)
		panic("accept == nil")
	}
	return b
}

// Canon converts a Frag into an NFA.
// The states will be numbered, and a new state is prepended
// if start transition is not epsilon.
func (a *Frag) Canon() *NFA {
	if a.Start.In != nil {
		start := new(NFAState)
		start.Succ = []*NFAState{a.Start}
		a.Start = start
	}
	res := &NFA{
		Frag:   *a,
		States: a.numberStates(),
	}
	return res
}

// WriteDotTo writes the NFA in the dot (graphviz) language.
func (a *NFA) WriteDotTo(w io.Writer) {
	fmt.Fprintf(w, "digraph nfa {\n"+
		"\trankdir=LR;\n"+
		"\tnode[shape=circle];\n"+
		"\ts%d[shape=doublecircle];\n",
		a.Accept.Id)
	for _, s := range a.States {
		for _, succ := range s.Succ {
			label := "ϵ"
			if succ.In != nil {
				label = succ.In.String()
			}
			fmt.Fprintf(w, "\ts%d->s%d[label=%q];\n",
				s.Id, succ.Id, label)
		}
	}
	fmt.Fprintln(w, "}")
}
