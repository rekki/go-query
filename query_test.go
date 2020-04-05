package query

import (
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"math/rand"
	"os"
	"path"
	"sort"
	"strings"
	"testing"
)

type IntSlice []int32

func (p IntSlice) Len() int           { return len(p) }
func (p IntSlice) Less(i, j int) bool { return p[i] < p[j] }
func (p IntSlice) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

func postingsList(n int, into ...[]int32) []int32 {
	list := []int32{}
	for _, in := range into {
		list = append(list, in...)
	}

	max := int32(math.MaxInt32) - 1
	for i := 0; i < n; i++ {
		list = append(list, rand.Int31n(max))
	}
	sort.Sort(IntSlice(list))
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

func TestBoost(t *testing.T) {
	if queryScores(
		Term(6, "x", []int32{1, 2, 3}).SetBoost(100),
	)[0] < 100 {
		t.Fatal("no boost")
	}

	if queryScores(
		CreateFileTerm(6, "x", []int32{1, 2, 3}).SetBoost(100),
	)[0] < 100 {
		t.Fatal("no boost")
	}

	if queryScores(
		Or(Term(6, "x", []int32{1, 2, 3})).SetBoost(100),
	)[0] < 100 {
		t.Fatal("no boost")
	}

	if queryScores(
		And(Term(6, "x", []int32{1, 2, 3})).SetBoost(100),
	)[0] < 100 {
		t.Fatal("no boost")
	}

	if queryScores(
		DisMax(1, Term(6, "x", []int32{1, 2, 3})).SetBoost(100),
	)[0] < 100 {
		t.Fatal("no boost")
	}
	if queryScores(
		Constant(0, And(Term(6, "x", []int32{1, 2, 3}), Constant(1, Term(6, "x", []int32{1, 2, 3}).SetBoost(200)))).SetBoost(2),
	)[0] != 2 {
		t.Fatal("no boost")
	}

	if queryScores(
		Constant(0, And(Term(6, "x", []int32{1, 2, 3}), Constant(1, Term(6, "x", []int32{1, 2, 3}).SetBoost(200)))),
	)[0] != 0 {
		t.Fatal("no boost")
	}

	if queryScores(
		Constant(0, Or(Term(6, "x", []int32{1, 2, 3}), Constant(1, Term(6, "x", []int32{1, 2, 3}).SetBoost(200)))).SetBoost(2),
	)[0] != 2 {
		t.Fatal("no boost")
	}

}
func TestTermAdvanceNotMatch(t *testing.T) {
	perChunkA := make([]int32, TERM_CHUNK_SIZE)
	perChunkB := make([]int32, TERM_CHUNK_SIZE)
	for i := 0; i < len(perChunkA); i++ {
		perChunkA[i] = int32(1 + i)
		perChunkA[i] = int32(TERM_CHUNK_SIZE + 10 + i)
	}
	perChunkA = append(perChunkA, 10000000, 10000002)
	perChunkB = append(perChunkB, 10000000, 10000003)

	eq(t, []int32{10000000}, query(And(
		Term(10, "x", perChunkA),
		Term(10, "x", perChunkB),
	)))
}

func TestEmpty(t *testing.T) {
	eq(t, []int32{}, query(And(And(), Term(10, "x", []int32{1, 2, 3}), And())))
	eq(t, []int32{}, query(Or()))
	eq(t, []int32{}, query(DisMax(0)))
}

func TestAddSubQuery(t *testing.T) {
	eq(t, []int32{2, 3}, query(And(Term(10, "x", []int32{1, 2, 3})).AddSubQuery(Term(10, "x", []int32{2, 3}))))
	eq(t, []int32{1, 2, 3}, query(Or(Term(10, "x", []int32{1, 2, 3})).AddSubQuery(Term(10, "x", []int32{2, 3}))))
	eq(t, []int32{1, 2, 3}, query(DisMax(1, Term(10, "x", []int32{1, 2, 3})).AddSubQuery(Term(10, "x", []int32{2, 3}))))
}

func CreateFileTerm(n int, _t string, postings []int32) Query {
	dir, err := ioutil.TempDir("", "tt")
	if err != nil {
		panic(err)
	}
	defer os.RemoveAll(dir)
	fn := path.Join(dir, fmt.Sprintf("t_%d", rand.Int()))

	err = AppendFileNameTerm(fn, postings)
	if err != nil {
		panic(err)
	}

	return NewFileTerm(n, fn)

}
func TestStrings(t *testing.T) {
	s := Constant(1,
		DisMax(1,
			AndNot(Term(10, "y", []int32{1}), Term(10, "x", []int32{1, 2, 3}), Term(10, "x", []int32{1, 2, 3})),
			Or(
				AndNot(Term(10, "y", []int32{1}), Term(10, "x", []int32{1, 2, 3}), Term(10, "x", []int32{1, 2, 3})),
				AndNot(Term(10, "y", []int32{1}), Term(10, "x", []int32{1, 2, 3}), Term(10, "x", []int32{1, 2, 3})),
				AndNot(CreateFileTerm(10, "y", []int32{1}), Term(10, "x", []int32{1, 2, 3}), Term(10, "x", []int32{1, 2, 3})),
			),
		)).String()

	if !strings.Contains(s, "AND") {
		t.Fatal("and")
	}
	if !strings.Contains(s, "CONST") {
		t.Fatal("const")
	}
	if !strings.Contains(s, "DisMax") {
		t.Fatal("dismax")
	}

	if !strings.Contains(s, "x") {
		t.Fatal("term")
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
		DisMax(0.1, Term(10, "x", []int32{1, 2, 3, 4}), Term(10, "x", []int32{1, 2, 4}), CreateFileTerm(10, "x", []int32{1, 4})),
	))

	qu := Term(10, "x", []int32{1, 2, 3, 4})
	qu.SetBoost(0)
	eqF(t, []float32{0, 0, 0, 0}, queryScores(
		DisMax(0.1, qu),
	))

	qu = Term(10, "x", []int32{1, 2, 3, 4})
	qu.SetBoost(1)
	eqF(t, []float32{1.2527629, 1.2527629, 1.2527629, 1.2527629}, queryScores(
		And(DisMax(0.1, qu)),
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
	old_chunk_size := TERM_CHUNK_SIZE
	for _, s := range []int{1, 2, 32, 64, 4096, math.MaxInt32} {
		TERM_CHUNK_SIZE = s
		k := rand.Intn(65000)
		a := postingsList(100 + k)
		b := postingsList(1000+k, a)
		c := postingsList(10000+k, a, b)
		d := postingsList(100000+k, a, b, c)
		e := postingsList(1000000+k, a, b, c, d)

		eq(t, a, query(Term(10, "x", a)))
		eq(t, a, query(CreateFileTerm(10, "x", a)))
		eq(t, b, query(Term(10, "x", b)))
		eq(t, c, query(Term(10, "x", c)))
		eq(t, d, query(Term(10, "x", d)))
		eq(t, e, query(Term(10, "x", e)))

		eq(t, b, query(Or(
			Term(10, "x", a),
			CreateFileTerm(10, "x", a),
			Term(10, "x", b),
			CreateFileTerm(10, "x", b),
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
			CreateFileTerm(10, "x", a),
			Term(10, "x", b),
			Term(10, "x", c),
			Term(10, "x", d),
			CreateFileTerm(10, "x", d),
			Term(10, "x", e),
		)))

		eq(t, a, query(And(
			DisMax(1,
				Term(10, "x", a),
				Term(10, "x", b),
				Term(10, "x", c),
				Term(10, "x", d),
				Term(10, "x", e)),
			Term(10, "x", a),
		)))

		eq(t, a, query(And(
			Or(Term(10, "x", a),
				Term(10, "x", b),
				Term(10, "x", c),
				Term(10, "x", d),
				Term(10, "x", e)),
			Term(10, "x", a),
		)))

		eq(t, a, query(And(
			Term(10, "x", a),
			Term(10, "x", b),
			Term(10, "x", c),
			Term(10, "x", d),
			Term(10, "x", e),
		)))

		eq(t, b, query(And(
			Term(10, "x", b),
			Term(10, "x", c),
			Term(10, "x", d),
			Term(10, "x", e),
		)))

		eq(t, c, query(And(
			Term(10, "x", c),
			Term(10, "x", d),
			Term(10, "x", e),
		)))

		eq(t, d, query(And(
			Term(10, "x", d),
			Term(10, "x", e),
		)))
	}
	TERM_CHUNK_SIZE = old_chunk_size
	a := postingsList(100)
	b := postingsList(1000, a)
	c := postingsList(10000, a, b)
	d := postingsList(100000, a, b, c)
	e := postingsList(1000000, a, b, c, d)

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
			CreateFileTerm(10, "x", []int32{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14}),
			CreateFileTerm(10, "x", []int32{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14}),
		),
	)))

	eq(t, []int32{6, 7, 8, 10}, query(AndNot(
		Term(10, "x", []int32{1, 2, 3, 9}),
		AndNot(
			Term(10, "x", []int32{4, 5}),
			Term(10, "x", []int32{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}),
			Term(10, "x", []int32{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}),
			CreateFileTerm(10, "x", []int32{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14}),
			CreateFileTerm(10, "x", []int32{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14}),
		),
	)))

	eq(t, []int32{1, 7, 1001}, query(AndNot(
		nil,
		CreateFileTerm(10, "x", []int32{1, 3, 5, 7, 100, 1001}),
		CreateFileTerm(10, "x", []int32{1, 4, 7, 10, 1000, 1001}),
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
