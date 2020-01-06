## github.com/rekki/go-query: simple []int32 query library

[![Build Status](https://travis-ci.org/rekki/go-query.svg?branch=master)](https://travis-ci.org/rekki/go-query) [![codecov](https://codecov.io/gh/rekki/go-query/branch/master/graph/badge.svg)](https://codecov.io/gh/rekki/go-query) [![GoDoc](https://godoc.org/github.com/rekki/go-query?status.svg)](https://godoc.org/github.com/rekki/go-query)


used to build and execute queries such as:

```
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

* scoring: only idf scode (for now)
* supported queries: or, and, and_not, dis_max, constant, term
* util/norm: space_between_digits, lowercase, trim, cleanup, ... [![GoDoc](https://godoc.org/github.com/rekki/go-query/util/norm?status.svg)](https://godoc.org/github.com/rekki/go-query/util/norm)
* util/tokenize: left edge, custom, charngram, unique, soundex, ... [![GoDoc](https://godoc.org/github.com/rekki/go-query/util/tokenize?status.svg)](https://godoc.org/github.com/rekki/go-query/util/tokenize)
* util/memory index: useful example of how to build more complex search engine with the library [![GoDoc](https://godoc.org/github.com/rekki/go-query/util/index?status.svg)](https://godoc.org/github.com/rekki/go-query/util/index)


---
# query
--
    import "github.com/rekki/go-query"

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

will print:

    matching 1, score: 2.639057
    matching 2, score: 2.639057
    matching 3, score: 2.852632
    matching 8, score: 2.639057
    matching 9, score: 4.105394

## Usage

```go
const (
	NO_MORE   = int32(math.MaxInt32)
	NOT_READY = int32(-1)
)
```

```go
var TERM_CHUNK_SIZE = 4096
```
splits the postings list into chunks that are binary searched and inside each
chunk linearly searching for next advance()

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

#### func  Constant

```go
func Constant(boost float32, q Query) *constantQuery
```

#### func  DisMax

```go
func DisMax(tieBreaker float32, queries ...Query) *disMaxQuery
```
Creates DisMax query, for example if the query is:

    DisMax(0.5, "name:amsterdam","name:university","name:free")

lets say we have an index with following idf: amsterdam: 1.3, free: 0.2,
university: 2.1 the score is computed by:

    max(score(amsterdam),score(university), score(free)) = 2.1 (university)
    + score(free) * tiebreaker = 0.1
    + score(amsterdam) * tiebreaker = 0.65
    = 2.85

#### func  Or

```go
func Or(queries ...Query) *orQuery
```
Creates OR query

#### func  Term

```go
func Term(totalDocumentsInIndex int, t string, postings []int32) *termQuery
```
Basic []int32{} that the whole interface works on top score is IDF (not tf*idf,
just 1*idf, since we dont store the term frequency for now) if you dont know
totalDocumentsInIndex, which could be the case sometimes, pass any constant > 0
WARNING: the query *can not* be reused WARNING: the query it not thread safe

#### type Query

```go
type Query interface {
	Next() int32
	GetDocId() int32
	Score() float32
	SetBoost(float32) Query

	String() string
	// contains filtered or unexported methods
}
```

Reuse/Concurrency: None of the queries are safe to be re-used. WARNING: the
query *can not* be reused WARNING: the query it not thread safe

Example Iteration:

    q := Term([]int32{1,2,3})
    for q.Next() != query.NO_MORE {
    	did := q.GetDocId()
    	score := q.Score()
    	fmt.Printf("matching %d, score: %f\n", did, score)
    }
---
# util
--
    import "github.com/rekki/go-query/util"

Simlpe utils to tokenize and normalize text

Example:

    package main

    import (
    	"fmt"

    	n "github.com/rekki/go-query/util/norm"
    	t "github.com/rekki/go-query/util/tokenize"
    )

    func main() {
    	tokenizer := []t.Tokenizer{
    		t.NewWhitespace(),
    		t.NewLeftEdge(1),
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

## Usage
---
# index
--
    import "github.com/rekki/go-query/util/index"

Illustration of how you can use go-query to build a somewhat functional search
index Example:

    package main

    import (
    	"log"

    	iq "github.com/rekki/go-query"
    	"github.com/rekki/go-query/util/analyzer"
    	"github.com/rekki/go-query/util/index"
    )

    type ExampleCity struct {
    	Name    string
    	Country string
    }

    func (e *ExampleCity) IndexableFields() map[string][]string {
    	out := map[string][]string{}

    	out["name"] = []string{e.Name}
    	out["name_fuzzy"] = []string{e.Name}
    	out["name_soundex"] = []string{e.Name}
    	out["country"] = []string{e.Country}

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
    	m := index.NewMemOnlyIndex(map[string]*analyzer.Analyzer{
    		"name":         index.AutocompleteAnalyzer,
    		"name_fuzzy":   index.FuzzyAnalyzer,
    		"name_soundex": index.SoundexAnalyzer,
    		"country":      index.IDAnalyzer,
    	})

    	list := []*ExampleCity{
    		&ExampleCity{Name: "Amsterdam", Country: "NL"},
    		&ExampleCity{Name: "Amsterdam University", Country: "NL"},
    		&ExampleCity{Name: "Amsterdam University Second", Country: "NL"},
    		&ExampleCity{Name: "London", Country: "UK"},
    		&ExampleCity{Name: "Sofia", Country: "BG"},
    	}

    	m.Index(toDocuments(list)...)

    	// search for "(name:aMS OR name:u) AND (country:NL OR country:BG)"

    	query := iq.Or(
    		iq.And(
    			iq.Or(m.Terms("name", "aMS u")...),
    			iq.Or(m.Terms("country", "NL")...),
    		).SetBoost(2),
    		iq.And(
    			iq.Or(m.Terms("name_fuzzy", "bondom u")...),
    			iq.Or(m.Terms("country", "UK")...),
    		).SetBoost(0.1),
    		iq.And(
    			iq.Or(m.Terms("name_soundex", "sfa")...),
    			iq.Or(m.Terms("country", "BG")...),
    		).SetBoost(0.01),
    	)
    	log.Printf("query: %v", query.String())
    	m.Foreach(query, func(did int32, score float32, doc index.Document) {
    		city := doc.(*ExampleCity)
    		log.Printf("%v matching with score %f", city, score)
    	})
    }

will print

    2019/12/03 18:40:11 &{Amsterdam NL} matching with score 3.923317
    2019/12/03 18:40:11 &{Amsterdam University NL} matching with score 6.428843
    2019/12/03 18:40:11 &{Amsterdam University NL Second} matching with score 6.428843
    2019/12/03 18:40:11 &{London UK} matching with score 0.537528
    2019/12/03 18:40:11 &{Sofia BG} matching with score 0.035835

## Usage

```go
var AutocompleteAnalyzer = analyzer.NewAnalyzer(DefaultNormalizer, DefaultSearchTokenizer, AutocompleteIndexTokenizer)
```

```go
var AutocompleteIndexTokenizer = []tokenize.Tokenizer{
	tokenize.NewWhitespace(),
	tokenize.NewLeftEdge(1),
}
```

```go
var DefaultAnalyzer = analyzer.NewAnalyzer(DefaultNormalizer, DefaultSearchTokenizer, DefaultIndexTokenizer)
```

```go
var DefaultIndexTokenizer = []tokenize.Tokenizer{
	tokenize.NewWhitespace(),
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
}
```

```go
var FuzzyAnalyzer = analyzer.NewAnalyzer(DefaultNormalizer, FuzzyTokenizer, FuzzyTokenizer)
```

```go
var FuzzyTokenizer = []tokenize.Tokenizer{
	tokenize.NewWhitespace(),
	tokenize.NewCharNgram(2),
	tokenize.NewUnique(),
	tokenize.NewSurround("$"),
}
```

```go
var IDAnalyzer = analyzer.NewAnalyzer([]norm.Normalizer{norm.NewNoop()}, []tokenize.Tokenizer{tokenize.NewNoop()}, []tokenize.Tokenizer{tokenize.NewNoop()})
```

```go
var SoundexAnalyzer = analyzer.NewAnalyzer(DefaultNormalizer, SoundexTokenizer, SoundexTokenizer)
```

```go
var SoundexTokenizer = []tokenize.Tokenizer{
	tokenize.NewWhitespace(),
	tokenize.NewSoundex(),
}
```

#### func  Parse

```go
func Parse(input *spec.Query, makeTermQuery func(string, string) iq.Query) (iq.Query, error)
```
Take spec.Query object and a makeTermQuery function and produce a parsed query
Example:

    return Parse(input, func(k, v string) iq.Query {
    	kv := k + ":"+ v
    	return iq.Term(0, kv, postings[kv])
    })

#### func  QueryFromBytes

```go
func QueryFromBytes(b []byte) (*spec.Query, error)
```
somewhat useless method (besides for testing) Example:

    query, err := QueryFromBytes([]byte(`{
      "type": "OR",
      "queries": [
        {
          "field": "name",
          "value": "sofia"
        },
        {
          "field": "name",
          "value": "amsterdam"
        }
      ]
    }`))
    if err != nil {
    	panic(err)
    }

#### func  QueryFromJson

```go
func QueryFromJson(input interface{}) (*spec.Query, error)
```
simple (*slow*) helper method that takes interface{} and converst it to
spec.Query with jsonpb in case you receive request like request = {"limit":10,
query: ....}, pass request.query to QueryFromJson and get a query object back

#### type Document

```go
type Document interface {
	IndexableFields() map[string][]string
}
```

Export this interface on the documents you want indexed

Example if you want to index fields "name" and "country":

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

#### type Hit

```go
type Hit struct {
	Score    float32  `json:"score"`
	Id       int32    `json:"id"`
	Document Document `json:"doc"`
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

    query := iq.And(
    	iq.Or(m.Terms("name", "aMS u")...),
    	iq.Or(m.Terms("country", "NL BG")...),
    )
    m.Foreach(query, func(did int32, score float32, doc index.Document) {
    	city := doc.(*ExampleCity)
    	log.Printf("%v matching with score %f", city, score)
    })

#### func (*MemOnlyIndex) Index

```go
func (m *MemOnlyIndex) Index(docs ...Document)
```
index a bunch of documents

#### func (*MemOnlyIndex) Parse

```go
func (m *MemOnlyIndex) Parse(input *spec.Query) (iq.Query, error)
```
Parse dsl input into an query object Example:

    query, err := QueryFromBytes([]byte(`{
      "type": "OR",
      "queries": [
        {
          "field": "name",
          "value": "sofia"
        },
        {
          "field": "name",
          "value": "amsterdam"
        }
      ]
    }`))
    if err != nil {
    	panic(err)
    }
    parsedQuery, err := m.Parse(query)
    if err != nil {
    	panic(err)
    }
    top = m.TopN(1, parsedQuery, nil)
    ...
    {
      "total": 3,
      "hits": [
        {
          "score": 1.609438,
          "id": 3,
          "doc": {
            "Name": "Sofia",
            "Country": "BG"
          }
        }
      ]
    }

#### func (*MemOnlyIndex) Terms

```go
func (m *MemOnlyIndex) Terms(field string, term string) []iq.Query
```
Generate array of queries from the tokenized term for this field, using the
perField analyzer

#### func (*MemOnlyIndex) TopN

```go
func (m *MemOnlyIndex) TopN(limit int, query iq.Query, cb func(int32, float32, Document) float32) *SearchResult
```
TopN documents The following texample gets top5 results and also check add 100
to the score of cities that have NL in the score. usually the score of your
search is some linear combination of f(a*text + b*popularity + c*context..)

Example:

    query := iq.And(
    	iq.Or(m.Terms("name", "ams university")...),
    	iq.Or(m.Terms("country", "NL BG")...),
    )
    top := m.TopN(5, q, func(did int32, score float32, doc Document) float32 {
    	city := doc.(*ExampleCity)
    	if city.Country == "NL" {
    		score += 100
    	}
    	n++
    	return score
    })

the SearchResult structure looks like

    {
      "total": 3,
      "hits": [
        {
          "score": 101.09861,
          "id": 0,
          "doc": {
            "Name": "Amsterdam",
            "Country": "NL"
          }
        }
        ...
      ]
    }

If the callback is null, then the original score is used (1*idf at the moment)

#### type SearchResult

```go
type SearchResult struct {
	Total int   `json:"total"`
	Hits  []Hit `json:"hits"`
}
```
