package re

import (
	"fmt"
	"os"

	"github.com/z-rui/lex"
)

type state int

// yyLex states
const (
	_INITIAL = iota
	_CHARSET
	_LITERAL
)

type yyLex struct {
	*lex.Scanner
	NoSpace bool   // terminate at space
	NoNewline bool // terminate at newline
	Defs map[string]*Frag
	state state
}

// hand-crafted lexer
func (l *yyLex) Lex(yylval *yySymType) (tok int) {
	l.Flush()

	c := l.Input()
	if c == -1 {
		return 0
	}

	switch l.state {
	case _INITIAL:
		switch c {
		case '[':
			l.state = _CHARSET
			switch c = l.Input(); c {
			case '^':
				tok = LBRACK_CARET
			default:
				l.Back(1)
				fallthrough
			case -1:
				tok = LBRACK
			}
		case '"':
			l.state = _LITERAL
			tok = LPAREN
		case '{':
			for c != '}' && c != -1 {
				c = l.Input()
			}
			if c == -1 {
				return 2
			}
			name := string(l.Token[1:len(l.Token)-1])
			if frag, ok := l.Defs[name]; ok {
				yylval.frag = frag
			} else {
				l.Error(fmt.Sprintf("undefined %q", name))
				return 2
			}
			tok = NAME
		case '(':
			tok = LPAREN
		case ')':
			tok = RPAREN
		case '|':
			tok = ALTER
		case '*':
			tok = STAR
		case '+':
			tok = PLUS
		case '.':
			tok = DOT
		case ' ', '\t':
			if l.NoSpace {
				l.Back(1)
				return 0
			}
			goto common
		case '\r', '\n', '\f', '\v':
			if l.NoNewline {
				l.Back(1)
				return 0
			}
			goto common
		default:
			goto common
		}
	case _CHARSET:
		switch c {
		case ']':
			l.state = _INITIAL
			tok = RBRACK
		case '-':
			tok = DASH
		default:
			goto common
		}
	case _LITERAL:
		switch c {
		case '"':
			l.state = _INITIAL
			tok = RPAREN
		default:
			goto common
		}
	}
	return
common:
	if c == '\\' {
		tok = l.scanEscape(yylval)
	} else {
		yylval.r = c
		tok = CHAR
	}
	return
}

func (l *yyLex) scanEscape(yylval *yySymType) int {
	c := l.Input()
	switch c {
	case 'n':
		c = '\n'
	case 'r':
		c = '\r'
	case 'f':
		c = '\f'
	case 't':
		c = '\t'
	case 'v':
		c = '\v'
	case 'b':
		c = '\b'
	case 'a':
		c = '\a'
	default:
		break
	case -1:
		return 2 // unk
	}
	yylval.r = c
	return CHAR
}

func (l *yyLex) Error(msg string) {
	lin, col := l.LineCol(l.Position)
	fmt.Fprintf(os.Stderr, "%s:%d:%d: %s\n", l.Filename, lin+1, col+1, msg)
}


