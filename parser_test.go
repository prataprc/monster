//  Copyright (c) 2013 Couchbase, Inc.

package monster

import "testing"
import "fmt"
import "io/ioutil"
import "time"

import "github.com/prataprc/goparsec"
import "github.com/prataprc/monster/common"

var _ = fmt.Sprintf("dummy")

func TestEmpty(t *testing.T) {
	s := parsec.NewScanner([]byte(``))
	root, _ := Y(s)
	scope := root.(common.Scope)
	forms := scope["_globalForms"].([]*common.Form)
	nterms := scope["_nonterminals"].(common.NTForms)
	if len(forms) > 0 {
		t.Fatalf("Expected empty forms for empty prod %v", forms)
	} else if len(nterms) > 0 {
		t.Fatalf("Expected empty nterms for empty prod %v", nterms)
	}
}

func TestString(t *testing.T) {
	text, err := ioutil.ReadFile("./testdata/string.prod")
	if err != nil {
		t.Fatal(err)
	}
	s := parsec.NewScanner(text)
	root, _ := Y(s)
	scope := root.(common.Scope)
	forms := scope["_globalForms"].([]*common.Form)
	nterms := scope["_nonterminals"].(common.NTForms)
	scope = BuildContext(scope, uint64(time.Now().UnixNano()), "./bags")
	scope = scope.ApplyGlobalForms()
	if len(forms) > 0 {
		t.Fatalf("Expected empty forms for string.prod %v", forms)
	} else if len(nterms) != 1 {
		t.Fatalf("Expected single nterms for string.prod %v", nterms)
	}
	out := EvalForms("root", scope, nterms["s"]).(string)
	if out != "\"hello\"" {
		t.Fatalf("Unexpected string.prod out: %v", out)
	}
}

func TestTerm(t *testing.T) {
	text, err := ioutil.ReadFile("./testdata/term.prod")
	if err != nil {
		t.Fatal(err)
	}
	s := parsec.NewScanner(text)
	root, _ := Y(s)
	scope := root.(common.Scope)
	forms := scope["_globalForms"].([]*common.Form)
	nterms := scope["_nonterminals"].(common.NTForms)
	scope = BuildContext(scope, uint64(time.Now().UnixNano()), "./bags")
	scope = scope.ApplyGlobalForms()
	if len(forms) > 0 {
		t.Fatalf("Expected empty forms for term.prod %v", forms)
	} else if len(nterms) != 1 {
		t.Fatalf("Expected single nterms for term.prod %v", nterms)
	}
	out := EvalForms("root", scope, nterms["s"]).(string)
	if out != "\"" {
		t.Fatalf("Unexpected term.prod out: %v", out)
	}
}

func TestForm(t *testing.T) {
	text, err := ioutil.ReadFile("./testdata/form.prod")
	if err != nil {
		t.Fatal(err)
	}
	s := parsec.NewScanner(text)
	root, _ := Y(s)
	scope := root.(common.Scope)
	nterms := scope["_nonterminals"].(common.NTForms)
	forms := scope["_globalForms"].([]*common.Form)
	scope = BuildContext(scope, uint64(time.Now().UnixNano()), "./bags")
	scope = scope.ApplyGlobalForms()
	if len(forms) != 1 {
		t.Fatalf("Expected empty forms for form.prod %v", forms)
	} else if len(nterms) != 1 {
		t.Fatalf("Expected single nterms for form.prod %v", nterms)
	}
	out := EvalForms("root", scope, nterms["s"]).(string)
	if out != "\"10\"" {
		t.Fatalf("Unexpected term.prod out: %v", out)
	}
}

func TestNTerm(t *testing.T) {
	text, err := ioutil.ReadFile("./testdata/nterm.prod")
	if err != nil {
		t.Fatal(err)
	}
	s := parsec.NewScanner(text)
	root, _ := Y(s)
	scope := root.(common.Scope)
	nterms := scope["_nonterminals"].(common.NTForms)
	forms := scope["_globalForms"].([]*common.Form)
	scope = BuildContext(scope, uint64(time.Now().UnixNano()), "./bags")
	scope = scope.ApplyGlobalForms()
	if len(forms) != 1 {
		t.Fatalf("Expected empty forms for form.prod %v", forms)
	} else if len(nterms) != 2 {
		t.Fatalf("Expected single nterms for form.prod %v", nterms)
	}
	out := EvalForms("root", scope, nterms["s"]).(string)
	if out != "10 hello 20" {
		t.Fatalf("Unexpected term.prod out: %v", out)
	}
}

func TestOr(t *testing.T) {
	text, err := ioutil.ReadFile("./testdata/or.prod")
	if err != nil {
		t.Fatal(err)
	}
	s := parsec.NewScanner(text)
	root, _ := Y(s)
	scope := root.(common.Scope)
	nterms := scope["_nonterminals"].(common.NTForms)
	forms := scope["_globalForms"].([]*common.Form)
	scope = BuildContext(scope, uint64(time.Now().UnixNano()), "./bags")
	scope = scope.ApplyGlobalForms()
	if len(forms) != 1 {
		t.Fatalf("Expected empty forms for form.prod %v", forms)
	} else if len(nterms) != 1 {
		t.Fatalf("Expected single nterms for form.prod %v", nterms)
	}
	out := EvalForms("root", scope, nterms["s"]).(string)
	if out != "1020" && out != " hello " {
		t.Fatalf("Unexpected term.prod out: %v", out)
	}
}

func BenchmarkString(b *testing.B) {
	text, err := ioutil.ReadFile("./testdata/string.prod")
	if err != nil {
		b.Fatal(err)
	}
	s := parsec.NewScanner(text)
	root, _ := Y(s)
	scope := root.(common.Scope)
	nterms := scope["_nonterminals"].(common.NTForms)
	scope = BuildContext(scope, uint64(time.Now().UnixNano()), "./bags")
	scope = scope.ApplyGlobalForms()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		EvalForms("root", scope, nterms["s"])
	}
}

func BenchmarkTerm(b *testing.B) {
	text, err := ioutil.ReadFile("./testdata/term.prod")
	if err != nil {
		b.Fatal(err)
	}
	s := parsec.NewScanner(text)
	root, _ := Y(s)
	scope := root.(common.Scope)
	nterms := scope["_nonterminals"].(common.NTForms)
	scope = BuildContext(scope, uint64(time.Now().UnixNano()), "./bags")
	scope = scope.ApplyGlobalForms()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		EvalForms("root", scope, nterms["s"])
	}
}

func BenchmarkForm(b *testing.B) {
	text, err := ioutil.ReadFile("./testdata/form.prod")
	if err != nil {
		b.Fatal(err)
	}
	s := parsec.NewScanner(text)
	root, _ := Y(s)
	scope := root.(common.Scope)
	nterms := scope["_nonterminals"].(common.NTForms)
	scope = BuildContext(scope, uint64(time.Now().UnixNano()), "./bags")
	scope = scope.ApplyGlobalForms()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		EvalForms("root", scope, nterms["s"])
	}
}

func BenchmarkNTerm(b *testing.B) {
	text, err := ioutil.ReadFile("./testdata/nterm.prod")
	if err != nil {
		b.Fatal(err)
	}
	s := parsec.NewScanner(text)
	root, _ := Y(s)
	scope := root.(common.Scope)
	nterms := scope["_nonterminals"].(common.NTForms)
	scope = BuildContext(scope, uint64(time.Now().UnixNano()), "./bags")
	scope = scope.ApplyGlobalForms()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		EvalForms("root", scope, nterms["s"])
	}
}

func BenchmarkOr(b *testing.B) {
	text, err := ioutil.ReadFile("./testdata/or.prod")
	if err != nil {
		b.Fatal(err)
	}
	s := parsec.NewScanner(text)
	root, _ := Y(s)
	scope := root.(common.Scope)
	nterms := scope["_nonterminals"].(common.NTForms)
	scope = BuildContext(scope, uint64(time.Now().UnixNano()), "./bags")
	scope = scope.ApplyGlobalForms()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		EvalForms("root", scope, nterms["s"])
	}
}

func BenchmarkUsersProdY(b *testing.B) {
	text, err := ioutil.ReadFile("./prods/users.prod")
	if err != nil {
		b.Fatal(err)
	}
	s := parsec.NewScanner(text)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Y(s)
	}
}

func BenchmarkUsersProd(b *testing.B) {
	text, err := ioutil.ReadFile("./prods/users.prod")
	if err != nil {
		b.Fatal(err)
	}
	s := parsec.NewScanner(text)
	root, _ := Y(s)
	scope := root.(common.Scope)
	nterms := scope["_nonterminals"].(common.NTForms)
	scope = BuildContext(scope, uint64(time.Now().UnixNano()), "./bags")
	scope = scope.ApplyGlobalForms()

	b.ResetTimer()
	out := 0
	for i := 0; i < b.N; i++ {
		out += len(EvalForms("root", scope, nterms["s"]).(string))
	}
	b.SetBytes(int64(float64(out) / float64(b.N)))
}

func BenchmarkProjsProdY(b *testing.B) {
	text, err := ioutil.ReadFile("./prods/projects.prod")
	if err != nil {
		b.Fatal(err)
	}
	s := parsec.NewScanner(text)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Y(s)
	}
}

func BenchmarkProjsProd(b *testing.B) {
	text, err := ioutil.ReadFile("./prods/projects.prod")
	if err != nil {
		b.Fatal(err)
	}
	s := parsec.NewScanner(text)
	root, _ := Y(s)
	scope := root.(common.Scope)
	nterms := scope["_nonterminals"].(common.NTForms)
	scope = BuildContext(scope, uint64(time.Now().UnixNano()), "./bags")
	scope = scope.ApplyGlobalForms()

	b.ResetTimer()
	out := 0
	for i := 0; i < b.N; i++ {
		out += len(EvalForms("root", scope, nterms["s"]).(string))
	}
	b.SetBytes(int64(float64(out) / float64(b.N)))
}

func BenchmarkJSONProdY(b *testing.B) {
	text, err := ioutil.ReadFile("./prods/json.prod")
	if err != nil {
		b.Fatal(err)
	}
	s := parsec.NewScanner(text)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Y(s)
	}
}

func BenchmarkJSONProd(b *testing.B) {
	text, err := ioutil.ReadFile("./prods/projects.prod")
	if err != nil {
		b.Fatal(err)
	}
	s := parsec.NewScanner(text)
	root, _ := Y(s)
	scope := root.(common.Scope)
	nterms := scope["_nonterminals"].(common.NTForms)
	scope = BuildContext(scope, uint64(time.Now().UnixNano()), "./bags")
	scope = scope.ApplyGlobalForms()

	b.ResetTimer()
	out := 0
	for i := 0; i < b.N; i++ {
		out += len(EvalForms("root", scope, nterms["s"]).(string))
	}
	b.SetBytes(int64(float64(out) / float64(b.N)))
}
