package gocuke

import (
	"strconv"
	"sync/atomic"
)

var nextId uint64

func newId() string {
	id := atomic.AddUint64(&nextId, 1)
	return strconv.FormatUint(id, 10)
}
