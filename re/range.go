package re

import (
	"container/heap"
	"sort"
	"strconv"
	"unicode"
)

// Range represents a range of runes.
type Range struct {
	First, Last rune
}

func appendRune(buf []byte, r rune) []byte {
	switch r {
	case '-', '^', ']': // need escape inside []
		buf = append(buf, '\\')
		buf = append(buf, byte(r))
		return buf
	default:
		buf2 := strconv.AppendQuoteRune(make([]byte, 0, 16), r)
		// strip the single quotes
		return append(buf, buf2[1:len(buf2)-1]...)
	}
}

func (r Range) appendTo(buf []byte) []byte {
	buf = appendRune(buf, r.First)
	switch n := r.Last - r.First; n {
	default:
		buf = append(buf, '-')
		fallthrough
	case 1:
		buf = appendRune(buf, r.Last)
	case 0:
	}
	return buf
}

// String converts a range [a,b] into string "a", "ab" or "a-b",
// depending on the length of the range.
// Non-printable runes are escaped.
func (r Range) String() string {
	return string(r.appendTo(nil))
}

// RangeSet holds a set of ranges.
type RangeSet []Range

func (s RangeSet) Len() int           { return len(s) }
func (s RangeSet) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s RangeSet) Less(i, j int) bool { return s[i].First < s[j].First }

// String converts ranges into a string in regex notation.
func (s RangeSet) String() string {
	buf := []byte{'['}
	n := len(s)
	if n > 0 && s[n-1].Last == unicode.MaxRune && s[n-1].First != 0 {
		// probably an inverted range
		buf = append(buf, '^')
		s = append(RangeSet{}, s...).Canon(true)
	}
	for _, r := range s {
		buf = r.appendTo(buf)
	}
	buf = append(buf, ']')
	return string(buf)
}

// Canon rearranges ranges so that they
// 1) are separate (disjoint and not adjacent), and
// 2) appear in increasing order.
// The original set may be altered.
func (s RangeSet) Canon(inverted bool) (ret RangeSet) {
	// linear time other than sorting
	sort.Sort(s)
	curr := Range{0, -1}
	ret = s[:0] // we can actually reuse the space
	for _, r := range s {
		if r.First > curr.Last+1 {
			if inverted {
				curr = Range{curr.Last + 1, r.First - 1}
			}
			if curr.First <= curr.Last {
				ret = append(ret, curr)
			}
			curr = r
		} else if r.Last > curr.Last {
			curr.Last = r.Last
		}
	}
	if inverted {
		curr = Range{curr.Last + 1, unicode.MaxRune}
	}
	if curr.First <= curr.Last {
		ret = append(ret, curr)
	}
	return
}

type rsHeapEnt struct {
	r RangeSet
	i int
}
type rsHeap struct {
	h []rsHeapEnt
}

func (h rsHeap) Less(i, j int) bool {
	return h.h[i].r[0].First < h.h[j].r[0].First
}
func (h rsHeap) Len() int            { return len(h.h) }
func (h rsHeap) Swap(i, j int)       { h.h[i], h.h[j] = h.h[j], h.h[i] }
func (h *rsHeap) Push(x interface{}) { h.h = append(h.h, x.(rsHeapEnt)) }
func (h *rsHeap) Pop() (r interface{}) {
	n := len(h.h)
	r = h.h[n-1]
	h.h = h.h[:n-1]
	return
}

func (h *rsHeap) pop() (r rHeapEnt) {
	top := heap.Pop(h).(rsHeapEnt)
	r = rHeapEnt{top.r[0], top.i}
	top.r = top.r[1:]
	if len(top.r) > 0 {
		heap.Push(h, top)
	}
	return
}

type rHeapEnt struct {
	r Range
	i int // set index
}
type rHeap struct { // ordered by Last
	h []rHeapEnt
}

func (h rHeap) Less(i, j int) bool {
	return h.h[i].r.Last < h.h[j].r.Last
}
func (h rHeap) Len() int            { return len(h.h) }
func (h rHeap) Swap(i, j int)       { h.h[i], h.h[j] = h.h[j], h.h[i] }
func (h *rHeap) Push(x interface{}) { h.h = append(h.h, x.(rHeapEnt)) }
func (h *rHeap) Pop() (r interface{}) {
	n := len(h.h)
	r = h.h[n-1]
	h.h = h.h[:n-1]
	return
}

// rearrange canonicalized range sets into disjoint ranges,
// each mapping to a set of range sets.
func flattenRangeSets(s []RangeSet) (ret RangeSet, mapping [][]int) {
	n := len(s)
	src := &rsHeap{h: make([]rsHeapEnt, n)}
	for i := range src.h {
		src.h[i] = rsHeapEnt{s[i], i}
	}
	heap.Init(src)

	curr := unicode.MaxRune
	h := &rHeap{}
	append1 := func(r Range) {
		ret = append(ret, r)
		m := []int{}
		for _, ent := range h.h {
			m = append(m, ent.i)
		}
		mapping = append(mapping, m)
	}
	popUntil := func(next rune) {
		for h.Len() > 0 {
			last := h.h[0].r.Last
			if last >= next {
				break
			}
			append1(Range{curr, last})
			curr = last + 1
			for h.Len() > 0 && h.h[0].r.Last == last {
				heap.Pop(h)
			}
		}
	}

	for src.Len() > 0 {
		top := src.pop()
		next := top.r.First
		popUntil(next)
		if h.Len() > 0 && curr < next {
			append1(Range{curr, next - 1})
		}
		curr = next
		heap.Push(h, top)
	}
	popUntil(unicode.MaxRune + 1)
	return
}
