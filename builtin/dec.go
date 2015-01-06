//  Copyright (c) 2013 Couchbase, Inc.

package builtin

import "fmt"

import "github.com/prataprc/monster/common"

var _ = fmt.Sprintf("dummy")

func Dec(scope common.Scope, args ...interface{}) interface{} {
	name := args[0].(string)
	vali, ok := scope[name]
	if ok {
		val, ok := vali.(int)
		if ok {
			scope[name] = val - 1
		}
	}
	return ""
}
