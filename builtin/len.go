//  Copyright (c) 2013 Couchbase, Inc.

package builtin

import "fmt"

import "github.com/prataprc/monster/common"

var _ = fmt.Sprintf("dummy")

// Len will return the length of string in args[0]
func Len(scope common.Scope, args ...interface{}) interface{} {
	if s, ok := args[0].(string); ok {
		return int64(len(s))
	}
	return 0
}
