package query

import (
	"math"
)

const (
	NO_MORE   = int32(math.MaxInt32)
	NOT_READY = int32(-1)
)

type Query interface {
	advance(int32) int32
	Next() int32
	GetDocId() int32
	Score() float32
	String() string
}
