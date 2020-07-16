package query

import "fmt"

type ConstantQuery struct {
	query Query
	boost float32
}

func Constant(boost float32, q Query) *ConstantQuery {
	return &ConstantQuery{
		query: q,
		boost: boost,
	}
}

func (q *ConstantQuery) Cost() int {
	return q.query.Cost()
}

func (q *ConstantQuery) GetDocId() int32 {
	return q.query.GetDocId()
}

func (q *ConstantQuery) Score() float32 {
	return q.boost
}

func (q *ConstantQuery) Advance(target int32) int32 {
	return q.query.Advance(target)
}

func (q *ConstantQuery) Next() int32 {
	return q.query.Next()
}

func (q *ConstantQuery) String() string {
	return fmt.Sprintf("{CONST(%f {%s})}", q.boost, q.query.String())
}

func (q *ConstantQuery) SetBoost(b float32) Query {
	q.boost = b
	return q
}

func (q *ConstantQuery) PayloadDecode(p Payload) {
	panic("unsupported")
}
