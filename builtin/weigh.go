//  Copyright (c) 2013 Couchbase, Inc.

package builtin

import "fmt"

import "github.com/prataprc/monster/common"

var _ = fmt.Sprintf("dummy")

const DefaultRestrain = 0.0

// Weigh can be used as the first form in any rule, to
// define its choice preference for or grammar.
func Weigh(scope common.Scope, args ...interface{}) interface{} {
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

	if len(args) < 1 {
		panic(fmt.Errorf("insufficient arguments\n"))
	}
	weight := makefloat64(args[0])

	if len(args) > 1 {
		restrain := makefloat64(args[1])
		return []interface{}{weight, restrain}
	}
	return []interface{}{weight, DefaultRestrain}
}
