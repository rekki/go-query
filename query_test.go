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

func eq(t *testing.T, a, b []int32) {
	if len(a) != len(b) {
		log.Panicf("len(a) != len(b) ; len(a) = %d, len(b) = %d [%v %v]", len(a), len(b), a, b)
		t.FailNow()
	}

	for i, _ := range a {
		if a[i] != b[i] {
			t.Log("a[i] != b[i]")
			t.FailNow()
		}
	}
}

func BenchmarkNext1000(b *testing.B) {
	x := postingsList(1000)

	for n := 0; n < b.N; n++ {
		sum := int32(0)
		q := Term("", x)
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
			Term("x", x),
			Term("y", y),
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
			Term("x", x),
			Term("y", y),
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
			Term("x", x),
			Term("y", y),
		)

		for q.Next() != NO_MORE {
			sum += q.GetDocId()
		}
	}
}

func TestModify(t *testing.T) {
	eq(t, []int32{6, 7, 8, 10}, query(AndNot(
		Term("x", []int32{1, 2, 3, 9}),
		AndNot(
			Term("x", []int32{4, 5}),
			Term("x", []int32{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11}),
			Term("x", []int32{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}),
		),
	)))

	qq := And(
		Term("a", []int32{1, 2, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13}),
		Term("b", []int32{1, 3, 9}),
	)

	eq(t, []int32{1, 9}, query(qq))
	for i := 0; i < 100; i++ {
		k := rand.Intn(65000)
		a := postingsList(100 + k)
		b := postingsList(1000 + k)
		c := postingsList(10000 + k)
		d := postingsList(100000 + k)
		e := postingsList(1000000 + k)

		eq(t, a, query(Term("x", a)))
		eq(t, b, query(Term("x", b)))
		eq(t, c, query(Term("x", c)))
		eq(t, d, query(Term("x", d)))
		eq(t, e, query(Term("x", e)))

		eq(t, b, query(Or(
			Term("x", a),
			Term("x", b),
		)))

		eq(t, c, query(Or(
			Term("x", a),
			Term("x", b),
			Term("x", c),
		)))

		eq(t, e, query(Or(
			Term("x", a),
			Term("x", b),
			Term("x", c),
			Term("x", d),
			Term("x", e),
		)))

		eq(t, a, query(And(
			Term("x", a),
			Term("x", b),
			Term("x", c),
			Term("x", d),
			Term("x", e),
		)))

		eq(t, a, query(And(
			Term("x", a),
			Term("x", b),
			Term("x", c),
			Term("x", d),
			Term("x", e),
		)))
	}
	a := postingsList(100)
	b := postingsList(1000)
	c := postingsList(10000)
	d := postingsList(100000)
	e := postingsList(1000000)

	eq(t, []int32{4, 6, 7, 8, 10}, query(AndNot(
		Term("x", []int32{1, 2, 3, 9}),
		Or(
			Term("x", []int32{3, 4}),
			Term("x", []int32{1, 2, 3, 6, 7, 8, 9, 10}),
		),
	)))
	eq(t, []int32{6, 7, 8, 10}, query(AndNot(
		Term("x", []int32{1, 2, 3, 9}),
		AndNot(
			Term("x", []int32{4, 5}),
			Term("x", []int32{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}),
			Term("x", []int32{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}),
		),
	)))

	eq(t, []int32{6, 7, 8, 10}, query(AndNot(
		Or(
			Term("x", []int32{1, 2}),
			Term("x", []int32{3, 9})),
		AndNot(
			Term("x", []int32{4, 5}),
			Term("x", []int32{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}),
			Term("x", []int32{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}),
		),
	)))

	eq(t, []int32{}, query(AndNot(
		Term("x", []int32{1, 2, 3, 9}),
		Term("x", []int32{1, 2, 3, 9}),
	)))

	eq(t, []int32{}, query(AndNot(
		Term("x", []int32{1, 2, 3, 9}),
	)))

	eq(t, []int32{1, 2, 3, 9}, query(AndNot(
		Term("x", []int32{}),
		Term("x", []int32{1, 2, 3, 9}),
	)))

	eq(t, b, query(And(
		Or(
			Term("x", a),
			Term("x", b),
		),
		Term("x", b),
		Term("x", c),
		Term("x", d),
		Term("x", e),
	)))

	eq(t, c, query(And(
		Or(
			Term("x", a),
			Term("x", b),
			And(
				Term("x", c),
				Term("x", d),
			),
		),
		Term("x", d),
		Term("x", e),
	)))

	eq(t, []int32{1, 2, 3, 9}, query(And(
		Or(
			Term("x", []int32{1, 2}),
			Term("x", []int32{3, 9})),
		AndNot(
			Term("x", []int32{4, 5}),
			Or(
				Term("x", []int32{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}),
				Term("x", []int32{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}),
			),
		),
	)))
	q := And(
		Or(
			Term("a", []int32{1, 2}),
			Term("b", []int32{3, 9})),
		AndNot(
			Or(Term("c", []int32{4, 5}), Term("x", []int32{4, 100})),
			Or(
				Term("d", []int32{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}),
				Term("e", []int32{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}),
			),
		),
	)

	eq(t, []int32{1, 2, 3, 9}, query(q))
}
