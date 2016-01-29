//  Copyright (c) 2013 Couchbase, Inc.

package builtin

import "fmt"
import "github.com/prataprc/monster/common"

// Chr returns a character as string.
func Chr(scope common.Scope, args ...interface{}) interface{} {
	if len(args) == 0 {
		panic("insufficient args to (chr)")
	}
	return fmt.Sprintf("%c", args[0])
}
