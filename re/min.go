package re

import (
	_ "fmt"
	"sort"
)

func (dfa *DFA) computePred() []map[int]RangeSet {
	n := len(dfa.States)
	predMap := make([]map[int]RangeSet, n)
	for i := range predMap {
		predMap[i] = make(map[int]RangeSet)
	}
	for i, s := range dfa.States {
		for _, succ := range s.Succ {
			pr := predMap[succ.Dest.Id]
			pr[i] = append(pr[i], succ.Input)
		}
	}
	return predMap
}

func (dfa *DFA) Hopcroft(P []Bitset) []Bitset {
	n := len(dfa.States)
	pred := dfa.computePred()
	W := map[string]Bitset{}
	for _, s := range P {
		W[s.StringKey()] = s
	}
	empty := NewBitset(n)

	X := NewBitset(n)
	inter := NewBitset(n)
	diff := NewBitset(n)
	//fmt.Printf("P = %v\n", P)
	for len(W) > 0 {
		var A Bitset
		//fmt.Println("W =", W)
		for k, v := range W {
			delete(W, k)
			A = v
			break
		}
		//fmt.Println("A =", A)

		predMap := map[int]RangeSet{}
		A.ForEach(func(i int) {
			for k, v := range pred[i] {
				predMap[k] = append(predMap[k], v...)
			}
		})
		//fmt.Println("A: predMap =", predMap)
		rss := []RangeSet(nil)
		prs := []int(nil)
		for k, v := range predMap {
			rss = append(rss, v.Canon(false))
			prs = append(prs, k)
		}
		//fmt.Println("A: rss =", rss, "prs =", prs)
		_, mapping := flattenRangeSets(rss)
		//fmt.Println("A: flat =", []Range(flat), "mapping =", mapping)
		for _, m := range mapping {
			copy(X, empty)
			for _, j := range m {
				X.Set(prs[j])
			}
			//fmt.Println("X =", X)
			for i := 0; i < len(P); i++ {
				Y := P[i]
				copy(inter, Y)
				inter.InterWith(X)
				if inter.Equal(empty) {
					continue
				}
				copy(diff, Y)
				diff.DiffWith(X)
				if diff.Equal(empty) {
					continue
				}
				inter := inter.Clone()
				diff := diff.Clone()
				P[i] = inter
				P = append(P, diff)
				//fmt.Printf("P = %v, inter=%v, diff=%v, W=%v\n", P, inter, diff, W)
				if _, ok := W[Y.StringKey()]; ok {
					delete(W, Y.StringKey())
					W[inter.StringKey()] = inter
					W[diff.StringKey()] = diff
				} else if inter.PopCount() <= diff.PopCount() {
					W[inter.StringKey()] = inter
				} else {
					W[diff.StringKey()] = diff
				}
			}
		}
	}
	k := 0
	//println("equivalence classes:")
	for _, ec := range P {
		if !ec.Equal(empty) {
			P[k] = ec
			k++
			/*ec.ForEach(func (i int) {
				print(i, " ")
			})
			println()*/
		}
	}
	P = P[:k]
	return P
}

func (dfa *DFA) Coalesce(P []Bitset) *DFA {
	n := len(P)
	states := make([]*DFAState, n)
	ecMap := make([]int, len(dfa.States))
	for i := range ecMap {
		ecMap[i] = -1
	}
	for i, part := range P {
		id := -1
		clos := NewBitset(len(dfa.nfa.States))
		part.ForEach(func(j int) {
			if id == -1 {
				id = j
			}
			ecMap[j] = i
			clos.UnionWith(dfa.States[j].Clos)
		})
		states[i] = &DFAState{
			Id:   id,
			Clos: clos,
		}
	}
	for i, part := range P {
		rangeVis := map[Range]int{}
		part.ForEach(func(j int) {
			s := dfa.States[j]
			for _, succ := range s.Succ {
				newId := ecMap[succ.Dest.Id]
				oldId, ok := rangeVis[succ.Input]
				if ok && oldId != newId {
					println("state", succ.Dest.Id, "old ec =", oldId, "new ec =", newId)
					panic("cannot coalesce")
				} else if !ok {
					states[i].Succ = append(states[i].Succ, DFAEdge{
						Input: succ.Input,
						Dest:  states[newId],
					})
					rangeVis[succ.Input] = newId
				}
			}
		})
	}
	sort.Slice(states, func(i, j int) bool {
		return states[i].Id < states[j].Id
	})
	return &DFA{
		nfa:    dfa.nfa,
		States: states,
	}
}

func (dfa *DFA) Minimize() *DFA {
	n := len(dfa.States)
	F := NewBitset(n)
	NF := NewBitset(n)
	for i, s := range dfa.States {
		if s.Clos.Get(dfa.nfa.Accept.Id) {
			F.Set(i)
		} else {
			NF.Set(i)
		}
	}
	P := dfa.Hopcroft([]Bitset{F, NF})
	dfa = dfa.Coalesce(P)
	for i, s := range dfa.States {
		s.Id = i
	}
	return dfa
}
