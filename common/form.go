//  Copyright (c) 2013 Couchbase, Inc.

package common

import "fmt"

var _ = fmt.Sprintf("dummy")

type NTForms map[string][]*Form

type FormFn func(scope Scope, args ...interface{}) interface{}

type Form struct {
	Name     string
	Fn       FormFn
	Weight   float64
	Restrain float64
}

func NewForm(name string, fun interface{}) *Form {
	form := &Form{Name: name}
	switch fn := fun.(type) {
	case func(Scope, ...interface{}) interface{}:
		form.Fn = FormFn(fn)
	case FormFn:
		form.Fn = FormFn(fn)
	}
	return form
}

func (form *Form) SetWeight(weight, restrain float64) {
	form.Weight, form.Restrain = weight, restrain
}

func (form *Form) SetDefaultWeight(weight float64) {
	if form.Weight == 0.0 {
		form.Weight = weight
	}
}

func (form *Form) Eval(scope Scope, args ...interface{}) interface{} {
	return form.Fn(scope, args...)
}
