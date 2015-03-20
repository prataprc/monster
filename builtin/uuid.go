//  Copyright (c) 2013 Couchbase, Inc.

package builtin

import "fmt"
import "bytes"
import "strconv"
import "io"
import "encoding/binary"
import "math/rand"
import crypt "crypto/rand"

import "github.com/prataprc/monster/common"

var _ = fmt.Sprintf("dummy")

type UUID []byte

func init() {
	seed := newUUID().Uint64()
	rand.Seed(int64(seed))
}

// Uuid returns a unique value based on current nanosecond timestamp.
func Uuid(scope common.Scope, args ...interface{}) interface{} {
	return newUUID().Uint64()
}

func newUUID() UUID {
	uuid := make([]byte, 8)
	if n, err := io.ReadFull(crypt.Reader, uuid); err != nil {
		panic("crypt.Reader errored out")
	} else if n != len(uuid) {
		panic("crypt.Reader failed")
	}
	return UUID(uuid)
}

func (u UUID) Uint64() uint64 {
	return binary.LittleEndian.Uint64(([]byte)(u))
}

func (u UUID) Str() string {
	var buf bytes.Buffer
	for i := 0; i < len(u); i++ {
		if i > 0 {
			buf.WriteString(":")
		}
		buf.WriteString(strconv.FormatUint(uint64(u[i]), 16))
	}
	return buf.String()
}
