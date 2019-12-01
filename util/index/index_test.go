package index

import (
	"log"
	"testing"

	iq "github.com/jackdoe/go-query"
)

// get full list from https://raw.githubusercontent.com/lutangar/cities.json/master/cities.json

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
