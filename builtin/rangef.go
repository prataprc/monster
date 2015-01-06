//  Copyright (c) 2013 Couchbase, Inc.

package builtin

import "fmt"
import "math/rand"

import "github.com/prataprc/monster/common"

func Rangef(scope common.Scope, args ...interface{}) interface{} {
	rnd := scope["_random"].(*rand.Rand)
	if len(args) == 2 {
		min, max := args[0].(float64), args[1].(float64)
		return rnd.Float64()*min + (max - min)

	} else if len(args) == 1 {
		max := args[0].(float64)
		return rnd.Float64() * max
	}
	panic(fmt.Errorf("Atleast one argument expected for range-form\n"))
}
