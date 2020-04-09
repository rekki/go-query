package query

import (
	"fmt"
)

type payloadTermQuery struct {
	term    *termQuery
	payload []byte
}

func PayloadTerm(totalDocumentsInIndex int, t string, postings []int32, payload []byte) *payloadTermQuery {
	term := Term(totalDocumentsInIndex, t, postings)
	return &payloadTermQuery{
		term:    term,
		payload: payload,
	}
}

func (t *payloadTermQuery) GetDocId() int32 {
	return t.term.docId
}

func (t *payloadTermQuery) Cost() int {
	return len(t.term.postings) - t.term.cursor
}

func (t *payloadTermQuery) String() string {
	return fmt.Sprintf("p_%s(%d)/%.2f", t.term.term, len(t.term.postings), t.term.idf)
}

func (t *payloadTermQuery) Score() float32 {
	return t.term.Score()
}

func (t *payloadTermQuery) Advance(target int32) int32 {
	return t.term.Advance(target)
}

func (t *payloadTermQuery) Next() int32 {
	return t.term.Next()
}

func (t *payloadTermQuery) SetBoost(b float32) Query {
	t.term.SetBoost(b)
	return t
}

func (t *payloadTermQuery) PayloadDecode(p Payload) {
	p.Push()
	defer p.Pop()

	p.Consume(t.term.docId, t.term.cursor, t.payload)
}
