package index

import (
	"log"
	"testing"
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
	m.Foreach(m.Or("name", "aMS u"), func(did int32, score float32, doc Document) {
		city := doc.(*ExampleCity)
		log.Printf("%v matching with score %f", city, score)
		n++
	})
	if n != 2 {
		t.Fatalf("expected 2 got %d", n)
	}
}
