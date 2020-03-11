package index

import (
	"fmt"
	"sync"

	iq "github.com/rekki/go-query"
	"github.com/rekki/go-query/util/analyzer"
	spec "github.com/rekki/go-query/util/go_query_dsl"
)

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
