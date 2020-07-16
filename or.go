package query

import "strings"

type OrQuery struct {
	queries []Query
	docId   int32
	boost   float32
}

// Creates OR query
func Or(queries ...Query) *OrQuery {
	return &OrQuery{
		queries: queries,
		docId:   NOT_READY,
		boost:   1,
	}
}

func (q *OrQuery) AddSubQuery(sub Query) *OrQuery {
	q.queries = append(q.queries, sub)
	return q
}

func (q *OrQuery) Cost() int {
	//XXX: optimistic, assume sets greatly overlap, which of course is not always true
	max := 0
	for _, sub := range q.queries {
		if max < sub.Cost() {
			max = sub.Cost()
		}
	}

	return max
}

func (q *OrQuery) GetDocId() int32 {
	return q.docId
}

func (q *OrQuery) PayloadDecode(p Payload) {
	p.Push()
	defer p.Pop()

	for _, s := range q.queries {
		if s.GetDocId() == q.docId {
			s.PayloadDecode(p)
		}
	}
}

func (q *OrQuery) Score() float32 {
	score := float32(0)
	n := len(q.queries)
	for i := 0; i < n; i++ {
		s := q.queries[i]
		if s.GetDocId() == q.docId {
			score += s.Score()
		}
	}
	return score * q.boost
}

func (q *OrQuery) Advance(target int32) int32 {
	newDoc := NO_MORE
	n := len(q.queries)
	for i := 0; i < n; i++ {
		subQuery := q.queries[i]
		curDoc := subQuery.GetDocId()
		if curDoc < target {
			curDoc = subQuery.Advance(target)
		}

		if curDoc < newDoc {
			newDoc = curDoc
		}
	}
	q.docId = newDoc
	return q.docId
}

func (q *OrQuery) Next() int32 {
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

func (q *OrQuery) String() string {
	out := []string{}
	for _, v := range q.queries {
		out = append(out, v.String())
	}
	return "{" + strings.Join(out, " OR ") + "}"
}

func (q *OrQuery) SetBoost(b float32) Query {
	q.boost = b
	return q
}
