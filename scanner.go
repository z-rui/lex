package lex

import (
	"io"
	"sort"
)

type Scanner struct {
	io.RuneReader
	Filename string
	Position int   // rune offset of the start of token
	Token    []rune
	buffered int   // characters buffered
	linePos  []int // position of the n-th '\n'
	eof      bool  // saw eof?
}

// Input returns the next rune from the input stream.
func (s *Scanner) Input() (r rune) {
	if s.buffered > 0 {
		n := len(s.Token)
		s.Token = s.Token[:n+1]
		s.buffered--
		r = s.Token[n]
	} else if s.eof {
		return -1
	} else {
		var err error
		r, _, err = s.ReadRune()
		switch err {
		case nil:
			if r == '\n' {
				pos := s.Position + len(s.Token)
				s.linePos = append(s.linePos, pos)
			}
			s.Token = append(s.Token, r)
		case io.EOF:
			s.eof = true
			return -1
		default:
			panic(err)
		}
	}
	return
}

// Back puts back the last n runes into the input stream (actually into a buffer).
func (s *Scanner) Back(n int) {
	s.Token = s.Token[:len(s.Token)-n]
	s.buffered += n
}

// Buffered returns the runes in the buffer (those that were put back by Back).
func (s *Scanner) Buffer() []rune {
	n := len(s.Token)
	return s.Token[n:n+s.buffered]
}

// Flush clears the current token and moves the position to its end.
func (s *Scanner) Flush() {
	n := len(s.Token)
	s.Token = s.Token[n:]
	s.Position += n
}

// LineCol returns the line and column index (starting from 0) of a given rune offset.
func (s *Scanner) LineCol(pos int) (l, c int) {
	l = sort.SearchInts(s.linePos, pos)
	c = pos
	if 0 < l {
		c -= s.linePos[l-1] + 1
	}
	return
}
