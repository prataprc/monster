//  Copyright (c) 2013 Couchbase, Inc.

package builtin

import "time"

import "github.com/prataprc/monster/common"

func UnixNano(scope common.Scope, args ...interface{}) interface{} {
	return time.Now().UnixNano()
}
