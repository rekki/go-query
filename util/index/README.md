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
    )

    type ExampleCity struct {
    	Name    string
    	Country string
    }

    func (e *ExampleCity) IndexableFields() map[string]string {
    	out := map[string]string{}

    	out["name"] = e.Name
    	out["name_fuzzy"] = e.Name
    	out["name_soundex"] = e.Name
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
	tokenize.NewUnique(),
}
```

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
	tokenize.NewUnique(),
}
```

#### func  Ld

```go
func Ld(s, t string) int
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
