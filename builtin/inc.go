//  Copyright (c) 2013 Couchbase, Inc.

package builtin

import "fmt"

import "github.com/prataprc/monster/common"

var _ = fmt.Sprintf("dummy")

func Inc(scope common.Scope, args ...interface{}) interface{} {
	name := args[0].(string)
	vali, ok := scope[name]
	if ok {
		scope[name] = vali.(int64) + 1
	}
	return ""
}
