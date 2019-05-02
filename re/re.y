package re

%%

%union {
	rng    Range
	ranges RangeSet
	frag   *Frag
	frags  []*Frag
}

%token <rng>  CHAR
%token <frag> NAME
%token LPAREN RPAREN LBRACK LBRACK_CARET RBRACK
       PIPE STAR PLUS DASH DOT

%type <ranges> charset
%type <frag> regex primary concats_opt
%type <frags> alters concats

%%

regex : alters { $$ = Alter($1) }

alters
: concats_opt { $$ = []*Frag{$1} }
| alters PIPE concats_opt { $$ = append($$, $3) }
;

concats_opt
: /* epsilon */ { $$ = Empty() }
| concats { $$ = Concat($1) }

concats
: primary  { $$ = []*Frag{$1} }
| concats primary { $$ = append($1, $2) }
;

primary
: CHAR { $$ = Literal(RangeSet{$1}, false) }
| DOT { $$ = Literal(RangeSet{{'\n','\n'}}, true) }
| NAME { $$ = $1.Clone() }
| LBRACK charset RBRACK { $$ = Literal($2, false) }
| LBRACK_CARET charset RBRACK
{
	if len($2) == 0 { // [^]
		$$ = Literal(RangeSet{{'^','^'}}, false)
	} else {
		$$ = Literal($2, true)
	}
}
| LPAREN regex RPAREN { $$ = $2 }
| primary STAR { $$ = Kleene($1) }
| primary PLUS { $$ = KleenePlus($1) }
;

charset
: /* epsilon */ { $$ = nil }
| charset CHAR  { $$ = append($1, $2) }
;
