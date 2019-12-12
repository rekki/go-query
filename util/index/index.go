// Illustration of how you can use go-query to build a somewhat functional search index
// Example:
//  package main
//
//  import (
//  	"log"
//
//  	iq "github.com/rekki/go-query"
//  	"github.com/rekki/go-query/util/analyzer"
//  	"github.com/rekki/go-query/util/index"
//  )
//
//  type ExampleCity struct {
//  	Name    string
//  	Country string
//  }
//
//  func (e *ExampleCity) IndexableFields() map[string][]string {
//  	out := map[string][]string{}
//
//  	out["name"] = []string{e.Name}
//  	out["name_fuzzy"] = []string{e.Name}
//  	out["name_soundex"] = []string{e.Name}
//  	out["country"] = []string{e.Country}
//
//  	return out
//  }
//
//  func toDocuments(in []*ExampleCity) []index.Document {
//  	out := make([]index.Document, len(in))
//  	for i, d := range in {
//  		out[i] = index.Document(d)
//  	}
//  	return out
//  }
//
//  func main() {
//  	m := index.NewMemOnlyIndex(map[string]*analyzer.Analyzer{
//  		"name":         index.AutocompleteAnalyzer,
//  		"name_fuzzy":   index.FuzzyAnalyzer,
//  		"name_soundex": index.SoundexAnalyzer,
//  		"country":      index.IDAnalyzer,
//  	})
//
//  	list := []*ExampleCity{
//  		&ExampleCity{Name: "Amsterdam", Country: "NL"},
//  		&ExampleCity{Name: "Amsterdam University", Country: "NL"},
//  		&ExampleCity{Name: "Amsterdam University Second", Country: "NL"},
//  		&ExampleCity{Name: "London", Country: "UK"},
//  		&ExampleCity{Name: "Sofia", Country: "BG"},
//  	}
//
//  	m.Index(toDocuments(list)...)
//
//  	// search for "(name:aMS OR name:u) AND (country:NL OR country:BG)"
//
//  	query := iq.Or(
//  		iq.And(
//  			iq.Or(m.Terms("name", "aMS u")...),
//  			iq.Or(m.Terms("country", "NL")...),
//  		).SetBoost(2),
//  		iq.And(
//  			iq.Or(m.Terms("name_fuzzy", "bondom u")...),
//  			iq.Or(m.Terms("country", "UK")...),
//  		).SetBoost(0.1),
//  		iq.And(
//  			iq.Or(m.Terms("name_soundex", "sfa")...),
//  			iq.Or(m.Terms("country", "BG")...),
//  		).SetBoost(0.01),
//  	)
//  	log.Printf("query: %v", query.String())
//  	m.Foreach(query, func(did int32, score float32, doc index.Document) {
//  		city := doc.(*ExampleCity)
//  		log.Printf("%v matching with score %f", city, score)
//  	})
//  }
// will print
//  2019/12/03 18:40:11 &{Amsterdam NL} matching with score 3.923317
//  2019/12/03 18:40:11 &{Amsterdam University NL} matching with score 6.428843
//  2019/12/03 18:40:11 &{Amsterdam University NL Second} matching with score 6.428843
//  2019/12/03 18:40:11 &{London UK} matching with score 0.537528
//  2019/12/03 18:40:11 &{Sofia BG} matching with score 0.035835
package index

import (
	"fmt"
	"sync"

	iq "github.com/rekki/go-query"
	"github.com/rekki/go-query/util/analyzer"
	spec "github.com/rekki/go-query/util/go_query_dsl"
	"github.com/rekki/go-query/util/norm"
	"github.com/rekki/go-query/util/tokenize"
)

// Export this interface on the documents you want indexed
//
// Example if you want to index fields "name" and "country":
//  type ExampleCity struct {
//  	Name    string
//  	Country string
//  }
//
//  func (e *ExampleCity) IndexableFields() map[string]string {
//  	out := map[string]string{}
//
//  	out["name"] = e.Name
//  	out["country"] = e.Country
//
//  	return out
//  }
type Document interface {
	IndexableFields() map[string][]string
}

var DefaultNormalizer = []norm.Normalizer{
	norm.NewUnaccent(),
	norm.NewLowerCase(),
	norm.NewSpaceBetweenDigits(),
	norm.NewCleanup(norm.BASIC_NON_ALPHANUMERIC),
	norm.NewTrim(" "),
}

var DefaultSearchTokenizer = []tokenize.Tokenizer{
	tokenize.NewWhitespace(),
}

var DefaultIndexTokenizer = []tokenize.Tokenizer{
	tokenize.NewWhitespace(),
}

var DefaultAnalyzer = analyzer.NewAnalyzer(DefaultNormalizer, DefaultSearchTokenizer, DefaultIndexTokenizer)

var IDAnalyzer = analyzer.NewAnalyzer([]norm.Normalizer{norm.NewNoop()}, []tokenize.Tokenizer{tokenize.NewNoop()}, []tokenize.Tokenizer{tokenize.NewNoop()})

var SoundexTokenizer = []tokenize.Tokenizer{
	tokenize.NewWhitespace(),
	tokenize.NewSoundex(),
}

var FuzzyTokenizer = []tokenize.Tokenizer{
	tokenize.NewWhitespace(),
	tokenize.NewCharNgram(2),
	tokenize.NewUnique(),
	tokenize.NewSurround("$"),
}

var AutocompleteIndexTokenizer = []tokenize.Tokenizer{
	tokenize.NewWhitespace(),
	tokenize.NewLeftEdge(1),
}

var SoundexAnalyzer = analyzer.NewAnalyzer(DefaultNormalizer, SoundexTokenizer, SoundexTokenizer)

var FuzzyAnalyzer = analyzer.NewAnalyzer(DefaultNormalizer, FuzzyTokenizer, FuzzyTokenizer)

var AutocompleteAnalyzer = analyzer.NewAnalyzer(DefaultNormalizer, DefaultSearchTokenizer, AutocompleteIndexTokenizer)

type MemOnlyIndex struct {
	perField map[string]*analyzer.Analyzer
	postings map[string]map[string][]int32
	forward  []Document
	sync.RWMutex
}

// create new in-memory index with the specified perField analyzer
// by default DefaultAnalyzer is used
func NewMemOnlyIndex(perField map[string]*analyzer.Analyzer) *MemOnlyIndex {
	if perField == nil {
		perField = map[string]*analyzer.Analyzer{}
	}
	m := &MemOnlyIndex{postings: map[string]map[string][]int32{}, perField: perField}
	return m
}

// index a bunch of documents
func (m *MemOnlyIndex) Index(docs ...Document) {
	m.Lock()
	defer m.Unlock()

	for _, d := range docs {
		fields := d.IndexableFields()
		did := len(m.forward)
		m.forward = append(m.forward, d)
		for field, value := range fields {
			analyzer, ok := m.perField[field]
			if !ok {
				analyzer = DefaultAnalyzer
			}
			for _, v := range value {
				tokens := analyzer.AnalyzeIndex(v)
				for _, t := range tokens {
					m.add(field, t, int32(did))
				}
			}
		}
	}
}

func (m *MemOnlyIndex) add(k, v string, did int32) {
	pk, ok := m.postings[k]
	if !ok {
		pk = map[string][]int32{}
		m.postings[k] = pk
	}

	current, ok := pk[v]
	if !ok || len(current) == 0 {
		pk[v] = []int32{did}
	} else {
		if current[len(current)-1] != did {
			pk[v] = append(current, did)
		}
	}
}

// Parse dsl input into an query object
// Example:
//  query, err := QueryFromBytes([]byte(`{
//    "type": "OR",
//    "queries": [
//      {
//        "field": "name",
//        "value": "sofia"
//      },
//      {
//        "field": "name",
//        "value": "amsterdam"
//      }
//    ]
//  }`))
//  if err != nil {
//  	panic(err)
//  }
//  parsedQuery, err := m.Parse(query)
//  if err != nil {
//  	panic(err)
//  }
//  top = m.TopN(1, parsedQuery, nil)
//  ...
//  {
//    "total": 3,
//    "hits": [
//      {
//        "score": 1.609438,
//        "id": 3,
//        "doc": {
//          "Name": "Sofia",
//          "Country": "BG"
//        }
//      }
//    ]
//  }
func (m *MemOnlyIndex) Parse(input *spec.Query) (iq.Query, error) {
	return Parse(input, func(k, v string) iq.Query {
		terms := m.Terms(k, v)
		if len(terms) == 1 {
			return terms[0]
		}
		return iq.Or(terms...)
	})
}

// Generate array of queries from the tokenized term for this field,
// using the perField analyzer
func (m *MemOnlyIndex) Terms(field string, term string) []iq.Query {
	m.RLock()
	defer m.RUnlock()
	analyzer, ok := m.perField[field]
	if !ok {
		analyzer = DefaultAnalyzer
	}
	tokens := analyzer.AnalyzeSearch(term)
	queries := []iq.Query{}
	for _, t := range tokens {
		queries = append(queries, m.newTermQuery(field, t))
	}
	return queries
}

func (m *MemOnlyIndex) newTermQuery(field string, term string) iq.Query {
	s := fmt.Sprintf("%s:%s", field, term)
	pk, ok := m.postings[field]
	if !ok {
		return iq.Term(len(m.forward), s, []int32{})
	}
	pv, ok := pk[term]
	if !ok {
		return iq.Term(len(m.forward), s, []int32{})
	}
	// there are allocation in iq.Term(), so dont just defer unlock, otherwise it will be locked while term is created
	return iq.Term(len(m.forward), s, pv)
}

// Foreach matching document
// Example:
//  query := iq.And(
//  	iq.Or(m.Terms("name", "aMS u")...),
//  	iq.Or(m.Terms("country", "NL BG")...),
//  )
//  m.Foreach(query, func(did int32, score float32, doc index.Document) {
//  	city := doc.(*ExampleCity)
//  	log.Printf("%v matching with score %f", city, score)
//  })
func (m *MemOnlyIndex) Foreach(query iq.Query, cb func(int32, float32, Document)) {
	for query.Next() != iq.NO_MORE {
		did := query.GetDocId()
		score := query.Score()

		cb(did, score, m.forward[did])
	}
}

// TopN documents
// The following texample gets top5 results and also check add 100 to the score of cities that have NL in the score.
// usually the score of your search is some linear combination of f(a*text + b*popularity + c*context..)
//
// Example:
//  query := iq.And(
//  	iq.Or(m.Terms("name", "ams university")...),
//  	iq.Or(m.Terms("country", "NL BG")...),
//  )
//  top := m.TopN(5, q, func(did int32, score float32, doc Document) float32 {
//  	city := doc.(*ExampleCity)
//  	if city.Country == "NL" {
//  		score += 100
//  	}
//  	n++
//  	return score
//  })
// the SearchResult structure looks like
//  {
//    "total": 3,
//    "hits": [
//      {
//        "score": 101.09861,
//        "id": 0,
//        "doc": {
//          "Name": "Amsterdam",
//          "Country": "NL"
//        }
//      }
//      ...
//    ]
//  }
// If the callback is null, then the original score is used (1*idf at the moment)
func (m *MemOnlyIndex) TopN(limit int, query iq.Query, cb func(int32, float32, Document) float32) *SearchResult {
	out := &SearchResult{}
	scored := []Hit{}
	m.Foreach(query, func(did int32, originalScore float32, d Document) {
		out.Total++
		if limit == 0 {
			return
		}
		score := originalScore
		if cb != nil {
			score = cb(did, originalScore, d)
		}

		// just keep the list sorted
		// FIXME: use bounded priority queue
		doInsert := false
		if len(scored) < limit {
			doInsert = true
		} else if scored[len(scored)-1].Score < score {
			doInsert = true
		}

		if doInsert {
			hit := Hit{Score: score, Id: did, Document: d}
			if len(scored) < limit {
				scored = append(scored, hit)
			}
			for i := 0; i < len(scored); i++ {
				if scored[i].Score < hit.Score {
					copy(scored[i+1:], scored[i:])
					scored[i] = hit
					break
				}
			}
		}
	})

	out.Hits = scored

	return out
}

type Hit struct {
	Score    float32  `json:"score"`
	Id       int32    `json:"id"`
	Document Document `json:"doc"`
}

type SearchResult struct {
	Total int   `json:"total"`
	Hits  []Hit `json:"hits"`
}
