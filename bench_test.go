package monster

import (
	"math/rand"
	"testing"
)

var intProd = `
integer : range(-10000, 10000000).
`

func BenchmarkInteger(b *testing.B) {
	conf := make(map[string]interface{})
	start := ParseText([]byte(intProd), conf)

	nonterminals, root := Build(start)
	c := map[string]interface{}{
		"_random":       rand.New(rand.NewSource(int64(10))),
		"_nonterminals": nonterminals,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		root.Generate(c)
	}
}

var charProd = `
char    : "a"
        | "b"
        | "c"
        | "d"
        | "e"
        | "f"
        | "g"
        | "h"
        | "i"
        | "j"
        | "k"
        | "l"
        | "m"
        | "n"
        | "o"
        | "p"
        | "q"
        | "r"
        | "s"
        | "t"
        | "u"
        | "v"
        | "w"
        | "x"
        | "y"
        | "z".
`

func BenchmarkChar(b *testing.B) {
	conf := make(map[string]interface{})
	start := ParseText([]byte(charProd), conf)

	nonterminals, root := Build(start)
	c := map[string]interface{}{
		"_random":       rand.New(rand.NewSource(int64(10))),
		"_nonterminals": nonterminals,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		root.Generate(c)
	}
}

var strProd = `
string  : char
        | string char {200}
        | string char {200}
        | string char {200}.
` + charProd

func BenchmarkString(b *testing.B) {
	conf := make(map[string]interface{})
	start := ParseText([]byte(strProd), conf)

	nonterminals, root := Build(start)
	c := map[string]interface{}{
		"_random":       rand.New(rand.NewSource(int64(10))),
		"_nonterminals": nonterminals,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		root.Generate(c)
	}
}

var floatProd = `
float : rangef(-1000.0, 1000.0).
`

func BenchmarkFloat(b *testing.B) {
	conf := make(map[string]interface{})
	start := ParseText([]byte(floatProd), conf)

	nonterminals, root := Build(start)
	c := map[string]interface{}{
		"_random":       rand.New(rand.NewSource(int64(10))),
		"_nonterminals": nonterminals,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		root.Generate(c)
	}
}
