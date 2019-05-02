package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/z-rui/lex"
	"github.com/z-rui/lex/re"
)

type Rule struct {
	Start  re.RangeSet // start condition
	Frag   *re.Frag    // NFA for matching this rule
	Action string      // semantic action
}

type startCond struct {
	Name string
	Id   int
	Excl bool
}

type LexGen struct {
	Prefix string
	Rules  []Rule
	Out    io.Writer
	lex    lexer
	start  map[string]startCond
	defs   map[string]*re.Frag
	dfa    *re.DFA
	acc    map[int]int // DFA state id => accepting rule id
}

func (gen *LexGen) scanRE(noSpace bool) *re.Frag {
	return re.Parse(gen.lex.Scanner, gen.defs, noSpace, true)
}

func (gen *LexGen) CopyUntilMark() {
	const tmpl0 = `
func ($$lex *$$Lex) Lex($$lval *$$SymType) int {
	var c rune

	BEGIN := func(i int32) { $$lex.Start = i }; _ = BEGIN
`
	for {
		gen.lex.Flush()
		gen.lex.scanLine()
		line := string(gen.lex.Token)
		if line == "%%\n" || line == "%%" || line == "" {
			break
		}
		io.WriteString(gen.Out, line)
	}
	io.WriteString(gen.Out, strings.ReplaceAll(tmpl0, "$$", gen.Prefix))
}

func (gen *LexGen) ReadDefs() {
	gen.start = map[string]startCond{
		"INITIAL": startCond{
			Name: "INITIAL",
			Id:   0,
			Excl: false,
		},
	}
	gen.defs = map[string]*re.Frag{}

L:
	for {
		gen.lex.Flush()
		c := gen.lex.Input()
		switch {
		case c == '%':
			c = gen.lex.Input()
			switch c {
			case 's', 'x':
				names := gen.lex.scanIdents()
				for _, name := range names {
					sc, ok := gen.start[name]
					sc.Name = name
					if !ok {
						sc.Id = len(gen.start)
					}
					sc.Excl = c == 'x'
					gen.start[name] = sc
				}
			case '%':
				gen.lex.scanLine()
				break L
			default:
				gen.lex.Error(fmt.Sprintf("unknown directive %c", c))
				gen.lex.scanLine()
			}
		case 'A' <= c && c <= 'Z' || c == '_' || 'a' <= c && c <= 'z':
			gen.lex.scanIdent()
			name := string(gen.lex.Token)
			gen.lex.skipWS(true)
			frag := gen.scanRE(false)
			gen.defs[name] = frag
		case c == '\n':
		case c == -1:
			gen.lex.Error("unexpected EOF")
			break L
		default:
			gen.lex.Error(fmt.Sprintf("bad character %q", c))
		}
	}
}

func (gen *LexGen) ReadRules() {
L:
	for {
		var r Rule

		gen.lex.Flush()
		c := gen.lex.Input()
		switch c {
		case -1:
			break L
		case ' ', '\t':
			gen.lex.scanLine()
			io.WriteString(gen.Out, string(gen.lex.Token))
			continue
		case '\n', '\v', '\f', '\r':
			continue
		case '<':
			names := gen.lex.scanStart()
			if len(names) == 1 && names[0] == "*" {
				for _, sc := range gen.start {
					r.Start = append(r.Start, re.Range{rune(sc.Id), rune(sc.Id)})
				}
			} else {
				for _, name := range names {
					sc, ok := gen.start[name]
					if !ok {
						gen.lex.Error(fmt.Sprintf("undefined %q", name))
						continue
					}
					r.Start = append(r.Start, re.Range{rune(sc.Id), rune(sc.Id)})
				}
			}
		default:
			gen.lex.Back(1)
		}
		r.Frag = gen.scanRE(true)
		gen.lex.skipWS(true)
		switch c := gen.lex.Input(); c {
		case -1:
			break L
		case '{':
			gen.lex.scanCodeFrag()
		case '\n':
		default:
			gen.lex.scanLine()
		}
		if r.Frag != nil {
			r.Action = string(gen.lex.Token)
			gen.Rules = append(gen.Rules, r)
			// log.Println("New rule:", r)
		}
	}
}

func (gen *LexGen) Gen() {
	if len(gen.Rules) < 1 {
		log.Fatal("no rules defined")
	}
	gen.genDFA()
	gen.minDFA()
}

func (gen *LexGen) genDFA() {
	incl := re.RangeSet{} // list of inclusive start conditions
	for _, sc := range gen.start {
		if !sc.Excl {
			incl = append(incl, re.Range{rune(sc.Id), rune(sc.Id)})
		}
	}
	incl = incl.Canon(false)

	n := len(gen.Rules)
	l := make([]*re.Frag, n)
	nfaAcc := make([]*re.NFAState, n) // rule id => NFA accept
	for i := range gen.Rules {
		r := &gen.Rules[i]
		startCond := r.Start
		if r.Start == nil { // no start condition
			startCond = incl
		}
		// augment Frag with its start condition
		l[i] = &re.Frag{
			Start: &re.NFAState{
				In:   startCond,
				Succ: []*re.NFAState{r.Frag.Start},
			},
			Accept: r.Frag.Accept,
		}
		nfaAcc[i] = l[i].Accept
	}

	nfa := re.Alter(l).Canon()
	//nfa.WriteDotTo(os.Stderr)
	gen.dfa = re.NFA2DFA(nfa)
	gen.acc = make(map[int]int)

	for i, s := range gen.dfa.States {
		for ruleId, acc := range nfaAcc {
			if s.Clos.Get(acc.Id) {
				gen.acc[i] = ruleId
				break
			}
		}
	}
}

func (gen *LexGen) minDFA() {
	n := len(gen.Rules)
	P := make([]re.Bitset, n+1)
	for i := range P {
		P[i] = re.NewBitset(len(gen.dfa.States))
	}
	for i := range gen.dfa.States {
		if j, ok := gen.acc[i]; ok {
			P[j].Set(i)
		} else {
			P[n].Set(i)
		}
	}
	k := 0
	empty := re.NewBitset(len(gen.dfa.States))
	for i, s := range P {
		if s.Equal(empty) {
			log.Printf("Rule %d cannot be derived", i)
		} else {
			P[k] = P[i]
			k++
		}
	}
	P = P[:k]
	P = gen.dfa.Hopcroft(P)
	gen.dfa = gen.dfa.Coalesce(P)
	acc := make(map[int]int)

	for i, s := range gen.dfa.States {
		if rId, ok := gen.acc[s.Id]; ok {
			acc[i] = rId
		}
		s.Id = i
	}
	gen.acc = acc
	_ = P
}

func (gen *LexGen) Run(filename string) {
	f, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	gen.lex.Scanner = &lex.Scanner{
		RuneReader: bufio.NewReader(f),
		Filename:   filename,
	}
	fmt.Fprintf(gen.Out, "// Generated from %s.  DO NOT EDIT.\n\n", filename)
	gen.CopyUntilMark()
	gen.ReadDefs()
	gen.ReadRules()
	gen.Gen()
	gen.Dump()
}
