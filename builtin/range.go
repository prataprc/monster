//  Copyright (c) 2013 Couchbase, Inc.

package builtin

import "fmt"
import "math/rand"
import "strconv"

import "github.com/prataprc/monster/common"

func Range(scope common.Scope, args ...interface{}) interface{} {
    var min, max int
    var err error

    rnd := scope["_random"].(*rand.Rand)
    if len(args) == 2 {
        min, err = strconv.Atoi(args[0].(string))
        if err != nil {
            panic(fmt.Errorf("parsing argument %v\n", args[0]))
        }
        max, err = strconv.Atoi(args[1].(string))
        if err != nil {
            panic(fmt.Errorf("parsing argument %v\n", args[1]))
        }

    } else if len(args) == 1 {
        max, err = strconv.Atoi(args[0].(string))
        if err != nil {
            panic(fmt.Errorf("parsing argument %v\n", args[0]))
        }

    } else {
        panic(fmt.Errorf("Atleast one argument expected for range-form\n"))
    }
    return rnd.Intn(max-min) + min
}
