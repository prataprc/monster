//  Copyright (c) 2013 Couchbase, Inc.

package builtin

import "fmt"

import "github.com/prataprc/monster/common"

var _ = fmt.Sprintf("dummy")

// Global will define a set of one or more variables in
// global scope.
// args[0], args[2] ... args[N-1]  - variable name
// args[1], args[3] ... args[N] - variable value
func Global(scope common.Scope, args ...interface{}) interface{} {
	for i := 0; i < len(args); i += 2 {
		scope.Set(args[i].(string), args[i+1], true /*global*/)
	}
	return ""
}
