package common

import "fmt"
import "github.com/prataprc/goparsec"

type Scope map[string]interface{}

type Weight float64

var _ = fmt.Sprintf("dummy")

func NewScopeFromRoot(ns []parsec.ParsecNode) Scope {
	scope := Scope{
		"_globalForms":  ns[0].([]*Form),
		"_nonterminals": ns[1].(NTForms),
	}
	return scope
}

func (scope Scope) Clone() Scope {
	newS := make(Scope)
	for key, value := range scope {
		newS[key] = value
	}
	return newS
}

func (scope Scope) Set(key string, value interface{}) Scope {
	scope[key] = value
	return scope
}

func (scope Scope) Del(key string) Scope {
	delete(scope, key)
	return scope
}

func (scope Scope) Get(key string) (val interface{}, ok bool) {
	val, ok = scope[key]
	return val, ok
}

func (scope Scope) GetNonTerminal(name string) (nt interface{}, ok bool) {
	ntforms := scope["_nonterminals"].(NTForms)
	nt, ok = ntforms[name]
	return nt, ok
}

func (scope Scope) ApplyGlobalForms() Scope {
	for _, form := range scope["_globalForms"].([]*Form) {
		form.Eval(scope)
	}
	return scope
}

func (scope Scope) FormDuplicates(builtins map[string]*Form) []string {
	duplicates := make([]string, 0, 4)
	for name, _ := range scope["_nonterminals"].(NTForms) {
		if _, ok := builtins[name]; ok {
			duplicates = append(duplicates, name)
		}
	}
	return duplicates
}
