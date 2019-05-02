package re

import (
	"fmt"
	"os"
	"strconv"
)

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
	case 'a': r = '\a'
	case 'b': r = '\b'
	case 't': r = '\t'
	case 'n': r = '\n'
	case 'v': r = '\v'
	case 'f': r = '\f'
	case 'r': r = '\r'
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

%%

%s LITERAL CHARSET

x		[0-9A-Fa-f]
o		[0-7]
esc		\\({o}{o}{o}|x{x}{x}|u{x}{x}{x}{x}|U{x}{x}{x}{x}{x}{x}{x}{x}|.)
id              [A-Za-z_][A-Za-z_0-9]*

%%

<INITIAL>"("    return LPAREN
<INITIAL>")"    return RPAREN
<INITIAL>"["    BEGIN(CHARSET); return LBRACK
<INITIAL>"[^"   BEGIN(CHARSET); return LBRACK_CARET
<INITIAL>"|"    return PIPE
<INITIAL>"*"    return STAR
<INITIAL>"+"    return PLUS
<INITIAL>"."    return DOT
<INITIAL>\"     BEGIN(LITERAL); return LPAREN 
<INITIAL>\{{id}\} {
	name := string(yytext[1:len(yytext)-1])
	frag, ok := yylex.Defs[name]
	if ok {
		yylval.frag = frag
		return NAME
	}
	yylex.Error(fmt.Sprintf("undefined %q", name))
}
<INITIAL>[ \t\r\n\v\f] {
	r := yytext[0]
	if yylex.Whitespace & (1 << uint(r)) == 0 {
		yylex.Back(1)
		return 0
	}
	yylval.rng = Range{r, r}
	return CHAR
}

<LITERAL>\"	BEGIN(INITIAL); return RPAREN

<CHARSET>"]"    BEGIN(INITIAL); return RBRACK
<CHARSET>({esc}|[^\]])-({esc}|[^\]]) {
	r1, yytext := unquote(yytext)
	r2, _ := unquote(yytext[1:])
	yylval.rng = Range{r1, r2}
	return CHAR
}

{esc}           {
	r, _ := unquote(yytext)
	yylval.rng = Range{r, r}
	return CHAR
}
.               {
	yylval.rng = Range{yytext[0], yytext[0]}
	return CHAR
}
