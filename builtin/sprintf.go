//  Copyright (c) 2013 Couchbase, Inc.

package builtin

import "fmt"

import "github.com/prataprc/monster/common"

var _ = fmt.Sprintf("dummy")

func Sprintf(scope common.Scope, args ...interface{}) interface{} {
    if len(args) < 1 {
        panic(fmt.Errorf("insufficient argument to Sprintf"))
    }
    return fmt.Sprintf(args[0].(string), args[1:]...)
}
