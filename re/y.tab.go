package re

// Tokens
const (
	_ = iota + 2 // eof, error, unk
	CHAR
	NAME
	LPAREN
	RPAREN
	LBRACK
	LBRACK_CARET
	RBRACK
	PIPE
	STAR
	PLUS
	DASH
	DOT
)

var yyName = []string{
	"$end",
	"error",
	"$unk",
	"CHAR",
	"NAME",
	"LPAREN",
	"RPAREN",
	"LBRACK",
	"LBRACK_CARET",
	"RBRACK",
	"PIPE",
	"STAR",
	"PLUS",
	"DASH",
	"DOT",
}

const yyAccept = 7
const yyLast = 15

// Parse tables
var yyR1 = [...]int{
	0, 16, 19, 19, 18, 18, 20, 20, 17, 17,
	17, 17, 17, 17, 17, 17, 15, 15,
}

var yyR2 = [...]int{
	2, 1, 1, 3, 0, 1, 1, 2, 1, 1,
	1, 3, 3, 3, 2, 2, 0, 2,
}

var yyReduce = [...]int{
	4, 8, 10, 4, 16, 16, 9, 0, 6, 2,
	1, 5, 0, 0, 0, 14, 15, 4, 7, 13,
	17, 11, 12, 3,
}

var yyGoto = [...]int{
	13, 7, 8, 9, 10, 11,
}

var yyAction = [...]int{
	1, 2, 3, 20, 4, 5, 20, 15, 16, 21,
	12, 6, 22, 18, 19, 14, 17, 0, 23,
}

var yyCheck = [...]int{
	3, 4, 5, 3, 7, 8, 3, 11, 12, 9,
	3, 14, 9, 11, 6, 5, 10, -1, 17,
}

var yyPact = [...]int{
	-3, 19, 19, -3, 19, 19, 19, 19, -4, 19,
	6, -3, 8, 0, 3, 19, 19, -3, -4, 19,
	19, 19, 19, 19,
}

var yyPgoto = [...]int{
	10, 7, 2, 1, 19, 19,
}

type yySymType struct {
	yys int // state

	rng    Range
	ranges RangeSet
	frag   *Frag
	frags  []*Frag
}

type yyLexer interface {
	Lex(*yySymType) int
	Error(string)
}

var yyDebug = 0 // debug info from parser

// yyParse read tokens from yylex and parses input.
// Returns result on success, or nil on failure.
func yyParse(yylex yyLexer) *yySymType {
	var (
		yyn, yyt int
		yystate  = 0
		yyerror  = 0
		yymajor  = -1
		yystack  []yySymType
		yyD      []yySymType // rhs of reduction
		yylval   yySymType   // lexcial value from lexer
		yyval    yySymType   // value to be pushed onto stack
	)
	goto yyaction
yystack:
	yyval.yys = yystate
	yystack = append(yystack, yyval)
	yystate = yyn
	if yyDebug >= 2 {
		println("\tGOTO state", yyn)
	}
yyaction:
	// look up shift or reduce
	yyn = int(yyPact[yystate])
	if yyn == len(yyAction) && yystate != yyAccept { // simple state
		goto yydefault
	}
	if yymajor < 0 {
		yymajor = yylex.Lex(&yylval)
		if yyDebug >= 1 {
			println("In state", yystate)
		}
		if yyDebug >= 2 {
			println("\tInput token", yyName[yymajor])
		}
	}
	yyn += yymajor
	if 0 <= yyn && yyn < len(yyAction) && int(yyCheck[yyn]) == yymajor {
		yyn = int(yyAction[yyn])
		if yyn <= 0 {
			yyn = -yyn
			goto yyreduce
		}
		if yyDebug >= 1 {
			println("\tSHIFT token", yyName[yymajor])
		}
		if yyerror > 0 {
			yyerror--
		}
		yymajor = -1
		yyval = yylval
		goto yystack
	}
yydefault:
	yyn = int(yyReduce[yystate])
yyreduce:
	if yyn == 0 {
		if yymajor == 0 && yystate == yyAccept {
			if yyDebug >= 1 {
				println("\tACCEPT!")
			}
			return &yystack[0]
		}
		switch yyerror {
		case 0: // new error
			if yyDebug >= 1 {
				println("\tERROR!")
			}
			msg := "unexpected " + yyName[yymajor]
			var expect []int
			if yyReduce[yystate] == 0 {
				yyn = yyPact[yystate] + 3
				for i := 3; i < yyLast; i++ {
					if 0 <= yyn && yyn < len(yyAction) && yyCheck[yyn] == i && yyAction[yyn] != 0 {
						expect = append(expect, i)
						if len(expect) > 4 {
							break
						}
					}
					yyn++
				}
			}
			if n := len(expect); 0 < n && n <= 4 {
				for i, tok := range expect {
					switch i {
					case 0:
						msg += ", expecting "
					case n - 1:
						msg += " or "
					default:
						msg += ", "
					}
					msg += yyName[tok]
				}
			}
			yylex.Error(msg)
			fallthrough
		case 1, 2: // partially recovered error
			for { // pop states until error can be shifted
				yyn = int(yyPact[yystate]) + 1
				if 0 <= yyn && yyn < len(yyAction) && yyCheck[yyn] == 1 {
					yyn = yyAction[yyn]
					if yyn > 0 {
						break
					}
				}
				if len(yystack) == 0 {
					return nil
				}
				if yyDebug >= 2 {
					println("\tPopping state", yystate)
				}
				yystate = yystack[len(yystack)-1].yys
				yystack = yystack[:len(yystack)-1]
			}
			yyerror = 3
			if yyDebug >= 1 {
				println("\tSHIFT token error")
			}
			goto yystack
		default: // still waiting for valid tokens
			if yymajor == 0 { // no more tokens
				return nil
			}
			if yyDebug >= 1 {
				println("\tDISCARD token", yyName[yymajor])
			}
			yymajor = -1
			goto yyaction
		}
	}
	if yyDebug >= 1 {
		println("\tREDUCE rule", yyn)
	}
	yyt = len(yystack) - int(yyR2[yyn])
	yyD = yystack[yyt:]
	if len(yyD) > 0 { // pop items and restore state
		yyval = yyD[0]
		yystate = yyval.yys
		yystack = yystack[:yyt]
	}
	switch yyn { // Semantic actions

	case 1:
		yyval.frag = Alter(yyD[0].frags)
	case 2:
		yyval.frags = []*Frag{yyD[0].frag}
	case 3:
		yyval.frags = append(yyval.frags, yyD[2].frag)
	case 4:
		yyval.frag = Empty()
	case 5:
		yyval.frag = Concat(yyD[0].frags)
	case 6:
		yyval.frags = []*Frag{yyD[0].frag}
	case 7:
		yyval.frags = append(yyD[0].frags, yyD[1].frag)
	case 8:
		yyval.frag = Literal(RangeSet{yyD[0].rng}, false)
	case 9:
		yyval.frag = Literal(RangeSet{{'\n', '\n'}}, true)
	case 10:
		yyval.frag = yyD[0].frag.Clone()
	case 11:
		yyval.frag = Literal(yyD[1].ranges, false)
	case 12:

		if len(yyD[1].ranges) == 0 { // [^]
			yyval.frag = Literal(RangeSet{{'^', '^'}}, false)
		} else {
			yyval.frag = Literal(yyD[1].ranges, true)
		}

	case 13:
		yyval.frag = yyD[1].frag
	case 14:
		yyval.frag = Kleene(yyD[0].frag)
	case 15:
		yyval.frag = KleenePlus(yyD[0].frag)
	case 16:
		yyval.ranges = nil
	case 17:
		yyval.ranges = append(yyD[0].ranges, yyD[1].rng)
	}
	// look up goto
	yyt = int(yyR1[yyn]) - yyLast
	yyn = int(yyPgoto[yyt]) + yystate
	if 0 <= yyn && yyn < len(yyAction) &&
		int(yyCheck[yyn]) == yystate {
		yyn = int(yyAction[yyn])
	} else {
		yyn = int(yyGoto[yyt])
	}
	goto yystack
}
