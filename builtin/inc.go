//  Copyright (c) 2013 Couchbase, Inc.

package builtin

import "fmt"
import "strconv"

import "github.com/prataprc/monster/common"

var _ = fmt.Sprintf("dummy")

func Inc(scope common.Scope, args ...interface{}) interface{} {
    name := args[0].(string)
    vali, ok := scope[name]
    if ok {
        if val, ok := vali.(int); ok {
            scope[name] = val+1
        } else if val, ok := vali.(string); ok {
            v, _ := strconv.Atoi(val)
            scope[name] = v+1
        }
    }
    return ""
}
