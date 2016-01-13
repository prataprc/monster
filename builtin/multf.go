//  Copyright (c) 2013 Couchbase, Inc.

package builtin

import "fmt"
import "strconv"

import "github.com/prataprc/monster/common"

func Multf(scope common.Scope, args ...interface{}) interface{} {
	floatof := func(arg interface{}) float64 {
		switch v := arg.(type) {
		case int64:
			return float64(v)
		case int32:
			return float64(v)
		case int:
			return float64(v)
		case float32:
			return float64(v)
		case float64:
			return float64(v)
		case string:
			if x, err := strconv.Atoi(v); err == nil {
				return float64(x)
			} else if y, err := strconv.ParseFloat(v, 64); err == nil {
				return y
			} else {
				panic(fmt.Errorf("unexpected str->number %v:%v", arg, err))
			}
		}
		panic(fmt.Errorf("unexpected numerical value: %T:%v", arg, arg))
	}
	if len(args) == 0 {
		return float64(0.0)
	}
	acc := floatof(args[0])
	for _, arg := range args[1:] {
		acc *= floatof(arg)
	}
	return acc
}
