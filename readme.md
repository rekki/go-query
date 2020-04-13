## go-query simple []int32 query library [![GitHub Actions Status](https://github.com/rekki/go-query/workflows/test/badge.svg?branch=master)](https://github.com/rekki/go-query/actions) [![codecov](https://codecov.io/gh/rekki/go-query/branch/master/graph/badge.svg)](https://codecov.io/gh/rekki/go-query) [![GoDoc](https://godoc.org/github.com/rekki/go-query?status.svg)](https://godoc.org/github.com/rekki/go-query)

> Blazingly fast query engine

Used to build and execute queries such as:

```go
n := 10 // total docs in index

And(
        Term(n, "name:hello", []int32{4, 5}),
        Term(n, "name:world", []int32{4, 100}),
        Or(
                Term(n, "country:nl", []int32{20,30}),
                Term(n, "country:uk", []int32{4,30}),
        )
)
```

- scoring: only `idf` score (for now)
- supported queries: `or`, `and`, `and_not`, `dis_max`, `constant`, `term`
- [`go-query-normalize`](https://github.com/rekki/go-query-normalize): space_between_digits, lowercase, trim, cleanup, etc
- [`go-query-tokenize`](https://github.com/rekki/go-query-tokenize): left edge, custom, charngram, unique, soundex etc
- [`go-query-index`](https://github.com/rekki/go-query-index): useful example of how to build more complex search engine with the library

---

# query

Usually when you have inverted index you endup having something like:

```go
data := []*Document{}
index := map[string][]int32{}
for docId, d := range document {
     for _, token := range tokenize(normalize(d.Name)) {
        index[token] = append(index[token],docId)
     }
}
```

then from documents like `{hello world}`, `{hello}`, `{new york}`, `{new world}` you get inverted index in the form of:

```go
{
    "hello": [0,1],
    "world": [0,3],
    "new": [2,3],
    "york": [2]
}
```

anyway, if you want to read more on those check out the IR-book

This package helps you query indexes of this form, in fairly efficient way, keep in mind it expects the `[]int32` slices to be _sorted_. Example:

```go
package main

import (
    "fmt"

    "github.com/rekki/go-query"
)

func main() {
    totalDocumentsInIndex := 10
    q := query.And(
        query.Or(
            query.Term(totalDocumentsInIndex, "a", []int32{1, 2, 8, 9}),
            query.Term(totalDocumentsInIndex, "b", []int32{3, 9, 8}),
        ),
        query.AndNot(
            query.Or(
                query.Term(totalDocumentsInIndex, "c", []int32{4, 5}),
                query.Term(totalDocumentsInIndex, "c", []int32{4, 100}),
            ),
            query.Or(
                query.Term(totalDocumentsInIndex, "d", []int32{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}),
                query.Term(totalDocumentsInIndex, "e", []int32{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}),
            ),
        ),
    )

    // q.String() is {{a OR b} AND {{d OR e} -({c OR x})}}

    for q.Next() != query.NO_MORE {
        did := q.GetDocId()
        score := q.Score()
        fmt.Printf("matching %d, score: %f\n", did, score)
    }
}
```

will print:

```sh
matching 1, score: 2.639057
matching 2, score: 2.639057
matching 3, score: 2.852632
matching 8, score: 2.639057
matching 9, score: 4.105394
```
