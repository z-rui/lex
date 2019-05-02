// Generated from re.l.  DO NOT EDIT.

package re

import (
	"fmt"
	"os"
	"strconv"

	"github.com/z-rui/lex"
)

type yyLex struct {
	*lex.Scanner
	Start     int32
	NoSpace   bool
	NoNewline bool
	Defs      map[string]*Frag
}

func unquote(text []rune) (rune, []rune) {
	conv := func(i, j int, base int) (rune, []rune) {
		if j <= len(text) {
			n, _ := strconv.ParseInt(string(text[i:j]), base, 32)
			return rune(n), text[j:]
		}
		return text[0], text[1:]
	}
	r := text[0]
	text = text[1:]
	if r != '\\' || len(text) < 1 {
		return r, text
	}
	r = text[0]
	switch r {
	case 'a':
		r = '\a'
	case 'b':
		r = '\b'
	case 't':
		r = '\t'
	case 'n':
		r = '\n'
	case 'v':
		r = '\v'
	case 'f':
		r = '\f'
	case 'r':
		r = '\r'
	case '0', '1', '2', '3', '4', '5', '6', '7':
		return conv(0, 3, 8)
	case 'x':
		return conv(1, 3, 16)
	case 'u':
		return conv(1, 5, 16)
	case 'U':
		return conv(1, 9, 16)
	}
	return r, text
}

func (l *yyLex) Error(message string) {
	lin, col := l.LineCol(l.Position)
	fmt.Fprintf(os.Stderr, "%s:%d:%d: %s\n", l.Filename, lin+1, col+1, message)
}

func (yylex *yyLex) Lex(yylval *yySymType) int {
	var c rune

	BEGIN := func(i int32) { yylex.Start = i }
	_ = BEGIN
	const (
		INITIAL = iota
		LITERAL
		CHARSET
	)

yystart:
	yylex.Flush()
	yyleng := 0
	yyacc := -1

	goto yys0

yys0:
	c = rune(yylex.Start)
	switch {
	case c == '\x00':
		goto yys1
	case c == '\x01':
		goto yys2
	case c == '\x02':
		goto yys3
	default:
		goto yyfinish
	}
yys1:
	c = yylex.Input()
	switch {
	case '\x00' <= c && c <= '\b':
		goto yys4
	case c == '\t':
		goto yys5
	case c == '\n':
		goto yys6
	case '\v' <= c && c <= '\r':
		goto yys6
	case '\x0e' <= c && c <= '\x1f':
		goto yys4
	case c == ' ':
		goto yys5
	case c == '!':
		goto yys4
	case c == '"':
		goto yys7
	case '#' <= c && c <= '\'':
		goto yys4
	case c == '(':
		goto yys8
	case c == ')':
		goto yys9
	case c == '*':
		goto yys10
	case c == '+':
		goto yys11
	case ',' <= c && c <= '-':
		goto yys4
	case c == '.':
		goto yys12
	case '/' <= c && c <= 'Z':
		goto yys4
	case c == '[':
		goto yys13
	case c == '\\':
		goto yys14
	case ']' <= c && c <= 'z':
		goto yys4
	case c == '{':
		goto yys15
	case c == '|':
		goto yys16
	case '}' <= c && c <= '\U0010ffff':
		goto yys4
	default:
		goto yyfinish
	}
yys2:
	c = yylex.Input()
	switch {
	case '\x00' <= c && c <= '\t':
		goto yys4
	case '\v' <= c && c <= '!':
		goto yys4
	case c == '"':
		goto yys17
	case '#' <= c && c <= '[':
		goto yys4
	case c == '\\':
		goto yys14
	case ']' <= c && c <= '\U0010ffff':
		goto yys4
	default:
		goto yyfinish
	}
yys3:
	c = yylex.Input()
	switch {
	case '\x00' <= c && c <= '\t':
		goto yys18
	case '\v' <= c && c <= '[':
		goto yys18
	case c == '\\':
		goto yys19
	case c == ']':
		goto yys20
	case '^' <= c && c <= '\U0010ffff':
		goto yys18
	default:
		goto yyfinish
	}
yys4:
	yyacc = 16
	yyleng = len(yylex.Token)
	goto yyfinish
yys5:
	yyacc = 10
	yyleng = len(yylex.Token)
	goto yyfinish
yys6:
	yyacc = 11
	yyleng = len(yylex.Token)
	goto yyfinish
yys7:
	yyacc = 8
	yyleng = len(yylex.Token)
	goto yyfinish
yys8:
	yyacc = 0
	yyleng = len(yylex.Token)
	goto yyfinish
yys9:
	yyacc = 1
	yyleng = len(yylex.Token)
	goto yyfinish
yys10:
	yyacc = 5
	yyleng = len(yylex.Token)
	goto yyfinish
yys11:
	yyacc = 6
	yyleng = len(yylex.Token)
	goto yyfinish
yys12:
	yyacc = 7
	yyleng = len(yylex.Token)
	goto yyfinish
yys13:
	yyacc = 2
	yyleng = len(yylex.Token)
	c = yylex.Input()
	switch {
	case c == '^':
		goto yys21
	default:
		goto yyfinish
	}
yys14:
	yyacc = 16
	yyleng = len(yylex.Token)
	c = yylex.Input()
	switch {
	case '\x00' <= c && c <= '\t':
		goto yys22
	case '\v' <= c && c <= '/':
		goto yys22
	case '0' <= c && c <= '7':
		goto yys23
	case '8' <= c && c <= 'T':
		goto yys22
	case c == 'U':
		goto yys24
	case 'V' <= c && c <= 't':
		goto yys22
	case c == 'u':
		goto yys25
	case 'v' <= c && c <= 'w':
		goto yys22
	case c == 'x':
		goto yys26
	case 'y' <= c && c <= '\U0010ffff':
		goto yys22
	default:
		goto yyfinish
	}
yys15:
	yyacc = 16
	yyleng = len(yylex.Token)
	c = yylex.Input()
	switch {
	case 'A' <= c && c <= 'Z':
		goto yys27
	case c == '_':
		goto yys27
	case 'a' <= c && c <= 'z':
		goto yys27
	default:
		goto yyfinish
	}
yys16:
	yyacc = 4
	yyleng = len(yylex.Token)
	goto yyfinish
yys17:
	yyacc = 12
	yyleng = len(yylex.Token)
	goto yyfinish
yys18:
	yyacc = 16
	yyleng = len(yylex.Token)
	c = yylex.Input()
	switch {
	case c == '-':
		goto yys28
	default:
		goto yyfinish
	}
yys19:
	yyacc = 16
	yyleng = len(yylex.Token)
	c = yylex.Input()
	switch {
	case '\x00' <= c && c <= '\t':
		goto yys29
	case '\v' <= c && c <= ',':
		goto yys29
	case c == '-':
		goto yys30
	case '.' <= c && c <= '/':
		goto yys29
	case '0' <= c && c <= '7':
		goto yys31
	case '8' <= c && c <= 'T':
		goto yys29
	case c == 'U':
		goto yys32
	case 'V' <= c && c <= 't':
		goto yys29
	case c == 'u':
		goto yys33
	case 'v' <= c && c <= 'w':
		goto yys29
	case c == 'x':
		goto yys34
	case 'y' <= c && c <= '\U0010ffff':
		goto yys29
	default:
		goto yyfinish
	}
yys20:
	yyacc = 13
	yyleng = len(yylex.Token)
	goto yyfinish
yys21:
	yyacc = 3
	yyleng = len(yylex.Token)
	goto yyfinish
yys22:
	yyacc = 15
	yyleng = len(yylex.Token)
	goto yyfinish
yys23:
	yyacc = 15
	yyleng = len(yylex.Token)
	c = yylex.Input()
	switch {
	case '0' <= c && c <= '7':
		goto yys35
	default:
		goto yyfinish
	}
yys24:
	yyacc = 15
	yyleng = len(yylex.Token)
	c = yylex.Input()
	switch {
	case '0' <= c && c <= '9':
		goto yys36
	case 'A' <= c && c <= 'F':
		goto yys36
	case 'a' <= c && c <= 'f':
		goto yys36
	default:
		goto yyfinish
	}
yys25:
	yyacc = 15
	yyleng = len(yylex.Token)
	c = yylex.Input()
	switch {
	case '0' <= c && c <= '9':
		goto yys37
	case 'A' <= c && c <= 'F':
		goto yys37
	case 'a' <= c && c <= 'f':
		goto yys37
	default:
		goto yyfinish
	}
yys26:
	yyacc = 15
	yyleng = len(yylex.Token)
	c = yylex.Input()
	switch {
	case '0' <= c && c <= '9':
		goto yys38
	case 'A' <= c && c <= 'F':
		goto yys38
	case 'a' <= c && c <= 'f':
		goto yys38
	default:
		goto yyfinish
	}
yys27:
	c = yylex.Input()
	switch {
	case '0' <= c && c <= '9':
		goto yys27
	case 'A' <= c && c <= 'Z':
		goto yys27
	case c == '_':
		goto yys27
	case 'a' <= c && c <= 'z':
		goto yys27
	case c == '}':
		goto yys39
	default:
		goto yyfinish
	}
yys28:
	c = yylex.Input()
	switch {
	case '\x00' <= c && c <= '\t':
		goto yys40
	case '\v' <= c && c <= '[':
		goto yys40
	case c == '\\':
		goto yys41
	case '^' <= c && c <= '\U0010ffff':
		goto yys40
	default:
		goto yyfinish
	}
yys29:
	yyacc = 15
	yyleng = len(yylex.Token)
	c = yylex.Input()
	switch {
	case c == '-':
		goto yys28
	default:
		goto yyfinish
	}
yys30:
	yyacc = 15
	yyleng = len(yylex.Token)
	c = yylex.Input()
	switch {
	case '\x00' <= c && c <= '\t':
		goto yys40
	case '\v' <= c && c <= ',':
		goto yys40
	case c == '-':
		goto yys42
	case '.' <= c && c <= '[':
		goto yys40
	case c == '\\':
		goto yys41
	case '^' <= c && c <= '\U0010ffff':
		goto yys40
	default:
		goto yyfinish
	}
yys31:
	yyacc = 15
	yyleng = len(yylex.Token)
	c = yylex.Input()
	switch {
	case c == '-':
		goto yys28
	case '0' <= c && c <= '7':
		goto yys43
	default:
		goto yyfinish
	}
yys32:
	yyacc = 15
	yyleng = len(yylex.Token)
	c = yylex.Input()
	switch {
	case c == '-':
		goto yys28
	case '0' <= c && c <= '9':
		goto yys44
	case 'A' <= c && c <= 'F':
		goto yys44
	case 'a' <= c && c <= 'f':
		goto yys44
	default:
		goto yyfinish
	}
yys33:
	yyacc = 15
	yyleng = len(yylex.Token)
	c = yylex.Input()
	switch {
	case c == '-':
		goto yys28
	case '0' <= c && c <= '9':
		goto yys45
	case 'A' <= c && c <= 'F':
		goto yys45
	case 'a' <= c && c <= 'f':
		goto yys45
	default:
		goto yyfinish
	}
yys34:
	yyacc = 15
	yyleng = len(yylex.Token)
	c = yylex.Input()
	switch {
	case c == '-':
		goto yys28
	case '0' <= c && c <= '9':
		goto yys46
	case 'A' <= c && c <= 'F':
		goto yys46
	case 'a' <= c && c <= 'f':
		goto yys46
	default:
		goto yyfinish
	}
yys35:
	c = yylex.Input()
	switch {
	case '0' <= c && c <= '7':
		goto yys22
	default:
		goto yyfinish
	}
yys36:
	c = yylex.Input()
	switch {
	case '0' <= c && c <= '9':
		goto yys47
	case 'A' <= c && c <= 'F':
		goto yys47
	case 'a' <= c && c <= 'f':
		goto yys47
	default:
		goto yyfinish
	}
yys37:
	c = yylex.Input()
	switch {
	case '0' <= c && c <= '9':
		goto yys48
	case 'A' <= c && c <= 'F':
		goto yys48
	case 'a' <= c && c <= 'f':
		goto yys48
	default:
		goto yyfinish
	}
yys38:
	c = yylex.Input()
	switch {
	case '0' <= c && c <= '9':
		goto yys22
	case 'A' <= c && c <= 'F':
		goto yys22
	case 'a' <= c && c <= 'f':
		goto yys22
	default:
		goto yyfinish
	}
yys39:
	yyacc = 9
	yyleng = len(yylex.Token)
	goto yyfinish
yys40:
	yyacc = 14
	yyleng = len(yylex.Token)
	goto yyfinish
yys41:
	yyacc = 14
	yyleng = len(yylex.Token)
	c = yylex.Input()
	switch {
	case '\x00' <= c && c <= '\t':
		goto yys40
	case '\v' <= c && c <= '/':
		goto yys40
	case '0' <= c && c <= '7':
		goto yys49
	case '8' <= c && c <= 'T':
		goto yys40
	case c == 'U':
		goto yys50
	case 'V' <= c && c <= 't':
		goto yys40
	case c == 'u':
		goto yys51
	case 'v' <= c && c <= 'w':
		goto yys40
	case c == 'x':
		goto yys52
	case 'y' <= c && c <= '\U0010ffff':
		goto yys40
	default:
		goto yyfinish
	}
yys42:
	yyacc = 14
	yyleng = len(yylex.Token)
	c = yylex.Input()
	switch {
	case '\x00' <= c && c <= '\t':
		goto yys40
	case '\v' <= c && c <= '[':
		goto yys40
	case c == '\\':
		goto yys41
	case '^' <= c && c <= '\U0010ffff':
		goto yys40
	default:
		goto yyfinish
	}
yys43:
	c = yylex.Input()
	switch {
	case '0' <= c && c <= '7':
		goto yys29
	default:
		goto yyfinish
	}
yys44:
	c = yylex.Input()
	switch {
	case '0' <= c && c <= '9':
		goto yys53
	case 'A' <= c && c <= 'F':
		goto yys53
	case 'a' <= c && c <= 'f':
		goto yys53
	default:
		goto yyfinish
	}
yys45:
	c = yylex.Input()
	switch {
	case '0' <= c && c <= '9':
		goto yys54
	case 'A' <= c && c <= 'F':
		goto yys54
	case 'a' <= c && c <= 'f':
		goto yys54
	default:
		goto yyfinish
	}
yys46:
	c = yylex.Input()
	switch {
	case '0' <= c && c <= '9':
		goto yys29
	case 'A' <= c && c <= 'F':
		goto yys29
	case 'a' <= c && c <= 'f':
		goto yys29
	default:
		goto yyfinish
	}
yys47:
	c = yylex.Input()
	switch {
	case '0' <= c && c <= '9':
		goto yys55
	case 'A' <= c && c <= 'F':
		goto yys55
	case 'a' <= c && c <= 'f':
		goto yys55
	default:
		goto yyfinish
	}
yys48:
	c = yylex.Input()
	switch {
	case '0' <= c && c <= '9':
		goto yys38
	case 'A' <= c && c <= 'F':
		goto yys38
	case 'a' <= c && c <= 'f':
		goto yys38
	default:
		goto yyfinish
	}
yys49:
	yyacc = 14
	yyleng = len(yylex.Token)
	c = yylex.Input()
	switch {
	case '0' <= c && c <= '7':
		goto yys56
	default:
		goto yyfinish
	}
yys50:
	yyacc = 14
	yyleng = len(yylex.Token)
	c = yylex.Input()
	switch {
	case '0' <= c && c <= '9':
		goto yys57
	case 'A' <= c && c <= 'F':
		goto yys57
	case 'a' <= c && c <= 'f':
		goto yys57
	default:
		goto yyfinish
	}
yys51:
	yyacc = 14
	yyleng = len(yylex.Token)
	c = yylex.Input()
	switch {
	case '0' <= c && c <= '9':
		goto yys58
	case 'A' <= c && c <= 'F':
		goto yys58
	case 'a' <= c && c <= 'f':
		goto yys58
	default:
		goto yyfinish
	}
yys52:
	yyacc = 14
	yyleng = len(yylex.Token)
	c = yylex.Input()
	switch {
	case '0' <= c && c <= '9':
		goto yys59
	case 'A' <= c && c <= 'F':
		goto yys59
	case 'a' <= c && c <= 'f':
		goto yys59
	default:
		goto yyfinish
	}
yys53:
	c = yylex.Input()
	switch {
	case '0' <= c && c <= '9':
		goto yys60
	case 'A' <= c && c <= 'F':
		goto yys60
	case 'a' <= c && c <= 'f':
		goto yys60
	default:
		goto yyfinish
	}
yys54:
	c = yylex.Input()
	switch {
	case '0' <= c && c <= '9':
		goto yys46
	case 'A' <= c && c <= 'F':
		goto yys46
	case 'a' <= c && c <= 'f':
		goto yys46
	default:
		goto yyfinish
	}
yys55:
	c = yylex.Input()
	switch {
	case '0' <= c && c <= '9':
		goto yys61
	case 'A' <= c && c <= 'F':
		goto yys61
	case 'a' <= c && c <= 'f':
		goto yys61
	default:
		goto yyfinish
	}
yys56:
	c = yylex.Input()
	switch {
	case '0' <= c && c <= '7':
		goto yys40
	default:
		goto yyfinish
	}
yys57:
	c = yylex.Input()
	switch {
	case '0' <= c && c <= '9':
		goto yys62
	case 'A' <= c && c <= 'F':
		goto yys62
	case 'a' <= c && c <= 'f':
		goto yys62
	default:
		goto yyfinish
	}
yys58:
	c = yylex.Input()
	switch {
	case '0' <= c && c <= '9':
		goto yys63
	case 'A' <= c && c <= 'F':
		goto yys63
	case 'a' <= c && c <= 'f':
		goto yys63
	default:
		goto yyfinish
	}
yys59:
	c = yylex.Input()
	switch {
	case '0' <= c && c <= '9':
		goto yys40
	case 'A' <= c && c <= 'F':
		goto yys40
	case 'a' <= c && c <= 'f':
		goto yys40
	default:
		goto yyfinish
	}
yys60:
	c = yylex.Input()
	switch {
	case '0' <= c && c <= '9':
		goto yys64
	case 'A' <= c && c <= 'F':
		goto yys64
	case 'a' <= c && c <= 'f':
		goto yys64
	default:
		goto yyfinish
	}
yys61:
	c = yylex.Input()
	switch {
	case '0' <= c && c <= '9':
		goto yys37
	case 'A' <= c && c <= 'F':
		goto yys37
	case 'a' <= c && c <= 'f':
		goto yys37
	default:
		goto yyfinish
	}
yys62:
	c = yylex.Input()
	switch {
	case '0' <= c && c <= '9':
		goto yys65
	case 'A' <= c && c <= 'F':
		goto yys65
	case 'a' <= c && c <= 'f':
		goto yys65
	default:
		goto yyfinish
	}
yys63:
	c = yylex.Input()
	switch {
	case '0' <= c && c <= '9':
		goto yys59
	case 'A' <= c && c <= 'F':
		goto yys59
	case 'a' <= c && c <= 'f':
		goto yys59
	default:
		goto yyfinish
	}
yys64:
	c = yylex.Input()
	switch {
	case '0' <= c && c <= '9':
		goto yys45
	case 'A' <= c && c <= 'F':
		goto yys45
	case 'a' <= c && c <= 'f':
		goto yys45
	default:
		goto yyfinish
	}
yys65:
	c = yylex.Input()
	switch {
	case '0' <= c && c <= '9':
		goto yys66
	case 'A' <= c && c <= 'F':
		goto yys66
	case 'a' <= c && c <= 'f':
		goto yys66
	default:
		goto yyfinish
	}
yys66:
	c = yylex.Input()
	switch {
	case '0' <= c && c <= '9':
		goto yys58
	case 'A' <= c && c <= 'F':
		goto yys58
	case 'a' <= c && c <= 'f':
		goto yys58
	default:
		goto yyfinish
	}
yyfinish:
	yylex.Back(len(yylex.Token) - yyleng)
	yytext := yylex.Token[:]
	switch yyacc {
	case -1:
		_ = yytext
		return 0
	case 0:
		return LPAREN

	case 1:
		return RPAREN

	case 2:
		BEGIN(CHARSET)
		return LBRACK

	case 3:
		BEGIN(CHARSET)
		return LBRACK_CARET

	case 4:
		return PIPE

	case 5:
		return STAR

	case 6:
		return PLUS

	case 7:
		return DOT

	case 8:
		BEGIN(LITERAL)
		return LPAREN

	case 9:
		{
			name := string(yytext[1 : len(yytext)-1])
			frag, ok := yylex.Defs[name]
			if ok {
				yylval.frag = frag
				return NAME
			}
			yylex.Error(fmt.Sprintf("undefined %q", name))
		}
	case 10:
		{
			if yylex.NoSpace {
				return 0
			}
			yylval.rng = Range{yytext[0], yytext[0]}
			return CHAR
		}
	case 11:
		{
			if yylex.NoNewline {
				return 0
			}
			yylval.rng = Range{yytext[0], yytext[0]}
			return CHAR
		}
	case 12:
		BEGIN(INITIAL)
		return RPAREN

	case 13:
		BEGIN(INITIAL)
		return RBRACK

	case 14:
		{
			r1, yytext := unquote(yytext)
			r2, _ := unquote(yytext[1:])
			yylval.rng = Range{r1, r2}
			return CHAR
		}
	case 15:
		{
			r, _ := unquote(yytext)
			yylval.rng = Range{r, r}
			return CHAR
		}
	case 16:
		{
			yylval.rng = Range{yytext[0], yytext[0]}
			return CHAR
		}
	}
	goto yystart
}
