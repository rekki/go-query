package index

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"testing"

	iq "github.com/rekki/go-query"
	spec "github.com/rekki/go-query/util/go_query_dsl"
)

// get full list from https://raw.githubusercontent.com/lutangar/cities.json/master/cities.json

type ExampleCity struct {
	ID      int32
	Name    string
	Country string
	Names   []string
}

func (e *ExampleCity) DocumentID() int32 {
	return e.ID
}

func (e *ExampleCity) IndexableFields() map[string][]string {
	out := map[string][]string{}

	out["name"] = []string{e.Name}
	out["names"] = e.Names
	out["country"] = []string{e.Country}

	return out
}

func toDocuments(in []*ExampleCity) []Document {
	out := make([]Document, len(in))
	for i, d := range in {
		out[i] = Document(d)
	}
	return out
}

func toDocumentsID(in []*ExampleCity) []DocumentWithID {
	out := make([]DocumentWithID, len(in))
	for i, d := range in {
		out[i] = DocumentWithID(d)
	}
	return out
}

func TestUnique(t *testing.T) {
	m := NewMemOnlyIndex(nil)
	list := []*ExampleCity{
		{Names: []string{"Amsterdam", "Amsterdam"}, Country: "NL"},
		{Names: []string{"Sofia", "Sofia"}, Country: "NL"},
	}

	m.Index(toDocuments(list)...)
	n := 0
	q := iq.Or(m.Terms("names", "sofia")...)

	m.Foreach(q, func(did int32, score float32, doc Document) {
		n++
	})
	if n != 1 {
		t.Fatalf("expected 2 got %d", n)
	}
}

func TestExample(t *testing.T) {
	m := NewMemOnlyIndex(nil)
	list := []*ExampleCity{
		{Name: "Amsterdam", Country: "NL"},
		{Name: "Amsterdam, USA", Country: "USA"},
		{Name: "London", Country: "UK"},
		{Name: "Sofia", Country: "BG"},
	}

	m.Index(toDocuments(list)...)
	n := 0
	q := iq.Or(m.Terms("name", "aMSterdam sofia")...)

	m.Foreach(q, func(did int32, score float32, doc Document) {
		city := doc.(*ExampleCity)
		log.Printf("%v matching with score %f", city, score)
		n++
	})
	if n != 3 {
		t.Fatalf("expected 2 got %d", n)
	}
	n = 0

	q = iq.Or(m.Terms("name", "aMSterdam sofia")...)
	top := m.TopN(1, q, func(did int32, score float32, doc Document) float32 {
		city := doc.(*ExampleCity)
		if city.Country == "NL" {
			score += 100
		}
		n++
		return score
	})

	if top.Hits[0].Score < 100 {
		t.Fatalf("expected > 100")
	}
	if top.Total != 3 {
		t.Fatalf("expected 3")
	}
	if len(top.Hits) != 1 {
		t.Fatalf("expected 1")
	}

	q = iq.Or(m.Terms("name", "aMSterdam sofia")...)
	top = m.TopN(0, q, func(did int32, score float32, doc Document) float32 {
		return score
	})

	if len(top.Hits) != 0 {
		t.Fatalf("expected 0")
	}
	if top.Total != 3 {
		t.Fatalf("expected 3")
	}
}

func TestExampleDir(t *testing.T) {
	dir, err := ioutil.TempDir("", "forward")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)

	m := NewDirIndex(dir, NewFDCache(10), nil)
	list := []*ExampleCity{
		{Name: "Amsterdam", Country: "NL", ID: 0},
		{Name: "Amsterdam, USA", Country: "USA", ID: 1},
		{Name: "London", Country: "UK", ID: 2},
		{Name: "Sofia Amsterdam", Country: "BG", ID: 3},
	}

	for i := len(list); i < 10000; i++ {
		list = append(list, &ExampleCity{Name: fmt.Sprintf("%dLondon", i), Country: "UK", ID: int32(i)})
	}
	err = m.Index(toDocumentsID(list)...)
	if err != nil {
		t.Fatal(err)
	}
	n := 0
	q := iq.And(m.Terms("name", "aMSterdam sofia")...)

	m.Foreach(q, func(did int32, score float32) {
		city := list[did]
		log.Printf("%v matching with score %f", city, score)
		n++
	})
	if n != 1 {
		t.Fatalf("expected 1 got %d", n)
	}

	n = 0
	qq := iq.Or(m.Terms("name", "aMSterdam sofia")...)

	m.Foreach(qq, func(did int32, score float32) {
		city := list[did]
		log.Printf("%v matching with score %f", city, score)
		n++
	})
	if n != 3 {
		t.Fatalf("expected 3 got %d", n)
	}

	m.Lazy = true

	n = 0
	qqq := iq.Or(m.Terms("name", "aMSterdam sofia")...)

	m.Foreach(qqq, func(did int32, score float32) {
		city := list[did]
		log.Printf("lazy %v matching with score %f", city, score)
		n++
	})
	if n != 3 {
		t.Fatalf("expected 3 got %d", n)
	}

}

func BenchmarkDirIndexBuild(b *testing.B) {
	b.StopTimer()
	dir, err := ioutil.TempDir("", "forward")
	if err != nil {
		panic(err)
	}
	defer os.RemoveAll(dir)

	m := NewDirIndex(dir, NewFDCache(10), nil)
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		err = m.Index(DocumentWithID(&ExampleCity{Name: "Amsterdam", Country: "NL", ID: int32(i)}))
		if err != nil {
			panic(err)
		}
	}
	b.StopTimer()

}

func BenchmarkMemIndexBuild(b *testing.B) {
	m := NewMemOnlyIndex(nil)
	for i := 0; i < b.N; i++ {
		m.Index(DocumentWithID(&ExampleCity{Name: "Amsterdam", Country: "NL", ID: int32(i)}))
	}

}

var dont = 0

func BenchmarkDirIndexSearch10000(b *testing.B) {
	b.StopTimer()
	dir, err := ioutil.TempDir("", "forward")
	if err != nil {
		panic(err)
	}
	defer os.RemoveAll(dir)
	m := NewDirIndex(dir, NewFDCache(10), nil)
	for i := 0; i < 10000; i++ {
		err = m.Index(DocumentWithID(&ExampleCity{Name: "Amsterdam", Country: "NL", ID: int32(i)}))
		if err != nil {
			panic(err)
		}
	}

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		n := 0
		q := iq.Or(m.Terms("name", "aMSterdam sofia")...)
		m.Foreach(q, func(did int32, score float32) {
			n++
			dont++

		})
	}
	b.StopTimer()
}

func BenchmarkMemIndexSearch10000(b *testing.B) {
	b.StopTimer()
	m := NewMemOnlyIndex(nil)
	for i := 0; i < 10000; i++ {
		m.Index(Document(&ExampleCity{Name: "Amsterdam", Country: "NL", ID: int32(i)}))
	}

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		n := 0
		q := iq.Or(m.Terms("name", "aMSterdam sofia")...)
		m.Foreach(q, func(did int32, score float32, _d Document) {
			n++
			dont++

		})
	}
	b.StopTimer()
}

func TestParsing(t *testing.T) {
	m := NewMemOnlyIndex(nil)
	list := []*ExampleCity{
		{Name: "Amsterdam", Country: "NL"},
		{Name: "Amsterdam, USA", Country: "USA"},
		{Name: "London", Country: "UK"},
		{Name: "Sofia", Country: "BG"},
	}

	m.Index(toDocuments(list)...)

	query, err := QueryFromBytes([]byte(`{"type":"OR", "queries":[{"field":"name","value":"sofia"}, {"field":"name","value":"amsterdam"}]}`))
	if err != nil {
		panic(err)
	}
	parsedQuery, _ := m.Parse(query)
	top := m.TopN(1, parsedQuery, nil)
	if top.Total != 3 {
		t.Fatalf("expected 2 got %d", top.Total)
	}
	if len(top.Hits) != 1 {
		t.Fatalf("expected 1")
	}

	parsedQuery, _ = m.Parse(toQuery(`{"type":"OR", "queries":[{"field":"name","value":"sofia"}, {"field":"name","value":"amsterdam"}]}`))
	top = m.TopN(1, parsedQuery, nil)
	if top.Total != 3 {
		t.Fatalf("expected 2 got %d", top.Total)
	}
	if len(top.Hits) != 1 {
		t.Fatalf("expected 1")
	}

	_, err = m.Parse(nil)
	if err == nil {
		t.Fatalf("error")
	}
	_, err = m.Parse(toQuery(`{}`))
	if !strings.Contains(err.Error(), "field") {
		t.Fatalf("need field error")
	}

	_, err = QueryFromBytes([]byte(`{xx:wrong}`))
	if err == nil {
		t.Fatalf("need json error")
	}

	_, err = QueryFromJson(nil)
	if err != nil {
		t.Fatalf("need field error")
	}

	_, err = m.Parse(toQuery(`{"value": "abc"}`))
	if !strings.Contains(err.Error(), "field") {
		t.Fatalf("need field error")
	}

	_, err = m.Parse(toQuery(`{"value": "abc", "field":"abc", "not": {"field":"x","value":"y"}}`))
	if !strings.Contains(err.Error(), "term queries can have only field and value") {
		t.Fatalf("need field error")
	}

	_, err = m.Parse(toQuery(`{"field": "abc", "value": "abc"}`))
	if err != nil {
		t.Fatalf("should be without error")
	}

	_, err = m.Parse(toQuery(`{"field": "abc", "value": "abc"}`))
	if err != nil {
		t.Fatalf("should be without error")
	}

	parsedQuery, _ = m.Parse(toQuery(`{"field": "name", "value": "sofia", "boost": 100000}`))
	topScore := m.TopN(1, parsedQuery, nil).Hits[0].Score
	if topScore < 100 {
		t.Fatalf("should be boosted")
	}

	parsedQuery, _ = m.Parse(toQuery(`{"field": "name", "value": "sofia", "boost": 100000}`))
	topScore = m.TopN(1, parsedQuery, nil).Hits[0].Score
	if topScore < 100 {
		t.Fatalf("should be boosted")
	}

	parsedQuery, _ = m.Parse(toQuery(`{
  "type": "OR",
  "queries": [
    {
      "type": "AND",
      "queries": [
        {
          "field": "name",
          "value": "sofia",
          "boost": 100000
        }
      ],
      "boost": 50
    },
    {
      "type": "OR",
      "queries": [
        {
          "field": "name",
          "value": "sofia",
          "boost": 100000
        }
      ],
      "boost": 50
    },
    {
      "type": "DISMAX",
      "queries": [
        {
          "field": "name",
          "value": "sofia",
          "boost": 100000
        },
        {
          "field": "name",
          "value": "sofia usa",
          "boost": 100000
        }
      ],
      "boost": 50
    },
    {
      "type": "DISMAX",
      "queries": [
        {
          "field": "name",
          "value": "sofia",
          "boost": 100000
        },
        {
          "field": "name",
          "value": "sofia usa",
          "boost": 100000
        },
        {
          "field": "name",
          "value": "zzz",
          "boost": 100000
        }
      ],
      "boost": 50
    },
    {
      "type": "DISMAX",
      "queries": [
        {
          "field": "name",
          "value": "sofia",
          "boost": 100000
        }
      ],
      "boost": 50
    },
    {
      "type": "AND",
      "not": {
        "field": "name",
        "value": "amsterdam",
        "boost": 100000
      },
      "queries": [
        {
          "field": "name",
          "value": "sofia",
          "boost": 100000
        }
      ],
      "boost": 50
    }
  ],
  "boost": 100
}`))
	topScore = m.TopN(1, parsedQuery, nil).Hits[0].Score
	if topScore < 100 {
		t.Fatalf("should be boosted")
	}

}
func toQuery(s string) *spec.Query {
	var unparsed interface{}
	err := json.Unmarshal([]byte(s), &unparsed)
	if err != nil {
		panic(err)
	}
	query, err := QueryFromJson(unparsed)
	if err != nil {
		panic(err)
	}

	return query
}
