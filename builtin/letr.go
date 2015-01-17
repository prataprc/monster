//  Copyright (c) 2013 Couchbase, Inc.

package builtin

import "fmt"

import "github.com/prataprc/monster/common"

var _ = fmt.Sprintf("dummy")

// Letr will define a single variables in local scope and return
// the same.
// args[0], args[2] ... args[N-1]  - variable name
// args[1], args[3] ... args[N] - variable value
func Letr(scope common.Scope, args ...interface{}) interface{} {
	scope.Set(args[0].(string), args[1], false /*global*/)
	return args[1]
}
