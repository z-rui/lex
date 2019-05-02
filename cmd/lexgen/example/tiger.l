package main

import "fmt"

%%

%x COMMENT STRING

ws         [ \t\n\r\f\v]
digit      [0-9]
ident      [A-Za-z][A-Za-z0-9_]*

%%
	commentLevel := 0
	strBuf := []rune{}

type     yylex.Printf("TYPE")
var      yylex.Printf("VAR")
function yylex.Printf("FUNCTION")
break    yylex.Printf("BREAK")
of       yylex.Printf("OF")
end      yylex.Printf("END")
in       yylex.Printf("IN")
nil      yylex.Printf("NIL")
let      yylex.Printf("LET")
do       yylex.Printf("DO")
to       yylex.Printf("TO")
for      yylex.Printf("FOR")
while    yylex.Printf("WHILE")
else     yylex.Printf("ELSE")
then     yylex.Printf("THEN")
if       yylex.Printf("IF")
array    yylex.Printf("ARRAY")

":="     yylex.Printf("ASSIGN")
"|"      yylex.Printf("OR")
"&"      yylex.Printf("AND")
">="     yylex.Printf("GE")
">"      yylex.Printf("GT")
"<="     yylex.Printf("LE")
"<"      yylex.Printf("LT")
"<>"     yylex.Printf("NEQ")
"="      yylex.Printf("EQ")
"/"      yylex.Printf("DIVIDE")
"*"      yylex.Printf("TIMES")
"-"      yylex.Printf("MINUS")
"+"      yylex.Printf("PLUS")
"."      yylex.Printf("DOT")
"}"      yylex.Printf("RBRACE")
"{"      yylex.Printf("LBRACE")
"]"      yylex.Printf("RBRACK")
"["      yylex.Printf("LBRACK")
")"      yylex.Printf("RPAREN")
"("      yylex.Printf("LPAREN")
";"      yylex.Printf("SEMICOLON")
":"      yylex.Printf("COLON")
","      yylex.Printf("COMMA")

{ident}  yylex.Printf("IDENT(%s)", string(yytext))
{ws}+    /* skip whitespace */
{digit}+ yylex.Printf("INT(%s)", string(yytext))

"/*"     {
	commentLevel = 1
	BEGIN(COMMENT)
}
<COMMENT>/+"*" commentLevel++
<COMMENT>"*"+/ {
	commentLevel--
	if commentLevel < 1 {
		BEGIN(INITIAL)
	}
}
<COMMENT>(/+[^*/]*|"*"+[^*/]*|[^/*]+)     // skip

\"       strBuf = strBuf[:0]; BEGIN(STRING)
<STRING>\" {
	yylex.Printf("STRING(%q)", string(strBuf))
	BEGIN(INITIAL)
}
<STRING>\\[nt"\\] {
	r := yytext[1]
	switch r {
		case 'n': r = '\n'
		case 't': r = '\t'
		case '"': r = '"'
	}
	strBuf = append(strBuf, r)
}
<STRING>\\{digit}{digit}{digit} { // \ddd
	n := rune(0)
	for _, r := range yytext[1:] {
		n = n * 10 + (r - '0')
	}
	strBuf = append(strBuf, n)
}
<STRING>\\\^[@A-Z\[\\\]\^_] { // \^c
	r := yytext[2] - 64
	strBuf = append(strBuf, r)
}
<STRING>\\{ws}+\\               /* \f...f\ */
<STRING>\\  yylex.Error("invalid escape sequence")
<STRING>\n  yylex.Error("newline in string literal"); BEGIN(INITIAL)
<STRING>[^\\\n]+ {
	strBuf = append(strBuf, yytext...)
}

.        yylex.Error(fmt.Sprintf("Unknown character %q", yytext[0]))
