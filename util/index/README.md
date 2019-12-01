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
