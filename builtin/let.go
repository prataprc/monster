//  Copyright (c) 2013 Couchbase, Inc.

package builtin

import "fmt"

import "github.com/prataprc/monster/common"

var _ = fmt.Sprintf("dummy")

func Let(scope common.Scope, args ...interface{}) interface{} {
	for i := 0; i < len(args); i += 2 {
		name := args[i].(string)
		scope[name] = args[i+1]
	}
	return ""
}
