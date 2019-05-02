package main

import (
	"bufio"
	"flag"
	"log"
	"os"

	"github.com/z-rui/lex/re"
)

var (
	genNFA = flag.Bool("n", false, "Generate NFA instead of DFA")
	minify = flag.Bool("m", true, "Minify DFA")
)

func main() {
	flag.Parse()

	if flag.NArg() != 1 {
		log.Fatal("need exactly 1 argument.")
	}
	frag := re.ParseString(flag.Arg(0))
	if frag == nil {
		return
	}
	a := frag.Canon()
	w := bufio.NewWriter(os.Stdout)
	if *genNFA {
		a.WriteDotTo(w)
	} else {
		d := re.NFA2DFA(a)
		if *minify {
			d = d.Minimize()
			d.WriteDotTo(w, false)
		} else {
			d.WriteDotTo(w, true)
		}
	}
	w.Flush()
}
