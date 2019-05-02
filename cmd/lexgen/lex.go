package main

import (
	"fmt"
	"os"

	"github.com/z-rui/lex"
)

type lexer struct {
	*lex.Scanner
	errs int
}

func (l *lexer) Error(err string) {
	lin, col := l.LineCol(l.Position)
	fmt.Fprintf(os.Stderr, "%s:%d:%d: %s\n", l.Filename, lin+1, col+1, err)
	l.errs++
	if l.errs >= 10 {
		fmt.Fprintf(os.Stderr, "too many errors\n")
		os.Exit(1)
	}
}

func (l *lexer) scanLine() {
	for {
		c := l.Input()
		if c == '\n' || c == -1 {
			break
		}
	}
}

func (l *lexer) skipWS(noNewline bool) {
L:
	for {
		c := l.Input()
		switch c {
		case ' ', '\t':
			/* fine */
		case '\v', '\n', '\r', '\f':
			if !noNewline {
				break
			}
			fallthrough
		default:
			l.Back(1)
			fallthrough
		case -1:
			break L
		}
	}
	l.Flush()
	return
}

func (l *lexer) scanIdent() {
	for {
		c := l.Input()
		switch {
		case '0' <= c && c <= '9', 'A' <= c && c <= 'Z', c == '_', 'a' <= c && c <= 'z':
		case c == -1:
			return
		default:
			l.Back(1)
			return
		}
	}
}

func (l *lexer) scanIdents() (names []string) {
	for {
		l.skipWS(true)
		c := l.Input()
		switch {
		case 'A' <= c && c <= 'Z', c == '_', 'a' <= c && c <= 'z':
		case c == '\n':
			return
		case c == -1:
			l.Error("unexpected EOF")
			return
		default:
			l.Error(fmt.Sprintf("bad character %q", c))
			continue
		}
		l.scanIdent()
		names = append(names, string(l.Token))
	}
}

func (l *lexer) scanStart() (names []string) {
	for {
		l.Flush()
		c := l.Input()
		switch {
		case 'A' <= c && c <= 'Z', c == '_', 'a' <= c && c <= 'z':
		case c == -1:
			l.Error("unexpected EOF")
			return
		case c == '*':
			if c = l.Input(); c == '>' {
				return []string{"*"}
			}
			fallthrough
		default:
			l.Error(fmt.Sprintf("bad character %q", c))
			continue
		}
		l.scanIdent()
		names = append(names, string(l.Token))
		c = l.Input()
		switch c {
		case ',', -1:
		case '>':
			return
		default:
			l.Back(1)
		}
	}
}

func (l *lexer) scanCodeFrag() {
	level := 1
L:
	for level > 0 {
		c := l.Input()
		switch c {
		case '{':
			level++
		case '}':
			level--
		case -1:
			l.Error("unexpected EOF, unmatched {")
			break L
		}
	}
}
