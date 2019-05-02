package re

%%

%union {
	r      rune
	ranges []Range
	frag   *Frag
	frags  []*Frag
}

%token <r> CHAR
%token <frag> NAME
%token LPAREN RPAREN LQUOTE RQUOTE
       LBRACK LBRACK_CARET RBRACK
       ALTER STAR PLUS DASH DOT

%type <r> char
%type <ranges> charset charset1 charset2
%type <frag> regex primary concats_opt
%type <frags> alters concats

%%

regex : alters { $$ = Alter($1) }

alters
: concats_opt { $$ = []*Frag{$1} }
| alters ALTER concats_opt { $$ = append($$, $3) }
;

concats_opt
: /* epsilon */ { $$ = Empty() }
| concats { $$ = Concat($1) }

concats
: primary  { $$ = []*Frag{$1} }
| concats primary { $$ = append($1, $2) }
;

primary
: CHAR { $$ = Literal([]Range{{$1, $1}}, false) }
| DOT { $$ = Literal([]Range{{'\n','\n'}}, true) }
| NAME { $$ = $1.Clone() }
| LBRACK charset RBRACK { $$ = Literal($2, false) }
| LBRACK_CARET charset RBRACK { $$ = Literal($2, true) }
| LPAREN regex RPAREN { $$ = $2 }
| LQUOTE regex RQUOTE { $$ = $2 }
| primary STAR { $$ = Kleene($1) }
| primary PLUS { $$ = KleenePlus($1) }
;

/* dealing with charset is nasty because DASH can fallback to CHAR
   but our parser generator does not have this feature. */
charset
: charset1
| charset2
| charset2 DASH { $$ = append($1, Range{'-', '-'}) }
;

charset1 /* following DASH => CHAR */
: { $$ = nil }
| charset2 DASH char { $1[len($1)-1].Last = $3; $$ = $1 }
;

charset2 /* following DASH => DASH */
: charset1 char { $$ = append($1, Range{$2, $2}) }
| charset2 CHAR { $$ = append($1, Range{$2, $2}) }
;

char : CHAR | DASH { $$ = '-' } ;
