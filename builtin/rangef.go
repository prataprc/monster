//  Copyright (c) 2013 Couchbase, Inc.

package builtin

import "fmt"

import "github.com/prataprc/monster/common"

// Rangef will randomly pick a value from args[0] to args[1]
// and return the same.
// args... are expected to be in float64
func Rangef(scope common.Scope, args ...interface{}) interface{} {
	rnd := scope.GetRandom()
	if len(args) == 2 {
		min, max := args[0].(float64), args[1].(float64)
		return rnd.Float64()*min + (max - min)

	} else if len(args) == 1 {
		max := args[0].(float64)
		return rnd.Float64() * max
	}
	panic(fmt.Errorf("atleast one argument expected for range-form\n"))
}
