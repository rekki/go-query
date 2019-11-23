package main

import (
	"testing"

	iq "github.com/jackdoe/go-query"
	rq "github.com/jackdoe/roaring-query"
)

var i32, ir = DoIndex("./list")

func BenchmarkRoaringScanTerm(b *testing.B) {

	x := ir["Lorem"]
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		sum := uint32(0)
		q := rq.NewTerm("", x)
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
		q := rq.NewBoolOrQuery(rq.NewTerm("", x), rq.NewTerm("", y), rq.NewTerm("", z))
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
		q := rq.NewBoolAndQuery(rq.NewTerm("", x), rq.NewTerm("", y), rq.NewTerm("", z))
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
		q := rq.NewBoolAndNotQuery(rq.NewTerm("", z), rq.NewTerm("", y), rq.NewTerm("", x))
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
		q := rq.NewBoolAndQuery(rq.NewTerm("", z), rq.NewBoolOrQuery(rq.NewBoolAndQuery(rq.NewTerm("", y), rq.NewTerm("", x)), rq.NewTerm("", y), rq.NewTerm("", x)))
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
