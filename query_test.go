package query

import (
	"log"
	"math/rand"
	"testing"
)

func postingsList(n int) []int32 {
	list := make([]int32, n)
	for i := 0; i < n; i++ {
		list[i] = int32(i) * 3
	}
	return list
}

func query(query Query) []int32 {
	out := []int32{}
	for query.Next() != NO_MORE {
		out = append(out, query.GetDocId())
	}
	return out
}

func queryScores(query Query) []float32 {
	out := []float32{}
	for query.Next() != NO_MORE {
		out = append(out, query.Score())
	}
	return out
}

func eq(t *testing.T, a, b []int32) {
	if len(a) != len(b) {
		log.Panicf("len(a) != len(b) ; len(a) = %d, len(b) = %d [%v %v]", len(a), len(b), a, b)
		t.FailNow()
	}

	for i := range a {
		if a[i] != b[i] {
			t.Logf("a[i] != b[i]; %v != %v", a, b)
			t.FailNow()
		}
	}
}

func eqF(t *testing.T, a, b []float32) {
	if len(a) != len(b) {
		log.Panicf("len(a) != len(b) ; len(a) = %d, len(b) = %d [%v %v]", len(a), len(b), a, b)
		t.FailNow()
	}

	for i := range a {
		if a[i] != b[i] {
			log.Panicf("a[i] != b[i]; %v != %v", a, b)
		}
	}
}

func BenchmarkNext1000(b *testing.B) {
	x := postingsList(1000)

	for n := 0; n < b.N; n++ {
		sum := int32(0)
		q := Term(10, "", x)
		for q.Next() != NO_MORE {
			sum += q.GetDocId()
		}
	}
}

func BenchmarkOr1000(b *testing.B) {
	x := postingsList(1000)
	y := postingsList(1000)

	for n := 0; n < b.N; n++ {
		sum := int32(0)
		q := Or(
			Term(10, "x", x),
			Term(10, "y", y),
		)

		for q.Next() != NO_MORE {
			sum += q.GetDocId()
		}
	}
}

func BenchmarkAnd1000(b *testing.B) {
	x := postingsList(1000000)
	y := postingsList(1000)

	for n := 0; n < b.N; n++ {
		sum := int32(0)
		q := And(
			Term(10, "x", x),
			Term(10, "y", y),
		)

		for q.Next() != NO_MORE {
			sum += q.GetDocId()
		}
	}
}

func BenchmarkAnd1000000(b *testing.B) {
	x := postingsList(1000000)
	y := postingsList(10000)

	for n := 0; n < b.N; n++ {
		sum := int32(0)
		q := And(
			Term(10, "x", x),
			Term(10, "y", y),
		)

		for q.Next() != NO_MORE {
			sum += q.GetDocId()
		}
	}
}

func TestModify(t *testing.T) {
	eqF(t,
		queryScores(
			Or(Term(10, "x", []int32{1, 2, 3})),
		),
		queryScores(
			DisMax(1, Term(10, "x", []int32{1, 2, 3})),
		),
	)

	eqF(t, []float32{
		computeIDF(10, 2) + 0.1*computeIDF(10, 3) + 0.1*computeIDF(10, 4),
		computeIDF(10, 3) + 0.1*computeIDF(10, 4),
		computeIDF(10, 4),
		computeIDF(10, 2) + 0.1*computeIDF(10, 3) + 0.1*computeIDF(10, 4),
	}, queryScores(
		DisMax(0.1, Term(10, "x", []int32{1, 2, 3, 4}), Term(10, "x", []int32{1, 2, 4}), Term(10, "x", []int32{1, 4})),
	))

	qu := Term(10, "x", []int32{1, 2, 3, 4})
	qu.SetBoost(0)
	eqF(t, []float32{0, 0, 0, 0}, queryScores(
		DisMax(0.1, qu),
	))

	qu = Term(10, "x", []int32{1, 2, 3, 4})
	qu.SetBoost(1)
	eqF(t, []float32{1.2527629, 1.2527629, 1.2527629, 1.2527629}, queryScores(
		DisMax(0.1, qu),
	))

	eq(t, []int32{0, 10}, query(AndNot(
		Or(Term(10, "x", []int32{1})),
		Term(10, "x", []int32{0, 1, 7, 10}),
		Term(10, "x", []int32{0, 1, 6, 10}),
	)))

	eq(t, []int32{0, 2}, query(AndNot(
		Or(Term(10, "x", []int32{1}), Term(10, "x", []int32{})),
		Term(10, "x", []int32{0, 1, 2}),
	)))

	eq(t, []int32{6, 7, 8, 10}, query(AndNot(
		Term(10, "x", []int32{1, 2, 3, 9}),
		AndNot(
			Term(10, "x", []int32{4, 5}),
			Term(10, "x", []int32{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11}),
			Term(10, "x", []int32{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}),
		),
	)))

	qq := And(
		Term(10, "a", []int32{1, 2, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13}),
		Term(10, "b", []int32{1, 3, 9}),
	)

	eq(t, []int32{1, 9}, query(qq))
	for i := 0; i < 100; i++ {
		k := rand.Intn(65000)
		a := postingsList(100 + k)
		b := postingsList(1000 + k)
		c := postingsList(10000 + k)
		d := postingsList(100000 + k)
		e := postingsList(1000000 + k)

		eq(t, a, query(Term(10, "x", a)))
		eq(t, b, query(Term(10, "x", b)))
		eq(t, c, query(Term(10, "x", c)))
		eq(t, d, query(Term(10, "x", d)))
		eq(t, e, query(Term(10, "x", e)))

		eq(t, b, query(Or(
			Term(10, "x", a),
			Term(10, "x", b),
		)))

		eq(t, c, query(Or(
			Term(10, "x", a),
			Term(10, "x", b),
			Term(10, "x", c),
		)))

		eq(t, e, query(Or(
			Term(10, "x", a),
			Term(10, "x", b),
			Term(10, "x", c),
			Term(10, "x", d),
			Term(10, "x", e),
		)))

		eq(t, a, query(And(
			Term(10, "x", a),
			Term(10, "x", b),
			Term(10, "x", c),
			Term(10, "x", d),
			Term(10, "x", e),
		)))

		eq(t, a, query(And(
			Term(10, "x", a),
			Term(10, "x", b),
			Term(10, "x", c),
			Term(10, "x", d),
			Term(10, "x", e),
		)))
	}
	a := postingsList(100)
	b := postingsList(1000)
	c := postingsList(10000)
	d := postingsList(100000)
	e := postingsList(1000000)

	eq(t, []int32{4, 6, 7, 8, 10}, query(AndNot(
		Term(10, "x", []int32{1, 2, 3, 9}),
		Or(
			Term(10, "x", []int32{3, 4}),
			Term(10, "x", []int32{1, 2, 3, 6, 7, 8, 9, 10}),
		),
	)))
	eq(t, []int32{6, 7, 8, 10}, query(AndNot(
		Term(10, "x", []int32{1, 2, 3, 9}),
		AndNot(
			Term(10, "x", []int32{4, 5}),
			Term(10, "x", []int32{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}),
			Term(10, "x", []int32{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}),
		),
	)))

	eq(t, []int32{6, 7, 8, 10}, query(AndNot(
		Or(
			Term(10, "x", []int32{1, 2}),
			Term(10, "x", []int32{3, 9})),
		AndNot(
			Term(10, "x", []int32{4, 5}),
			Term(10, "x", []int32{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}),
			Term(10, "x", []int32{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}),
		),
	)))

	eq(t, []int32{}, query(AndNot(
		Term(10, "x", []int32{1, 2, 3, 9}),
		Term(10, "x", []int32{1, 2, 3, 9}),
	)))

	eq(t, []int32{}, query(AndNot(
		Term(10, "x", []int32{1, 2, 3, 9}),
	)))

	eq(t, []int32{1, 2, 3, 9}, query(AndNot(
		Term(10, "x", []int32{}),
		Term(10, "x", []int32{1, 2, 3, 9}),
	)))

	eq(t, b, query(And(
		Or(
			Term(10, "x", a),
			Term(10, "x", b),
		),
		Term(10, "x", b),
		Term(10, "x", c),
		Term(10, "x", d),
		Term(10, "x", e),
	)))

	eq(t, c, query(And(
		Or(
			Term(10, "x", a),
			Term(10, "x", b),
			And(
				Term(10, "x", c),
				Term(10, "x", d),
			),
		),
		Term(10, "x", d),
		Term(10, "x", e),
	)))

	eq(t, []int32{1, 2, 3, 9}, query(And(
		Or(
			Term(10, "x", []int32{1, 2}),
			Term(10, "x", []int32{3, 9})),
		AndNot(
			Term(10, "x", []int32{4, 5}),
			Or(
				Term(10, "x", []int32{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}),
				Term(10, "x", []int32{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}),
			),
		),
	)))
	q := And(
		Or(
			Term(10, "a", []int32{1, 2}),
			Term(10, "b", []int32{3, 9})),
		AndNot(
			Or(Term(10, "c", []int32{4, 5}), Term(10, "x", []int32{4, 100})),
			Or(
				Term(10, "d", []int32{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}),
				Term(10, "e", []int32{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}),
			),
		),
	)

	eq(t, []int32{1, 2, 3, 9}, query(q))
}
