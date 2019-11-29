// Illustration of how you can use go-query to build a somewhat functional search index
// Example:
//  package main
//
//  import (
//  	"log"
//
//  	iq "github.com/jackdoe/go-query"
//  	"github.com/jackdoe/go-query/util/analyzer"
//  	"github.com/jackdoe/go-query/util/index"
//  	"github.com/jackdoe/go-query/util/tokenize"
//  )
//
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
//
//  	indexTokenizer := []tokenize.Tokenizer{
//  		tokenize.NewWhitespace(),
//  		tokenize.NewLeftEdge(1), // left edge ngram indexing for prefix matches
//  		tokenize.NewUnique(),
//  	}
//
//  	searchTokenizer := []tokenize.Tokenizer{
//  		tokenize.NewWhitespace(),
//  		tokenize.NewUnique(),
//  	}
//
//  	a := analyzer.NewAnalyzer(index.DefaultNormalizer, searchTokenizer, indexTokenizer)
//  	m := index.NewMemOnlyIndex(map[string]*analyzer.Analyzer{
//  		"name":    a,
//  		"country": index.DefaultAnalyzer,
//  	})
//
//  	list := []*ExampleCity{
//  		&ExampleCity{Name: "Amsterdam", Country: "NL"},
//  		&ExampleCity{Name: "Amsterdam University", Country: "NL"},
//  		&ExampleCity{Name: "Amsterdam University", Country: "NL"},
//  		&ExampleCity{Name: "London", Country: "UK"},
//  		&ExampleCity{Name: "Sofia", Country: "BG"},
//  	}
//
//  	m.Index(toDocuments(list)...)
//
//  	// search for "(name:aMS OR name:u) AND *country:NL OR country:BG)"
//
//  	query := iq.And(m.Or("name", "aMS u"), m.Or("country", "NL BG"))
//
//  	m.Foreach(query, func(did int32, score float32, doc index.Document) {
//  		city := doc.(*ExampleCity)
//  		log.Printf("%v matching with score %f", city, score)
//  	})
//  }
package index

import (
	"fmt"
	"sync"

	iq "github.com/jackdoe/go-query"
	"github.com/jackdoe/go-query/util/analyzer"
	"github.com/jackdoe/go-query/util/norm"
	"github.com/jackdoe/go-query/util/tokenize"
)

type Document interface {
	IndexableFields() map[string]string
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
	tokenize.NewUnique(),
}

var DefaultIndexTokenizer = []tokenize.Tokenizer{
	tokenize.NewWhitespace(),
	tokenize.NewUnique(),
}

var DefaultAnalyzer = analyzer.NewAnalyzer(DefaultNormalizer, DefaultSearchTokenizer, DefaultIndexTokenizer)

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
			tokens := analyzer.AnalyzeIndex(value)
			for _, t := range tokens {
				m.add(field, t, int32(did))
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
	pk[v] = append(pk[v], did)
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

// Handy method that just ORs all analyzed terms for this field
func (m *MemOnlyIndex) Or(field string, term string) iq.Query {
	queries := m.Terms(field, term)
	if len(queries) == 1 {
		return queries[0]
	}
	return iq.Or(queries...)
}

// Handy method that just ANDs all analyzed terms for this field
func (m *MemOnlyIndex) And(field string, term string) iq.Query {
	queries := m.Terms(field, term)
	if len(queries) == 1 {
		return queries[0]
	}
	return iq.And(queries...)
}

func (m *MemOnlyIndex) newTermQuery(field string, term string) iq.Query {
	s := fmt.Sprintf("%s:%s", field, term)
	pk, ok := m.postings[field]
	if !ok {
		return iq.Term(s, []int32{})
	}
	pv, ok := pk[term]
	if !ok {
		return iq.Term(s, []int32{})
	}
	// there are allocation in iq.Term(), so dont just defer unlock, otherwise it will be locked while term is created
	return iq.Term(s, pv)
}

// Foreach matching document
// Example:
//	query := m.Or("name", "amster")
//	m.Foreach(query, func(did int32, score float32, doc index.Document) {
//		city := doc.(*ExampleCity)
//		log.Printf("%v matching with score %f", city, score)
//	})
func (m *MemOnlyIndex) Foreach(query iq.Query, cb func(int32, float32, Document)) {
	for query.Next() != iq.NO_MORE {
		did := query.GetDocId()
		score := query.Score()

		cb(did, score, m.forward[did])
	}
}
