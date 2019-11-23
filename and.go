package query

import (
	"fmt"
	"strings"
)

type AndQuery struct {
	queries []Query
	not     Query
	docId   int32
}

func AndNot(not Query, queries ...Query) *AndQuery {
	return And(queries...).SetNot(not)
}

func And(queries ...Query) *AndQuery {
	return &AndQuery{
		queries: queries,
		docId:   NOT_READY,
	}
}

func (q *AndQuery) GetDocId() int32 {
	return q.docId
}

func (q *AndQuery) SetNot(not Query) *AndQuery {
	q.not = not
	return q
}

func (q *AndQuery) Score() float32 {
	return float32(len(q.queries))
}

func (q *AndQuery) nextAndedDoc(target int32) int32 {
	start := 1
	n := len(q.queries)
AGAIN:
	for {
		// initial iteration skips queries[0], because it is used in caller
		for i := start; i < n; i++ {
			sub_query := q.queries[i]
			if sub_query.GetDocId() < target {
				sub_query.advance(target)
			}

			if sub_query.GetDocId() == target {
				continue
			}

			target = q.queries[0].advance(sub_query.GetDocId())

			i = 0 //restart the loop from the first query
		}

		if q.not != nil && q.not.GetDocId() != NO_MORE && target != NO_MORE {
			if q.not.advance(target) == target {
				// the not query is matching, so we have to move on
				// advance everything, set the new target to the highest doc, and start again
				newTarget := target + 1
				for i := 0; i < n; i++ {
					current := q.queries[i].advance(newTarget)
					if current > newTarget {
						newTarget = current
					}
				}
				target = newTarget
				start = 0
				continue AGAIN
			}
		}

		q.docId = target
		return q.docId
	}
}

func (q *AndQuery) String() string {
	out := []string{}
	for _, v := range q.queries {
		out = append(out, v.String())
	}
	s := strings.Join(out, " AND ")
	if q.not != nil {
		s = fmt.Sprintf("%s[-(%s)]", s, q.not.String())
	}
	return s
}

func (q *AndQuery) advance(target int32) int32 {
	if len(q.queries) == 0 {
		q.docId = NO_MORE
		return NO_MORE
	}

	return q.nextAndedDoc(q.queries[0].advance(target))
}

func (q *AndQuery) Next() int32 {
	if len(q.queries) == 0 {
		q.docId = NO_MORE
		return NO_MORE
	}

	// XXX: pick cheapest leading query
	return q.nextAndedDoc(q.queries[0].Next())
}
