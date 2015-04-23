//  Copyright (c) 2013 Couchbase, Inc.

package common

import "fmt"
import "strconv"
import "math"

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
	//fmt.Println()
	for failed < len(forms) {
		f := rnd.Float64() * maxWeigh
		for i, form := range forms {
			weight := currWeight(name, i, scope, form)
			//fmt.Println(name, i, maxWeigh, f, weight, f <= weight, weight <= 0.0)
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
	if weight -= form.Restrain; weight <= 0 {
		weight = 0
	}
	scope.SetWeight(nm, weight)
	return scope
}

func maxWeighsOfForms(name string, scope Scope, forms []*Form) float64 {
	max := math.Inf(-1)
	for i, form := range forms {
		weight := currWeight(name, i, scope, form)
		//fmt.Println(i, weight, x)
		if weight > max {
			max = weight
		}
	}
	//fmt.Println(max)
	return max
}
