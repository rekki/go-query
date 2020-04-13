// Package query provides simple query dsl on top of sorted arrays of integers
package query

import (
	"math"
)

const (
	NO_MORE   = int32(math.MaxInt32)
	NOT_READY = int32(-1)
)

// Reuse/Concurrency:
// None of the queries are safe to be re-used.
// WARNING: the query *can not* be reused
// WARNING: the query it not thread safe
//
// Example Iteration:
//
//  q := Term([]int32{1,2,3})
//  for q.Next() != query.NO_MORE {
//  	did := q.GetDocId()
//  	score := q.Score()
//  	fmt.Printf("matching %d, score: %f\n", did, score)
//  }
type Query interface {
	Advance(int32) int32
	Next() int32
	GetDocId() int32
	Score() float32
	SetBoost(float32) Query
	Cost() int
	String() string

	PayloadDecode(p Payload)
}

type Payload interface {
	Push()
	Pop()
	Consume(int32, int, []byte)
	Score() float32
}
