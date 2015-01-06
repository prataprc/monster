//  Copyright (c) 2013 Couchbase, Inc.

package builtin

import "fmt"
import "time"
import "math/rand"

import "github.com/prataprc/monster/common"

func Ranget(scope common.Scope, args ...interface{}) interface{} {
	rnd := scope["_random"].(*rand.Rand)
	start, err := time.Parse(time.RFC3339, args[0].(string))
	if err != nil {
		panic(fmt.Errorf("parsing first argument %v: %v\n", args[0], err))
	}
	end, err := time.Parse(time.RFC3339, args[1].(string))
	if err != nil {
		panic(fmt.Errorf("parsing second argument %v: %v\n", args[0], err))
	}
	t := start.Add(time.Duration(rnd.Int63n(int64(end.Sub(start)))))
	return t.Format(time.RFC3339)
}
