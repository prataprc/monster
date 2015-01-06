//  Copyright (c) 2013 Couchbase, Inc.

package builtin

import "time"
import "fmt"

import "github.com/prataprc/monster/common"

var _ = fmt.Sprintf("dummy")

func Uuid(scope common.Scope, args ...interface{}) interface{} {
	uuid := time.Now().UnixNano()
	return uuid
}
