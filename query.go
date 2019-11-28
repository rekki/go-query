// Package query provides simple query dsl on top of sorted arrays of integers.
// Usually when you have inverted index you endup having something like:
//  data := []*Document{}
//  index := map[string][]int32{}
//  for docId, d := range document {
//  	for _, token := range tokenize(normalize(d.Name)) {
//          index[token] = append(index[token],docId)
//      }
//  }
//
// then from documents like {hello world},{hello},{new york},{new world} you get inverted index in the form of:
//   {
//      "hello": [0,1],
//      "world": [0,3],
//      "new": [2,3],
//      "york": [2]
//   }
// anyway, if you want to read more on those check out the IR-book
//
// This package helps you query indexes of this form, in fairly efficient way, keep in mind it expects the []int32 slices to be _sorted_.
// Example:
//  package main
//
//  import (
//  	"fmt"
//
//  	"github.com/jackdoe/go-query"
//  )
//
//  func main() {
//  	q := query.And(
//  		query.Or(
//  			query.Term("a", []int32{1, 2}),
//  			query.Term("b", []int32{3, 9})),
//  		query.AndNot(
//  			query.Or(query.Term("c", []int32{4, 5}), query.Term("c", []int32{4, 100})),
//  			query.Or(
//  				query.Term("d", []int32{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}),
//  				query.Term("e", []int32{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}),
//  			),
//  		),
//  	)
//
//  	// q.String() is {{a OR b} AND {{d OR e} -({c OR x})}}
//
//  	for q.Next() != query.NO_MORE {
//  		did := q.GetDocId()
//  		score := q.Score()
//  		fmt.Printf("matching %d, score: %f\n", did, score)
//  	}
//  }
//
// will print:
//  matching 1, score: 2.000000
//  matching 2, score: 2.000000
//  matching 3, score: 2.000000
//  matching 9, score: 2.000000
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
//
// Example Iteration:
//      q := Term([]int32{1,2,3})
//  	for q.Next() != query.NO_MORE {
//  		did := q.GetDocId()
//  		score := q.Score()
//  		fmt.Printf("matching %d, score: %f\n", did, score)
//  	}
type Query interface {
	advance(int32) int32
	Next() int32
	GetDocId() int32
	Score() float32
	cost() int
	String() string
}
