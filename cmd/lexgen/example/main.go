package main

//go:generate lexgen tiger.l
//go:generate go fmt lex.yy.go

import (
	"bufio"
	"fmt"
	"os"

	"github.com/z-rui/lex"
)

type yySymType struct{}
type yyLex struct {
	*lex.Scanner
	Start int32
}

func (l *yyLex) PosPrintf(pos int, s string, v ...interface{}) {
	lin, col := l.LineCol(pos)
	fmt.Printf("%s:%d:%d: %s\n", l.Filename, lin+1, col+1, fmt.Sprintf(s, v...))
}

func (l *yyLex) Printf(s string, v ...interface{}) {
	l.PosPrintf(l.Position, s, v...)
}

func (l *yyLex) Error(s string) {
	l.Printf("%s", s)
}


func main() {
	var stdin = bufio.NewReader(os.Stdin)
	var stdout = bufio.NewWriter(os.Stdout)
	defer stdout.Flush()

	l := &yyLex{
		Scanner: &lex.Scanner{
			RuneReader: stdin,
		},
	}
L:
	for {
		tok := l.Lex(nil)
		switch tok {
		case -1:
			panic("lexer gets stuck")
		case 0:
			break L
		}
	}
	if l.Start != 0 {
		l.Error("Unexpected EOF")
	}
}
