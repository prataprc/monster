//  Copyright (c) 2013 Couchbase, Inc.

package builtin

import "fmt"
import "strconv"

import "github.com/prataprc/monster/common"

var _ = fmt.Sprintf("dummy")

func Weigh(scope common.Scope, args ...interface{}) interface{} {
    if len(args) < 1 {
        panic(fmt.Errorf("insufficient arguments\n"))
    }
    weight, err := strconv.ParseFloat(args[0].(string), 64)
    if err != nil {
        panic(fmt.Errorf("parsing weight argument %v\n", args[0]))
    }
    if len(args) > 1 {
        restrain, err := strconv.ParseFloat(args[1].(string), 64)
        if err != nil {
            panic(fmt.Errorf("parsing restrain argument %v\n", args[0]))
        }
        return []interface{}{weight, restrain}
    }
    return []interface{}{weight, 0.0}
}

