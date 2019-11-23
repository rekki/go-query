//Example:

package main

import (
	"fmt"

	"github.com/jackdoe/go-query"
)

func main() {
	q := query.And(
		query.Or(
			query.Term("a", []int32{1, 2}),
			query.Term("b", []int32{3, 9})),
		query.AndNot(
			query.Term("c", []int32{4, 5}),
			query.Or(
				query.Term("d", []int32{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}),
				query.Term("e", []int32{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}),
			),
		),
	)
	for q.Next() != query.NO_MORE {
		did := q.GetDocId()
		score := q.Score()
		fmt.Printf("matching %d, score: %f\n", did, score)
	}
}
