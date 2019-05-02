package re

//go:generate lexgen re.l
//go:generate lrgen re.y
//go:generate go fmt lex.yy.go y.tab.go

import (
	"strings"

	"github.com/z-rui/lex"
)

type yyLex struct {
	*lex.Scanner
	Start      int32
	Whitespace uint
	Defs       map[string]*Frag
}

// Parse parses a regex into an Frag.
func Parse(sc *lex.Scanner, defs map[string]*Frag, noSpace, noNewline bool) *Frag {
	l := &yyLex{
		Scanner:   sc,
		Defs:      defs,
	}
	if !noSpace {
		l.Whitespace |= (1 << ' ') | (1 << '\t')
	}
	if !noNewline {
		l.Whitespace |= (1 << '\n') | (1 << '\r') | (1 << '\v') | (1 << '\f')
	}
	if res := yyParse(l); res != nil {
		return res.frag
	}
	return nil
}

// ParseString parses a regex in string into an Frag.
func ParseString(s string) *Frag {
	sc := &lex.Scanner{
		RuneReader: strings.NewReader(s),
		Filename:   "<string>",
	}
	return Parse(sc, nil, false, false)
}
