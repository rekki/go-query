package main

import (
	"testing"

	iq "github.com/jackdoe/go-query"
	rq "github.com/jackdoe/roaring-query"
)

var i32, ir = DoIndex("./list")

func TestIsSame(t *testing.T) {
	var matches uint64
	{
		m := ir
		x := m["Lorem"]
		y := m["corpora"]
		sum := uint64(0)
		q := rq.And(rq.Term("", x), rq.Term("", y))
		iter := q.Iterator()
		for iter.HasNext() {
			iter.Next()
			sum++
		}
		matches = sum
	}

	{
		m := i32
		x := m["Lorem"]
		y := m["corpora"]
		sum := uint64(0)

		q := iq.And(iq.Term("", x), iq.Term("", y))
		for q.Next() != iq.NO_MORE {
			sum++
		}
		if matches != sum {
			t.Fatalf("expected roaring: %d, got iunverted: %d", matches, sum)
		}
	}

}

func BenchmarkRoaringScanAndTwo(b *testing.B) {
	m := ir

	x := m["Lorem"]
	y := m["corpora"]

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		sum := uint32(0)
		q := rq.And(rq.Term("", x), rq.Term("", y))
		iter := q.Iterator()
		for iter.HasNext() {
			sum += iter.Next()
		}
	}
}

func BenchmarkInvertedScanAndTwo(b *testing.B) {
	m := i32

	x := m["Lorem"]
	y := m["corpora"]
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		sum := int32(0)
		q := iq.And(iq.Term("", x), iq.Term("", y))
		for q.Next() != iq.NO_MORE {
			sum += q.GetDocId()
		}
	}
}

func BenchmarkRoaringScanAndOne(b *testing.B) {
	m := ir

	x := m["Lorem"]

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		sum := uint32(0)
		q := rq.And(rq.Term("", x))
		iter := q.Iterator()
		for iter.HasNext() {
			sum += iter.Next()
		}
	}
}

func BenchmarkInvertedScanAndOne(b *testing.B) {
	m := i32

	x := m["Lorem"]
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		sum := int32(0)
		q := iq.And(iq.Term("", x))
		for q.Next() != iq.NO_MORE {
			sum += q.GetDocId()
		}
	}
}

func BenchmarkRoaringScanTerm(b *testing.B) {

	x := ir["Lorem"]
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		sum := uint32(0)
		q := rq.Term("", x)
		iter := q.Iterator()
		for iter.HasNext() {
			sum += iter.Next()
		}
	}
}

func BenchmarkInvertedScanTerm(b *testing.B) {
	m := i32

	x := m["Lorem"]
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		sum := int32(0)
		q := iq.Term("", x)
		for q.Next() != iq.NO_MORE {
			sum += q.GetDocId()
		}
	}
}

func BenchmarkRoaringScanOr(b *testing.B) {
	m := ir

	x := m["Lorem"]
	y := m["corpora"]
	z := m["qui"]
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		sum := uint32(0)
		q := rq.Or(rq.Term("", x), rq.Term("", y), rq.Term("", z))
		iter := q.Iterator()
		for iter.HasNext() {
			sum += iter.Next()
		}
	}
}

func BenchmarkInvertedScanOr(b *testing.B) {
	m := i32

	x := m["Lorem"]
	y := m["corpora"]
	z := m["qui"]
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		sum := int32(0)
		q := iq.Or(iq.Term("", x), iq.Term("", y), iq.Term("", z))
		for q.Next() != iq.NO_MORE {
			sum += q.GetDocId()
		}
	}
}

func BenchmarkRoaringScanAnd(b *testing.B) {
	m := ir

	x := m["Lorem"]
	y := m["corpora"]
	z := m["qui"]

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		sum := uint32(0)
		q := rq.And(rq.Term("", x), rq.Term("", y), rq.Term("", z))
		iter := q.Iterator()
		for iter.HasNext() {
			sum += iter.Next()
		}
	}
}

func BenchmarkInvertedScanAnd(b *testing.B) {
	m := i32

	x := m["Lorem"]
	y := m["corpora"]
	z := m["qui"]
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		sum := int32(0)
		q := iq.And(iq.Term("", x), iq.Term("", y), iq.Term("", z))
		for q.Next() != iq.NO_MORE {
			sum += q.GetDocId()
		}
	}
}

func BenchmarkRoaringScanAndNot(b *testing.B) {
	m := ir

	x := m["Lorem"]
	y := m["corpora"]
	z := m["qui"]
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		sum := uint32(0)
		q := rq.AndNot(rq.Term("", z), rq.Term("", y), rq.Term("", x))
		iter := q.Iterator()
		for iter.HasNext() {
			sum += iter.Next()
		}
	}
}

func BenchmarkInvertedScanAndNot(b *testing.B) {
	m := i32

	x := m["Lorem"]
	y := m["corpora"]
	z := m["qui"]
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		sum := int32(0)
		q := iq.And(iq.Term("", z), iq.Term("", y), iq.Term("", x))
		for q.Next() != iq.NO_MORE {
			sum += q.GetDocId()
		}
	}
}

func BenchmarkRoaringScanAndCompex(b *testing.B) {
	m := ir

	x := m["Lorem"]
	y := m["corpora"]
	z := m["qui"]
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		sum := uint32(0)
		q := rq.And(rq.Term("", z), rq.Or(rq.And(rq.Term("", y), rq.Term("", x)), rq.Term("", y), rq.Term("", x)))
		iter := q.Iterator()
		for iter.HasNext() {
			sum += iter.Next()
		}
	}
}

func BenchmarkInvertedScanAndCompex(b *testing.B) {
	m := i32

	x := m["Lorem"]
	y := m["corpora"]
	z := m["qui"]
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		sum := int32(0)
		q := iq.And(iq.Term("", z), iq.Or(iq.And(iq.Term("", y), iq.Term("", x)), iq.Term("", y), iq.Term("", x)))
		for q.Next() != iq.NO_MORE {
			sum += q.GetDocId()
		}
	}
}
