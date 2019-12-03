package query

import "strings"

type orQuery struct {
	queries []Query
	docId   int32
	boost   float32
}

// Creates OR query
func Or(queries ...Query) *orQuery {
	return &orQuery{
		queries: queries,
		docId:   NOT_READY,
		boost:   1,
	}
}

func (q *orQuery) AddSubQuery(sub Query) *orQuery {
	q.queries = append(q.queries, sub)
	return q
}

func (q *orQuery) cost() int {
	//XXX: optimistic, assume sets greatly overlap, which of course is not always true
	max := 0
	for _, sub := range q.queries {
		if max < sub.cost() {
			max = sub.cost()
		}
	}

	return max
}

func (q *orQuery) GetDocId() int32 {
	return q.docId
}

func (q *orQuery) Score() float32 {
	score := float32(0)
	n := len(q.queries)
	for i := 0; i < n; i++ {
		s := q.queries[i]
		if s.GetDocId() == q.docId {
			score += s.Score()
		}
	}
	return float32(score) * q.boost
}

func (q *orQuery) advance(target int32) int32 {
	newDoc := NO_MORE
	n := len(q.queries)
	for i := 0; i < n; i++ {
		subQuery := q.queries[i]
		curDoc := subQuery.GetDocId()
		if curDoc < target {
			curDoc = subQuery.advance(target)
		}

		if curDoc < newDoc {
			newDoc = curDoc
		}
	}
	q.docId = newDoc
	return q.docId
}

func (q *orQuery) Next() int32 {
	newDoc := NO_MORE
	n := len(q.queries)
	for i := 0; i < n; i++ {
		subQuery := q.queries[i]
		curDoc := subQuery.GetDocId()
		if curDoc == q.docId {
			curDoc = subQuery.Next()
		}

		if curDoc < newDoc {
			newDoc = curDoc
		}
	}
	q.docId = newDoc
	return newDoc
}

func (q *orQuery) String() string {
	out := []string{}
	for _, v := range q.queries {
		out = append(out, v.String())
	}
	return "{" + strings.Join(out, " OR ") + "}"
}

func (q *orQuery) SetBoost(b float32) Query {
	q.boost = b
	return q
}
