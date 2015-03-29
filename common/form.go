//  Copyright (c) 2013 Couchbase, Inc.

package common

import "fmt"

var _ = fmt.Sprintf("dummy")

// FormFn function signature for all forms.
type FormFn func(scope Scope, args ...interface{}) interface{}

// NTForms is a map of non-terminal name to its rule definitions.
type NTForms map[string][]*Form

// Form defines the structure of any form.
type Form struct {
	Name     string
	Fn       FormFn
	Weight   float64
	Restrain float64
}

// NewForm will instantiate a form structure and return the same.
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

// SetWeight will set the weightage for this form. Typically
// this form will be rule-form.
func (form *Form) SetWeight(weight, restrain float64) {
	form.Weight, form.Restrain = weight, restrain
}

// SetDefaultWeight will set a default weightage for this form,
// if a weightage is not defined already. Typically this form
// will be rule-form.
func (form *Form) SetDefaultWeight(weight float64) {
	if form.Weight == 0.0 {
		form.Weight = weight
	}
}

// Eval will evaulate this form passing `scope` and `args`.
func (form *Form) Eval(scope Scope, args ...interface{}) interface{} {
	return form.Fn(scope, args...)
}

func (form *Form) String() string {
	return fmt.Sprintf("%s {%v, %v}", form.Name, form.Weight, form.Restrain)
}
