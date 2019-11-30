# query
--
    import "github.com/jackdoe/go-query"

Package query provides simple query dsl on top of sorted arrays of integers.
Usually when you have inverted index you endup having something like:

    data := []*Document{}
    index := map[string][]int32{}
    for docId, d := range document {
    	for _, token := range tokenize(normalize(d.Name)) {
            index[token] = append(index[token],docId)
        }
    }

then from documents like {hello world},{hello},{new york},{new world} you get
inverted index in the form of:

    {
       "hello": [0,1],
       "world": [0,3],
       "new": [2,3],
       "york": [2]
    }

anyway, if you want to read more on those check out the IR-book

This package helps you query indexes of this form, in fairly efficient way, keep
in mind it expects the []int32 slices to be _sorted_. Example:

    package main

    import (
    	"fmt"

    	"github.com/jackdoe/go-query"
    )

    func main() {
    	q := query.And(
    		query.Or(
    			query.Term("a", []int32{1, 2}),
    			query.Term("b", []int32{3, 9})),
    		query.AndNot(
    			query.Or(query.Term("c", []int32{4, 5}), query.Term("c", []int32{4, 100})),
    			query.Or(
    				query.Term("d", []int32{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}),
    				query.Term("e", []int32{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}),
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

will print:

    matching 1, score: 2.000000
    matching 2, score: 2.000000
    matching 3, score: 2.000000
    matching 9, score: 2.000000

## Usage

```go
const (
	NO_MORE   = int32(math.MaxInt32)
	NOT_READY = int32(-1)
)
```

#### func  And

```go
func And(queries ...Query) *andQuery
```
Creates AND query

#### func  AndNot

```go
func AndNot(not Query, queries ...Query) *andQuery
```
Creates AND NOT query

#### func  AndTerm

```go
func AndTerm(queries ...*termQuery) *andTerm
```
Creates AND query

#### func  Or

```go
func Or(queries ...Query) *orQuery
```
Creates OR query

#### func  Term

```go
func Term(t string, postings []int32) *termQuery
```
Basic []int32{} that the whole interface works on top score is unnormalized IDF,
there is no term frequency

#### type Query

```go
type Query interface {
	Next() int32
	GetDocId() int32
	Score() float32

	String() string
	// contains filtered or unexported methods
}
```

Reuse/Concurrency: None of the queries are safe to be re-used.

Example Iteration:

        q := Term([]int32{1,2,3})
    	for q.Next() != query.NO_MORE {
    		did := q.GetDocId()
    		score := q.Score()
    		fmt.Printf("matching %d, score: %f\n", did, score)
    	}

----------------------------

# util
--
Simlpe utils to tokenize and normalize text

Example: package main

import (

    "fmt"

    n "github.com/jackdoe/go-query/util/norm"
    t "github.com/jackdoe/go-query/util/tokenize"

)

func main() {

    tokenizer := []t.Tokenizer{
    	t.NewWhitespace(),
    	t.NewLeftEdge(1),
    	t.NewUnique(),
    }
    normalizer := []n.Normalizer{
    	n.NewUnaccent(),
    	n.NewLowerCase(),
    	n.NewSpaceBetweenDigits(),
    	n.NewCleanup(n.BASIC_NON_ALPHANUMERIC),
    	n.NewTrim(" "),
    }

    tokens := t.Tokenize(
    	n.Normalize("Hęllö World yęar2019 ", normalizer...),
    	tokenizer...,
    )

    fmt.Printf("%v", tokens)
    // prints [h he hel hell hello w wo wor worl world y ye yea year 2 20 201 2019]

}

----------------------------

# index
--
    import "github.com/jackdoe/go-query/util/index"

Illustration of how you can use go-query to build a somewhat functional search
index Example:

    package main

    import (
    	"log"

    	iq "github.com/jackdoe/go-query"
    	"github.com/jackdoe/go-query/util/analyzer"
    	"github.com/jackdoe/go-query/util/index"
    	"github.com/jackdoe/go-query/util/tokenize"
    )

    type ExampleCity struct {
    	Name    string
    	Country string
    }

    func (e *ExampleCity) IndexableFields() map[string]string {
    	out := map[string]string{}

    	out["name"] = e.Name
    	out["country"] = e.Country

    	return out
    }

    func toDocuments(in []*ExampleCity) []index.Document {
    	out := make([]index.Document, len(in))
    	for i, d := range in {
    		out[i] = index.Document(d)
    	}
    	return out
    }

    func main() {

    	indexTokenizer := []tokenize.Tokenizer{
    		tokenize.NewWhitespace(),
    		tokenize.NewLeftEdge(1), // left edge ngram indexing for prefix matches
    		tokenize.NewUnique(),
    	}

    	searchTokenizer := []tokenize.Tokenizer{
    		tokenize.NewWhitespace(),
    		tokenize.NewUnique(),
    	}

    	autocomplete := analyzer.NewAnalyzer(index.DefaultNormalizer, searchTokenizer, indexTokenizer)
    	m := index.NewMemOnlyIndex(map[string]*analyzer.Analyzer{
    		"name":    autocomplete,
    		"country": index.DefaultAnalyzer,
    	})

    	list := []*ExampleCity{
    		&ExampleCity{Name: "Amsterdam", Country: "NL"},
    		&ExampleCity{Name: "Amsterdam University", Country: "NL"},
    		&ExampleCity{Name: "Amsterdam University", Country: "NL"},
    		&ExampleCity{Name: "London", Country: "UK"},
    		&ExampleCity{Name: "Sofia", Country: "BG"},
    	}

    	m.Index(toDocuments(list)...)

    	// search for "(name:aMS OR name:u) AND (country:NL OR country:BG)"

    	query := iq.And(m.Or("name", "aMS u"), m.Or("country", "NL BG"))

    	m.Foreach(query, func(did int32, score float32, doc index.Document) {
    		city := doc.(*ExampleCity)
    		log.Printf("%v matching with score %f", city, score)
    	})
    }

## Usage

```go
var DefaultAnalyzer = analyzer.NewAnalyzer(DefaultNormalizer, DefaultSearchTokenizer, DefaultIndexTokenizer)
```

```go
var DefaultIndexTokenizer = []tokenize.Tokenizer{
	tokenize.NewWhitespace(),
	tokenize.NewUnique(),
}
```

```go
var DefaultNormalizer = []norm.Normalizer{
	norm.NewUnaccent(),
	norm.NewLowerCase(),
	norm.NewSpaceBetweenDigits(),
	norm.NewCleanup(norm.BASIC_NON_ALPHANUMERIC),
	norm.NewTrim(" "),
}
```

```go
var DefaultSearchTokenizer = []tokenize.Tokenizer{
	tokenize.NewWhitespace(),
	tokenize.NewUnique(),
}
```

#### type Document

```go
type Document interface {
	IndexableFields() map[string]string
}
```


#### type MemOnlyIndex

```go
type MemOnlyIndex struct {
	sync.RWMutex
}
```


#### func  NewMemOnlyIndex

```go
func NewMemOnlyIndex(perField map[string]*analyzer.Analyzer) *MemOnlyIndex
```
create new in-memory index with the specified perField analyzer by default
DefaultAnalyzer is used

#### func (*MemOnlyIndex) Foreach

```go
func (m *MemOnlyIndex) Foreach(query iq.Query, cb func(int32, float32, Document))
```
Foreach matching document Example:

    q := query.Or("name", "amster")
    m.Foreach(query, func(did int32, score float32, doc index.Document) {
    	city := doc.(*ExampleCity)
    	log.Printf("%v matching with score %f", city, score)
    })

#### func (*MemOnlyIndex) Index

```go
func (m *MemOnlyIndex) Index(docs ...Document)
```
index a bunch of documents

#### func (*MemOnlyIndex) Terms

```go
func (m *MemOnlyIndex) Terms(field string, term string) []iq.Query
```
Generate array of queries from the tokenized term for this field, using the
perField analyzer
