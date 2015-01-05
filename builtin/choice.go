//  Copyright (c) 2013 Couchbase, Inc.

package builtin

import "fmt"
import "math/rand"

import "github.com/prataprc/monster/common"

var _ = fmt.Sprintf("dummy")

func Choice(scope common.Scope, args ...interface{}) interface{} {
    rnd := scope["_random"].(*rand.Rand)
    return args[rnd.Intn(len(args))]
}

