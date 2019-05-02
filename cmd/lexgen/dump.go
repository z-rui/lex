package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/z-rui/lex/re"
)

func (gen *LexGen) Dump() {
	const (
		tmpl1 = `
$$start:
	$$lex.Flush()
	$$leng := 0
	$$acc := -1

	goto $$s0

`
		// DFA here
		tmpl2 = `$$finish:
	$$lex.Back(len($$lex.Token) - $$leng)
	$$text := $$lex.Token[:]
	switch $$acc {
	case -1:
		_ = yytext
		return 0
`
		// actions here
		tmpl3 = `}
	goto $$start
}
`
	)

	gen.dumpStart()
	io.WriteString(gen.Out, strings.ReplaceAll(tmpl1, "$$", gen.Prefix))
	for _, s := range gen.dfa.States {
		gen.dumpState(s)
	}
	io.WriteString(gen.Out, strings.ReplaceAll(tmpl2, "$$", gen.Prefix))
	for i := range gen.Rules {
		gen.dumpRule(i)
	}
	io.WriteString(gen.Out, strings.ReplaceAll(tmpl3, "$$", gen.Prefix))
}

func (gen *LexGen) dumpStart() {
	starts := make([]string, len(gen.start))
	for _, sc := range gen.start {
		starts[sc.Id] = sc.Name
	}
	io.WriteString(gen.Out, "const (\n")
	fmt.Fprintf(gen.Out, "%s = iota\n", starts[0])
	for _, sc := range starts[1:] {
		fmt.Fprintf(gen.Out, "%s\n", sc)
	}
	io.WriteString(gen.Out, ")\n")
}

func (gen *LexGen) dumpState(s *re.DFAState) {
	fmt.Fprintf(gen.Out, "%ss%d:\n", gen.Prefix, s.Id)
	//fmt.Fprintf(gen.Out, "println(\"In state%d\")\n", s.Id)
	if ruleId, ok := gen.acc[s.Id]; ok {
		//fmt.Fprintf(gen.Out, "println(\"can accept rule %d\")\n", ruleId)
		fmt.Fprintf(gen.Out, "%sacc = %d\n", gen.Prefix, ruleId)
		fmt.Fprintf(gen.Out, "%sleng = len(%slex.Token)\n", gen.Prefix, gen.Prefix)
	}
	if len(s.Succ) > 0 {
		if s.Id == 0 {
			fmt.Fprintf(gen.Out, "c = rune(%slex.Start)\n", gen.Prefix)
		} else {
			fmt.Fprintf(gen.Out, "c = %slex.Input()\n", gen.Prefix)
		}
		fmt.Fprintf(gen.Out, "switch {\n")
		for _, e := range s.Succ {
			if r := e.Input; r.First == r.Last {
				fmt.Fprintf(gen.Out, "case c == %q:\n", r.First)
			} else {
				fmt.Fprintf(gen.Out, "case %q <= c && c <= %q:\n", r.First, r.Last)
			}
			fmt.Fprintf(gen.Out, "goto %ss%d\n", gen.Prefix, e.Dest.Id)
		}
		fmt.Fprintf(gen.Out, "default: goto %sfinish\n", gen.Prefix)
		fmt.Fprintf(gen.Out, "}\n")
	} else {
		fmt.Fprintf(gen.Out, "goto %sfinish\n", gen.Prefix)
	}
}

func (gen *LexGen) dumpRule(id int) {
	fmt.Fprintf(gen.Out, "case %d: %s\n", id, gen.Rules[id].Action)
}

func (gen *LexGen) DumpDFA() {
	out, err := os.Create("lex.yy.dot")
	if err != nil {
		log.Fatal(err)
	}
	defer out.Close()
	w := bufio.NewWriter(out)
	defer w.Flush()
	gen.dfa.WriteDotTo(w, false)
}
