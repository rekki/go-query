package query

import "fmt"

type constantQuery struct {
	query Query
	boost float32
}

func Constant(boost float32, q Query) *constantQuery {
	return &constantQuery{
		query: q,
		boost: boost,
	}
}

func (q *constantQuery) cost() int {
	return q.query.cost()
}

func (q *constantQuery) GetDocId() int32 {
	return q.query.GetDocId()
}

func (q *constantQuery) Score() float32 {
	return q.boost
}

func (q *constantQuery) advance(target int32) int32 {
	return q.query.advance(target)
}

func (q *constantQuery) Next() int32 {
	return q.query.Next()
}

func (q *constantQuery) String() string {
	return fmt.Sprintf("{CONST(%f {%s})}", q.boost, q.query.String())
}

func (q *constantQuery) SetBoost(b float32) Query {
	q.boost = b
	return q
}
