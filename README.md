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
	SetBoost(float32)

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

----------------------------

# util
--
Simlpe utils to tokenize and normalize text

Example:

    package main

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

    	autocomplete := analyzer.NewAnalyzer(
    		index.DefaultNormalizer,
    		searchTokenizer,
    		indexTokenizer,
    	)
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

    	query := iq.And(
    		iq.Or(m.Terms("name", "aMS u")...),
    		iq.Or(m.Terms("country", "NL BG")...),
    	)

    	m.Foreach(query, func(did int32, score float32, doc index.Document) {
    		city := doc.(*ExampleCity)
    		log.Printf("%v matching with score %f", city, score)
    	})
    }

will print

    2019/11/30 18:20:23 &{Amsterdam NL} matching with score 1.961658
    2019/11/30 18:20:23 &{Amsterdam University NL} matching with score 3.214421
    2019/11/30 18:20:23 &{Amsterdam University NL} matching with score 3.214421

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
