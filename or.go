package query

import "strings"

type orQuery struct {
	queries []Query
	docId   int32
}

// Creates OR query
func Or(queries ...Query) *orQuery {
	return &orQuery{
		queries: queries,
		docId:   NOT_READY,
	}
}

func (q *orQuery) AddSubQuery(sub Query) {
	q.queries = append(q.queries, sub)
}
func (q *orQuery) GetDocId() int32 {
	return q.docId
}

func (q *orQuery) Score() float32 {
	score := 0
	n := len(q.queries)
	for i := 0; i < n; i++ {
		sub_query := q.queries[i]
		if sub_query.GetDocId() == q.docId {
			score++
		}
	}
	return float32(score)
}

func (q *orQuery) advance(target int32) int32 {
	new_doc := NO_MORE
	n := len(q.queries)
	for i := 0; i < n; i++ {
		sub_query := q.queries[i]
		cur_doc := sub_query.GetDocId()
		if cur_doc < target {
			cur_doc = sub_query.advance(target)
		}

		if cur_doc < new_doc {
			new_doc = cur_doc
		}
	}
	q.docId = new_doc
	return q.docId
}

func (q *orQuery) Next() int32 {
	new_doc := NO_MORE
	n := len(q.queries)
	for i := 0; i < n; i++ {
		sub_query := q.queries[i]
		cur_doc := sub_query.GetDocId()
		if cur_doc == q.docId {
			cur_doc = sub_query.Next()
		}

		if cur_doc < new_doc {
			new_doc = cur_doc
		}
	}
	q.docId = new_doc
	return new_doc
}

func (q *orQuery) String() string {
	out := []string{}
	for _, v := range q.queries {
		out = append(out, v.String())
	}
	return "{" + strings.Join(out, " OR ") + "}"
}
