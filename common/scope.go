//  Copyright (c) 2013 Couchbase, Inc.

package common

import "fmt"
import "math/rand"
import "github.com/prataprc/goparsec"

var _ = fmt.Sprintf("dummy")

// Scope type
type Scope map[string]interface{}

// NewScopeFromRoot will create a new scope from root non-terminal.
func NewScopeFromRoot(ns []parsec.ParsecNode) Scope {
	globals := make(Scope)
	globals["_bagdir"] = ""
	globals["_prodfile"] = ""
	globals["_random"] = nil
	scope := Scope{
		"_globalForms":  ns[0].([]*Form),
		"_nonterminals": ns[1].(NTForms),
		"_weights":      make(map[string]float64),
		"_globals":      make(Scope),
	}
	return scope
}

// SetBagdir will set the bag-dir to be used by `bag` form.
func (scope Scope) SetBagdir(bagdir string) Scope {
	(scope["_globals"].(Scope))["_bagdir"] = bagdir
	return scope
}

// GetBagdir will return the current bagdir.
func (scope Scope) GetBagdir() string {
	return (scope["_globals"].(Scope))["_bagdir"].(string)
}

// SetProdfile will production filename.
func (scope Scope) SetProdfile(prodfile string) Scope {
	(scope["_globals"].(Scope))["_prodfile"] = prodfile
	return scope
}

// GetProdfile will return the current production filename.
func (scope Scope) GetProdfile() string {
	return (scope["_globals"].(Scope))["_prodfile"].(string)
}

// SetRandom will set *math/rand.Rand object.
func (scope Scope) SetRandom(rnd *rand.Rand) Scope {
	(scope["_globals"].(Scope))["_random"] = rnd
	return scope
}

// GetRandom will return current *math/rand.Rand object.
func (scope Scope) GetRandom() *rand.Rand {
	return (scope["_globals"].(Scope))["_random"].(*rand.Rand)
}

// SetWeight will set the weightage for form `name`.
func (scope Scope) SetWeight(name string, value float64) Scope {
	w := scope["_weights"].(map[string]float64)
	w[name] = value
	return scope
}

// GetWeight will return the weightage for form `name`.
func (scope Scope) GetWeight(name string) (value float64, ok bool) {
	value, ok = (scope["_weights"].(map[string]float64))[name]
	return value, ok
}

// SetNonTerminal will set the rule-forms for non-terminal `name`.
func (scope Scope) SetNonTerminal(name string, nt []*Form) Scope {
	(scope["_nonterminals"].(NTForms))[name] = nt
	return scope
}

// GetNonTerminal will return the rule-forms for non-terminal `name`.
func (scope Scope) GetNonTerminal(name string) (nt []*Form, ok bool) {
	ntforms := scope["_nonterminals"].(NTForms)
	nt, ok = ntforms[name]
	return nt, ok
}

// Set will set `name` to `value` in local scope or global
// scope (based on `g`).
func (scope Scope) Set(name string, value interface{}, g bool) Scope {
	if g {
		(scope["_globals"].(Scope))[name] = value
	} else {
		scope[name] = value
	}
	return scope
}

// Del will delete `name` from local scope or global scope (based on `g`).
func (scope Scope) Del(name string, value interface{}, g bool) Scope {
	if g {
		delete(scope["_globals"].(Scope), name)
	} else {
		delete(scope, name)
	}
	return scope
}

// Get will fetch `name` from local scope or global scope indicated by
// `g`.
func (scope Scope) Get(name string) (value interface{}, g, ok bool) {
	value, ok = scope[name]
	if ok {
		return value, false, true
	}
	value, ok = (scope["_globals"].(Scope))[name]
	if ok {
		return value, false, true
	}
	return nil, false, false
}

// GetString will get the variable name, convert it and return as string.
func (scope Scope) GetString(name string) (s string, g, ok bool) {
	value, g, ok := scope.Get(name)
	if ok {
		return value.(string), g, ok
	}
	return
}

// GetInt64 will get the variable name, convert it and return as int64.
func (scope Scope) GetInt64(name string) (i int64, g, ok bool) {
	value, g, ok := scope.Get(name)
	if ok {
		return value.(int64), g, ok
	}
	return
}

// Clone local scope and return back a new scope.
func (scope Scope) Clone() Scope {
	newS := make(Scope)
	for key, value := range scope {
		newS[key] = value
	}
	return newS
}

// ApplyGlobalForms on scope.
func (scope Scope) ApplyGlobalForms() Scope {
	newscope := scope.Clone()
	for _, form := range scope["_globalForms"].([]*Form) {
		form.Eval(newscope)
	}
	return newscope
}

// FormDuplicates will verify conflicts between user provided form-names
// and builtin form-names.
func (scope Scope) FormDuplicates(builtins map[string]*Form) []string {
	duplicates := make([]string, 0, 4)
	for name := range scope["_nonterminals"].(NTForms) {
		if _, ok := builtins[name]; ok {
			duplicates = append(duplicates, name)
		}
	}
	return duplicates
}
