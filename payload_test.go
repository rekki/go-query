package query

import (
	"testing"
)

func TestPayload(t *testing.T) {
	q := And(
		Or(
			PayloadTerm(10, "a", []int32{1, 2}, []byte{10, 20}),
			PayloadTerm(10, "b", []int32{3, 9}, []byte{30, 90})),
		AndNot(
			nil,
			Or(
				PayloadTerm(10, "d", []int32{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}, []byte{10, 20, 30, 40, 50, 60, 70, 80, 90, 100}),
				PayloadTerm(10, "e", []int32{2, 4, 5, 8, 9, 10}, []byte{20, 40, 50, 80, 90, 100}),
			),
		),
	)

	p := &payload{}

	f := []float32{}
	for q.Next() != NO_MORE {
		q.PayloadDecode(p)
		f = append(f, p.Score())
	}
	eqF(t, f, []float32{20, 80, 140, 410})

}

type payload struct {
	stack int
	score int
}

func (p *payload) Pop() {
	p.stack--
}

func (p *payload) Push() {
	p.stack++
}

func (p *payload) Reset() {
	p.stack = 0
	p.score = 0
}

func (p *payload) Consume(_did int32, idx int, data []byte) {
	// one byte per document
	b := data[idx]
	p.score += int(b)
}
func (p *payload) Score() float32 {
	return float32(p.score)
}
