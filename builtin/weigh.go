//  Copyright (c) 2013 Couchbase, Inc.

package builtin

import "fmt"

import "github.com/prataprc/monster/common"

var _ = fmt.Sprintf("dummy")

const DefaultRestrain = 0.0

// Weigh can be used as the first form in any rule, to
// define its choice preference for or grammar.
func Weigh(scope common.Scope, args ...interface{}) interface{} {
	if len(args) < 1 {
		panic(fmt.Errorf("insufficient arguments\n"))
	}
	weight := args[0].(float64)
	if len(args) > 1 {
		restrain := args[1].(float64)
		return []interface{}{weight, restrain}
	}
	return []interface{}{weight, DefaultRestrain}
}
