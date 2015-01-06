package common

import "fmt"
import "math/rand"

var _ = fmt.Sprintf("dummy")

func EvalForms(name string, scope Scope, forms []*Form) interface{} {
	if len(forms) == 0 {
		return nil
	}
	rnd := scope["_random"].(*rand.Rand)
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
			}
		}
	}
	return nil
}

func currWeight(name string, i int, scope Scope, form *Form) float64 {
	nm := fmt.Sprintf("##form_%s%d", name, i)
	if w, ok := scope[nm]; ok {
		return w.(float64)
	}
	return form.Weight
}

func decWeight(name string, i int, scope Scope, form *Form) Scope {
	nm := fmt.Sprintf("##form_%s%d", name, i)
	weight := form.Weight
	if w, ok := scope[nm]; ok {
		weight = w.(float64)
	}
	scope[nm] = weight - form.Restrain
	return scope
}
