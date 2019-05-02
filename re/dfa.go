package re

import (
	"fmt"
	"io"
	"strconv"
)

// DFAEdge represents an edge between two DFA states.
type DFAEdge struct {
	Input Range     // on what input
	Dest  *DFAState // goto which state
}

// DFAState represents a state in the DFA.
type DFAState struct {
	Id   int
	Clos Bitset    // closure of NFA states
	Succ []DFAEdge // outgoing edges
}

// String converts the state into string form, which is a set of NFA states.
func (s *DFAState) String() string {
	buf := []byte{'{'}
	s.Clos.ForEach(func(i int) {
		if len(buf) > 1 {
			buf = append(buf, ',')
		}
		buf = strconv.AppendInt(buf, int64(i), 10)
	})
	buf = append(buf, '}')
	return string(buf)
}

// DFA represents a deterministic finite automata.
type DFA struct {
	nfa    *NFA
	States []*DFAState
}

func (a *NFA) computeClosure() []Bitset {
	n := len(a.States)
	closure := make([]Bitset, n)
	for i := range closure {
		closure[i] = NewBitset(n)
		closure[i].Set(i)
	}
	for {
		changed := false
		for i, s := range a.States {
			clos := closure[i]
			for _, succ := range s.Succ {
				if succ.In == nil {
					clos1 := closure[succ.Id]
					if !clos1.Subset(clos) {
						changed = true
						clos.UnionWith(clos1)
					}
				}
			}
		}
		if !changed {
			break
		}
	}
	return closure
}

// NFA2DFA converts an NFA into a DFA.
// The NFA must be canonicalized.
func NFA2DFA(a *NFA) *DFA {
	closure := a.computeClosure()

	states := []*DFAState(nil)
	stateMap := map[string]*DFAState{}
	enterState := func(clos Bitset) *DFAState {
		key := clos.StringKey()
		state := stateMap[key]
		if state == nil {
			state = &DFAState{
				Clos: clos,
				Id:   len(states),
			}
			states = append(states, state)
			stateMap[key] = state
			//fmt.Printf("new state: %s -> %d\n", state.String(), state.Id)
		}
		return state
	}
	enterState(closure[a.Start.Id])
	for k := 0; k < len(states); k++ {
		state := states[k]
		// for each (nfastate, range, dest),
		// add (range, union of closure[dest]) to successor,
		// and we need to keep the ranges disjoint.
		s := []RangeSet(nil)
		dest := []Bitset(nil)
		state.Clos.ForEach(func(i int) {
			for _, succ := range a.States[i].Succ {
				if succ.In != nil {
					s = append(s, succ.In)
					dest = append(dest, closure[succ.Id])
				}
			}
		})
		rs, mapping := flattenRangeSets(s)
		for i, r := range rs {
			clos := NewBitset(len(a.States))
			for _, m := range mapping[i] {
				clos.UnionWith(dest[m])
			}
			next := enterState(clos)
			//fmt.Printf("%s --[%s]--> %s\n", state, r, next)
			state.Succ = append(state.Succ, DFAEdge{r, next})
		}
	}

	return &DFA{
		nfa:    a,
		States: states,
	}
}

func (a *DFA) Accept(i int) bool {
	return a.States[i].Clos.Get(a.nfa.Accept.Id)
}

func (a *DFA) WriteDotTo(w io.Writer, showClosure bool) {
	fmt.Fprintf(w, "digraph dfa {\n"+
		"\trankdir=LR;\n")
	for _, s := range a.States {
		peripheries := 1
		if a.Accept(s.Id) {
			peripheries = 2
		}
		if showClosure {
			fmt.Fprintf(w, "\ts%d[label=%q, peripheries=%d];\n",
				s.Id, s.String(), peripheries)
		} else {
			fmt.Fprintf(w, "\ts%d[peripheries=%d];\n",
				s.Id, peripheries)
		}
		edges := map[*DFAState]RangeSet{}
		for _, e := range s.Succ {
			rs := edges[e.Dest]
			rs = append(rs, e.Input)
			edges[e.Dest] = rs
		}
		for succ, rs := range edges {
			rs = rs.Canon(false)
			fmt.Fprintf(w, "\ts%d->s%d[label=%q];\n",
				s.Id, succ.Id, rs.String())
		}
	}
	fmt.Fprintln(w, "}")
}
