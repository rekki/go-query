package query

import "strings"

type OrQuery struct {
	queries []Query
	docId   int32
}

func Or(queries ...Query) *OrQuery {
	return &OrQuery{
		queries: queries,
		docId:   NOT_READY,
	}
}

func (q *OrQuery) AddSubQuery(sub Query) {
	q.queries = append(q.queries, sub)
}
func (q *OrQuery) GetDocId() int32 {
	return q.docId
}

func (q *OrQuery) Score() float32 {
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

func (q *OrQuery) advance(target int32) int32 {
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
func (q *OrQuery) String() string {
	out := []string{}
	for _, v := range q.queries {
		out = append(out, v.String())
	}
	return strings.Join(out, " OR ")
}

func (q *OrQuery) Next() int32 {
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
