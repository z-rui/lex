package main

import (
	"bufio"
	"flag"
	"log"
	"os"
)

var prefix = flag.String("p", "yy", "prefix")

func main() {
	flag.Parse()

	if flag.NArg() != 1 {
		log.Fatal("need exactly 1 argument")
	}
	//out := os.Stdout
	out, err := os.Create("lex.yy.go")
	if err != nil {
		log.Fatal(err)
	}
	defer out.Close()
	w := bufio.NewWriter(out)
	defer w.Flush()
	gen := &LexGen{
		Prefix: *prefix,
		Out:    w,
	}
	gen.Run(flag.Arg(0))

	gen.DumpDFA()
}
