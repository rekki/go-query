package query

import (
	"fmt"
	"strings"
)

type andQuery struct {
	queries []Query
	not     Query
	docId   int32
}

// Creates AND NOT query
func AndNot(not Query, queries ...Query) *andQuery {
	return And(queries...).SetNot(not)
}

// Creates AND query
func And(queries ...Query) *andQuery {
	return &andQuery{
		queries: queries,
		docId:   NOT_READY,
	}
}

func (q *andQuery) GetDocId() int32 {
	return q.docId
}

func (q *andQuery) SetNot(not Query) *andQuery {
	q.not = not
	return q
}

func (q *andQuery) Score() float32 {
	return float32(len(q.queries))
}

func (q *andQuery) nextAndedDoc(target int32) int32 {
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

func (q *andQuery) String() string {
	out := []string{}
	for _, v := range q.queries {
		out = append(out, v.String())
	}
	s := strings.Join(out, " AND ")
	if q.not != nil {
		s = fmt.Sprintf("%s -(%s)", s, q.not.String())
	}
	return "{" + s + "}"
}

func (q *andQuery) advance(target int32) int32 {
	if len(q.queries) == 0 {
		q.docId = NO_MORE
		return NO_MORE
	}

	return q.nextAndedDoc(q.queries[0].advance(target))
}

func (q *andQuery) Next() int32 {
	if len(q.queries) == 0 {
		q.docId = NO_MORE
		return NO_MORE
	}

	// XXX: pick cheapest leading query
	return q.nextAndedDoc(q.queries[0].Next())
}
