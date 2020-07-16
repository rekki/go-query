package query

import "strings"

type DisMaxQuery struct {
	queries    []Query
	scores     []float32
	docId      int32
	tieBreaker float32
	boost      float32
}

// Creates DisMax query, for example if the query is:
//   DisMax(0.5, "name:amsterdam","name:university","name:free")
// lets say we have an index with following idf: amsterdam: 1.3, free: 0.2, university: 2.1
// the score is computed by:
//    max(score(amsterdam),score(university), score(free)) = 2.1 (university)
//    + score(free) * tiebreaker = 0.1
//    + score(amsterdam) * tiebreaker = 0.65
//    = 2.85
func DisMax(tieBreaker float32, queries ...Query) *DisMaxQuery {
	return &DisMaxQuery{
		queries:    queries,
		docId:      NOT_READY,
		tieBreaker: tieBreaker,
		scores:     make([]float32, len(queries)),
		boost:      1,
	}
}

func (q *DisMaxQuery) AddSubQuery(sub Query) *DisMaxQuery {
	q.queries = append(q.queries, sub)
	q.scores = make([]float32, len(q.queries))
	return q
}

func (q *DisMaxQuery) Cost() int {
	//XXX: optimistic, assume sets greatly overlap, which of course is not always true
	max := 0
	for _, sub := range q.queries {
		if max < sub.Cost() {
			max = sub.Cost()
		}
	}

	return max
}

func (q *DisMaxQuery) GetDocId() int32 {
	return q.docId
}

func (q *DisMaxQuery) Score() float32 {
	n := len(q.queries)
	sum := float32(0)
	max := float32(0)

	for i := 0; i < n; i++ {
		s := q.queries[i]
		if s.GetDocId() == q.docId {
			subQueryScore := s.Score()
			if subQueryScore > max {
				max = subQueryScore
			}
			sum += subQueryScore
		}
	}
	return (max + ((sum - max) * q.tieBreaker)) * q.boost
}

func (q *DisMaxQuery) Advance(target int32) int32 {
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

func (q *DisMaxQuery) Next() int32 {
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

func (q *DisMaxQuery) String() string {
	out := []string{}
	for _, v := range q.queries {
		out = append(out, v.String())
	}
	return "{" + strings.Join(out, " DisMax ") + "}"
}

func (q *DisMaxQuery) SetBoost(b float32) Query {
	q.boost = b
	return q
}

func (q *DisMaxQuery) PayloadDecode(p Payload) {
	panic("unsupported")
}
