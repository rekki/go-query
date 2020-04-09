package query

import (
	"fmt"
	"sort"
	"strings"
)

type andQuery struct {
	queries []Query
	not     Query
	docId   int32
	leading Query
	boost   float32
}

// Creates AND NOT query
func AndNot(not Query, queries ...Query) *andQuery {
	return And(queries...).SetNot(not)
}

// Creates AND query
func And(queries ...Query) *andQuery {
	a := &andQuery{
		queries: queries,
		docId:   NOT_READY,
		boost:   1,
	}
	a.sortSubqueries()
	return a
}

func (q *andQuery) AddSubQuery(sub Query) *andQuery {
	q.queries = append(q.queries, sub)
	q.sortSubqueries()
	return q
}

func (q *andQuery) Cost() int {
	if len(q.queries) == 0 {
		return 0
	}
	return q.leading.Cost()
}

func (q *andQuery) sortSubqueries() {
	sort.Slice(q.queries, func(i, j int) bool {
		return q.queries[i].Cost() < q.queries[j].Cost()
	})
	if len(q.queries) > 0 {
		q.leading = q.queries[0]
	}
}

func (q *andQuery) GetDocId() int32 {
	return q.docId
}

func (q *andQuery) SetNot(not Query) *andQuery {
	q.not = not
	return q
}

func (q *andQuery) SetBoost(b float32) Query {
	q.boost = b
	return q
}

func (q *andQuery) PayloadDecode(p Payload) {
	p.Push()
	defer p.Pop()

	for _, s := range q.queries {
		s.PayloadDecode(p)
	}
}

func (q *andQuery) Score() float32 {
	score := float32(0)
	n := len(q.queries)

	for i := 0; i < n; i++ {
		score += q.queries[i].Score()
	}
	return score * q.boost
}

func (q *andQuery) nextAndedDoc(target int32) int32 {
	start := 1
	n := len(q.queries)
AGAIN:
	for {
		// initial iteration skips queries[0], because it is used in caller
		for i := start; i < n; i++ {
			subQuery := q.queries[i]
			subQueryDocId := subQuery.GetDocId()
			if subQueryDocId < target {
				subQueryDocId = subQuery.Advance(target)
			}

			if subQueryDocId == target {
				continue
			}

			target = q.leading.Advance(subQueryDocId)

			i = 0 //restart the loop from the first query
		}

		if q.not != nil && q.not.GetDocId() != NO_MORE && target != NO_MORE {
			if q.not.Advance(target) == target {
				target = q.leading.Advance(target + 1)
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

func (q *andQuery) Advance(target int32) int32 {
	if len(q.queries) == 0 {
		q.docId = NO_MORE
		return NO_MORE
	}

	return q.nextAndedDoc(q.leading.Advance(target))
}

func (q *andQuery) Next() int32 {
	if len(q.queries) == 0 {
		q.docId = NO_MORE
		return NO_MORE
	}

	return q.nextAndedDoc(q.leading.Next())
}
