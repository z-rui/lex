package re

// Tokens
const (
	_ = iota + 2 // eof, error, unk
	CHAR
	NAME
	LPAREN
	RPAREN
	LQUOTE
	RQUOTE
	LBRACK
	LBRACK_CARET
	RBRACK
	ALTER
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
	"LQUOTE",
	"RQUOTE",
	"LBRACK",
	"LBRACK_CARET",
	"RBRACK",
	"ALTER",
	"STAR",
	"PLUS",
	"DASH",
	"DOT",
}

const yyAccept = 8
const yyLast = 17

// Parse tables
var yyR1 = [...]int{
	0, 21, 24, 24, 23, 23, 25, 25, 22, 22,
	22, 22, 22, 22, 22, 22, 22, 18, 18, 18,
	19, 19, 20, 20, 17, 17,
}

var yyR2 = [...]int{
	2, 1, 1, 3, 0, 1, 1, 2, 1, 1,
	1, 3, 3, 3, 3, 2, 2, 1, 1, 2,
	0, 3, 2, 2, 1, 1,
}

var yyReduce = [...]int{
	4, 8, 10, 4, 4, 20, 20, 9, 0, 6,
	2, 1, 5, 0, 0, 0, 17, 18, 0, 15,
	16, 4, 7, 13, 14, 11, 24, 25, 22, 23,
	19, 12, 3, 21,
}

var yyGoto = [...]int{
	28, 15, 16, 17, 8, 9, 10, 11, 12,
}

var yyAction = [...]int{
	1, 2, 3, 26, 4, 26, 5, 6, 13, 14,
	29, 19, 20, 7, 25, 27, 21, 27, 22, 31,
	32, 18, 30, 33, 24, 23,
}

var yyCheck = [...]int{
	3, 4, 5, 3, 7, 3, 9, 10, 3, 4,
	3, 13, 14, 16, 11, 15, 12, 15, 12, 11,
	21, 6, 15, 30, 8, 6,
}

var yyPact = [...]int{
	-3, 26, 26, -3, -3, 26, 26, 26, 26, -2,
	26, 4, -3, 19, 16, 3, 2, 7, 8, 26,
	26, -3, -2, 26, 26, 26, 26, 26, 26, 26,
	0, 26, 26, 26,
}

var yyPgoto = [...]int{
	-7, 15, 26, 26, 5, 6, -1, 26, 26,
}

type yySymType struct {
	yys int // state

	r      rune
	ranges []Range
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
					case n-1:
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
 yyval.frag = Literal([]Range{{yyD[0].r, yyD[0].r}}, false) 
	case 9:
 yyval.frag = Literal([]Range{{'\n','\n'}}, true) 
	case 10:
 yyval.frag = yyD[0].frag.Clone() 
	case 11:
 yyval.frag = Literal(yyD[1].ranges, false) 
	case 12:
 yyval.frag = Literal(yyD[1].ranges, true) 
	case 13:
 yyval.frag = yyD[1].frag 
	case 14:
 yyval.frag = yyD[1].frag 
	case 15:
 yyval.frag = Kleene(yyD[0].frag) 
	case 16:
 yyval.frag = KleenePlus(yyD[0].frag) 
	case 19:
 yyval.ranges = append(yyD[0].ranges, Range{'-', '-'}) 
	case 20:
 yyval.ranges = nil 
	case 21:
 yyD[0].ranges[len(yyD[0].ranges)-1].Last = yyD[2].r; yyval.ranges = yyD[0].ranges 
	case 22:
 yyval.ranges = append(yyD[0].ranges, Range{yyD[1].r, yyD[1].r}) 
	case 23:
 yyval.ranges = append(yyD[0].ranges, Range{yyD[1].r, yyD[1].r}) 
	case 25:
 yyval.r = '-' 
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
