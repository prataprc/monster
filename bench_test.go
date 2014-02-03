//  Copyright (c) 2013 Couchbase, Inc.
//  Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file
//  except in compliance with the License. You may obtain a copy of the License at
//    http://www.apache.org/licenses/LICENSE-2.0
//  Unless required by applicable law or agreed to in writing, software distributed under the
//  License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
//  either express or implied. See the License for the specific language governing permissions
//  and limitations under the License.

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
	start, err := ParseText([]byte(intProd), conf)
	if err != nil {
		b.Fatal(err)
	}

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
	start, err := ParseText([]byte(charProd), conf)
	if err != nil {
		b.Fatal(err)
	}

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
	start, err := ParseText([]byte(strProd), conf)
	if err != nil {
		b.Fatal(err)
	}

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
	start, err := ParseText([]byte(floatProd), conf)
	if err != nil {
		b.Fatal(err)
	}

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
