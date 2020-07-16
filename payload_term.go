package query

import (
	"fmt"
)

type PayloadTermQuery struct {
	term    *TermQuery
	payload []byte
}

func PayloadTerm(totalDocumentsInIndex int, t string, postings []int32, payload []byte) *PayloadTermQuery {
	term := Term(totalDocumentsInIndex, t, postings)
	return &PayloadTermQuery{
		term:    term,
		payload: payload,
	}
}

func (t *PayloadTermQuery) GetDocId() int32 {
	return t.term.docId
}

func (t *PayloadTermQuery) Cost() int {
	return len(t.term.postings) - t.term.cursor
}

func (t *PayloadTermQuery) String() string {
	return fmt.Sprintf("p_%s(%d)/%.2f", t.term.term, len(t.term.postings), t.term.idf)
}

func (t *PayloadTermQuery) Score() float32 {
	return t.term.Score()
}

func (t *PayloadTermQuery) Advance(target int32) int32 {
	return t.term.Advance(target)
}

func (t *PayloadTermQuery) Next() int32 {
	return t.term.Next()
}

func (t *PayloadTermQuery) SetBoost(b float32) Query {
	t.term.SetBoost(b)
	return t
}

func (t *PayloadTermQuery) PayloadDecode(p Payload) {
	p.Push()
	defer p.Pop()

	p.Consume(t.term.docId, t.term.cursor, t.payload)
}
