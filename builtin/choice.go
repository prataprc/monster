//  Copyright (c) 2013 Couchbase, Inc.

package builtin

import "fmt"

import "github.com/prataprc/monster/common"

var _ = fmt.Sprintf("dummy")

// Choice will randomly pick one of the passed argument
// and return back.
func Choice(scope common.Scope, args ...interface{}) interface{} {
	rnd := scope.GetRandom()
	return args[rnd.Intn(len(args))]
}
