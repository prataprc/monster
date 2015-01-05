//  Copyright (c) 2013 Couchbase, Inc.

package builtin

import "fmt"
import "math/rand"
import "strconv"

import "github.com/prataprc/monster/common"

func Rangef(scope common.Scope, args ...interface{}) interface{} {
    var min, max float64
    var err error

    rnd := scope["_random"].(*rand.Rand)
    if len(args) == 2 {
        min, err = strconv.ParseFloat(args[0].(string), 64)
        if err != nil {
            panic(fmt.Errorf("parsing argument %v\n", args[0]))
        }
        max, err = strconv.ParseFloat(args[1].(string), 64)
        if err != nil {
            panic(fmt.Errorf("parsing argument %v\n", args[1]))
        }
        return rnd.Float64()*min + (max-min)

    } else if len(args) == 1 {
        max, err = strconv.ParseFloat(args[0].(string), 64)
        if err != nil {
            panic(fmt.Errorf("parsing argument %v\n", args[0]))
        }
        return rnd.Float64() * max
    }
    panic(fmt.Errorf("Atleast one argument expected for range-form\n"))
}
