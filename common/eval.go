//  Copyright (c) 2013 Couchbase, Inc.

package common

import "fmt"
import "strconv"

var _ = fmt.Sprintf("dummy")

// EvalForms will evaluate the non-terminal form `name`
// by randomly picking one of the []*Form defined by its
// rules.
func EvalForms(name string, scope Scope, forms []*Form) interface{} {
	if len(forms) == 0 {
		return nil
	}
	rnd := scope.GetRandom()
	lookup, failed := make([]bool, len(forms)), 0
	maxWeigh := maxWeighsOfForms(name, scope, forms)
	for failed < len(forms) {
		f := float64(rnd.Int31n(maxWeigh)) / float64(0x7FFFFFFF)
		for i, form := range forms {
			weight := currWeight(name, i, scope, form)
			if weight <= 0.0 {
				failed++

			} else if lookup[i] == false && f <= weight {
				lookup[i] = true
				scope = decWeight(name, i, scope, form) // weight restrainer
				val := form.Eval(scope.Clone())
				if val == nil {
					failed++
					continue
				}
				return val
			}
		}
	}
	return nil
}

func currWeight(name string, i int, scope Scope, form *Form) float64 {
	nm := name + strconv.Itoa(i)
	if w, ok := scope.GetWeight(nm); ok {
		return w
	}
	return form.Weight
}

func decWeight(name string, i int, scope Scope, form *Form) Scope {
	nm := name + strconv.Itoa(i)
	weight := form.Weight
	if w, ok := scope.GetWeight(nm); ok {
		weight = w
	}
	scope.SetWeight(nm, weight-form.Restrain)
	return scope
}

func maxWeighsOfForms(name string, scope Scope, forms []*Form) int32 {
	max := int32(0)
	for i, form := range forms {
		weight := currWeight(name, i, scope, form)
		x := int32(weight * float64(0x7FFFFFFF))
		if x > max {
			max = x
		}
	}
	return max
}
