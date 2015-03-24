//  Copyright (c) 2013 Couchbase, Inc.

package common

import "fmt"

var _ = fmt.Sprintf("dummy")

// EvalForms will evaluate the non-terminal form `name`
// by randomly picking one of the []*Form defined as its rules.
func EvalForms(name string, scope Scope, forms []*Form) interface{} {
	if len(forms) == 0 {
		return nil
	}
	rnd := scope.GetRandom()
	lookup, failed := make([]bool, len(forms)), 0
	for failed < len(forms) {
		f := rnd.Float64()
		for i, form := range forms {
			weight := currWeight(name, i, scope, form)
			if lookup[i] == false && f <= weight {
				lookup[i] = true
				scope = decWeight(name, i, scope, form) // weight restrainer
				val := form.Eval(scope.Clone())
				if val == nil {
					failed++
					continue
				}
				return val
			} else if weight <= 0.0 {
				failed++
			}
		}
	}
	return nil
}

func currWeight(name string, i int, scope Scope, form *Form) float64 {
	nm := fmt.Sprintf("%s%d", name, i)
	if w, ok := scope.GetWeight(nm); ok {
		return w
	}
	return form.Weight
}

func decWeight(name string, i int, scope Scope, form *Form) Scope {
	nm := fmt.Sprintf("%s%d", name, i)
	weight := form.Weight
	if w, ok := scope.GetWeight(nm); ok {
		weight = w
	}
	scope.SetWeight(nm, weight-form.Restrain)
	return scope
}
