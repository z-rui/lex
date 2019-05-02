package re

import (
	"testing"
)

func TestParse(t *testing.T) {
	frag := ParseString("")

	if frag.Start != frag.Accept || frag.Accept.In != nil {
		t.Error("wrong NFA for epsilon")
	}

	frag = ParseString("a")
	if frag.Start != frag.Accept || frag.Accept.In.String() != "[a]" {
		t.Error("wrong NFA for single char:", frag.Accept.In)
	}

	frag = ParseString("\t")
	if frag.Start != frag.Accept || frag.Accept.In.String() != `[\t]` {
		t.Error("wrong NFA for single char:", frag.Accept.In)
	}

	frag = ParseString("[a-z]")
	if frag.Start != frag.Accept || frag.Accept.In.String() != "[a-z]" {
		t.Error("wrong NFA for charset:", frag.Accept.In)
	}

	frag = ParseString(`"*"`)
	if frag.Start != frag.Accept || frag.Accept.In.String() != `[*]` {
		t.Error("wrong NFA for literal:", frag.Accept.In)
	}
}
