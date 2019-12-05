package norm

import (
	"testing"
)

func BenchmarkAlloc(b *testing.B) {
	b.ReportAllocs()

	s := "hello cat dog"
	ps := NewPorterStemmer()
	cnt := 0

	for n := 0; n < b.N; n++ {
		cnt += len(ps.Apply(s))
	}
}
