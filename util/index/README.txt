package index // import "github.com/jackdoe/go-query/util/index"

Illustration of how you can use go-query to build a somewhat functional
search index Example:

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

    	a := analyzer.NewAnalyzer(index.DefaultNormalizer, searchTokenizer, indexTokenizer)
    	m := index.NewMemOnlyIndex(map[string]*analyzer.Analyzer{
    		"name":    a,
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

    	// search for "(name:aMS OR name:u) AND *country:NL OR country:BG)"

    	query := iq.And(m.Or("name", "aMS u"), m.Or("country", "NL BG"))

    	m.Foreach(query, func(did int32, score float32, doc index.Document) {
    		city := doc.(*ExampleCity)
    		log.Printf("%v matching with score %f", city, score)
    	})
    }

VARIABLES

var DefaultAnalyzer = analyzer.NewAnalyzer(DefaultNormalizer, DefaultSearchTokenizer, DefaultIndexTokenizer)
var DefaultIndexTokenizer = []tokenize.Tokenizer{
	tokenize.NewWhitespace(),
	tokenize.NewUnique(),
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

TYPES

type Document interface {
	IndexableFields() map[string]string
}

type MemOnlyIndex struct {
	sync.RWMutex
	// Has unexported fields.
}

func NewMemOnlyIndex(perField map[string]*analyzer.Analyzer) *MemOnlyIndex
    create new in-memory index with the specified perField analyzer by default
    DefaultAnalyzer is used

func (m *MemOnlyIndex) And(field string, term string) iq.Query
    Handy method that just ANDs all analyzed terms for this field

func (m *MemOnlyIndex) Foreach(query iq.Query, cb func(int32, float32, Document))
    Foreach matching document Example:

        query := m.Or("name", "amster")
        m.Foreach(query, func(did int32, score float32, doc index.Document) {
        	city := doc.(*ExampleCity)
        	log.Printf("%v matching with score %f", city, score)
        })

func (m *MemOnlyIndex) Index(docs ...Document)
    index a bunch of documents

func (m *MemOnlyIndex) Or(field string, term string) iq.Query
    Handy method that just ORs all analyzed terms for this field

func (m *MemOnlyIndex) Terms(field string, term string) []iq.Query
    Generate array of queries from the tokenized term for this field, using the
    perField analyzer

