package index

import (
	"encoding/json"
	"log"
	"strings"
	"testing"

	iq "github.com/jackdoe/go-query"
	spec "github.com/jackdoe/go-query/util/go_query_dsl"
)

// get full list from https://raw.githubusercontent.com/lutangar/cities.json/master/cities.json

type ExampleCity struct {
	Name    string
	Country string
}

func (e *ExampleCity) IndexableFields() map[string][]string {
	out := map[string][]string{}

	out["name"] = []string{e.Name}
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

func TestExample(t *testing.T) {
	m := NewMemOnlyIndex(nil)
	list := []*ExampleCity{
		&ExampleCity{Name: "Amsterdam", Country: "NL"},
		&ExampleCity{Name: "Amsterdam, USA", Country: "USA"},
		&ExampleCity{Name: "London", Country: "UK"},
		&ExampleCity{Name: "Sofia", Country: "BG"},
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
}

func TestParsing(t *testing.T) {
	m := NewMemOnlyIndex(nil)
	list := []*ExampleCity{
		&ExampleCity{Name: "Amsterdam", Country: "NL"},
		&ExampleCity{Name: "Amsterdam, USA", Country: "USA"},
		&ExampleCity{Name: "London", Country: "UK"},
		&ExampleCity{Name: "Sofia", Country: "BG"},
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
