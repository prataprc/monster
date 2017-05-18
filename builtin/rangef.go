//  Copyright (c) 2013 Couchbase, Inc.

package builtin

import "fmt"

import "github.com/prataprc/monster/common"

// Rangef will randomly pick a value from args[0] to args[1]
// and return the same.
// args... are expected to be in float64
func Rangef(scope common.Scope, args ...interface{}) interface{} {
	makefloat64 := func(arg interface{}) float64 {
		val, ok := arg.(float64)
		if ok == false {
			intw, ok := arg.(int64)
			if ok == false {
				panic("expected int64, or float64")
			}
			val = float64(intw)
		}
		return val
	}

	rnd := scope.GetRandom()
	if len(args) == 2 {
		min, max := makefloat64(args[0]), makefloat64(args[1])
		f := (rnd.Float64() * (max - min)) + min
		return f

	} else if len(args) == 1 {
		max := makefloat64(args[0])
		return rnd.Float64() * max
	}
	panic(fmt.Errorf("atleast one argument expected for range-form\n"))
}
