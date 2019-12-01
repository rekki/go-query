package index

import (
	"bytes"
	"encoding/json"
	fmt "fmt"

	"github.com/gogo/protobuf/jsonpb"
	iq "github.com/jackdoe/go-query"
	"github.com/jackdoe/go-query/util/spec"
)

// somewhat useless method (besides for testing)
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
func QueryFromBytes(b []byte) (*spec.Query, error) {
	q := &spec.Query{}
	err := jsonpb.Unmarshal(bytes.NewReader(b), q)
	if err != nil {
		return nil, err
	}
	return q, nil
}

// simple (*slow*) helper method that takes interface{} and converst it to spec.Query with jsonpb
// in case you receive request like request = {"limit":10, query: ....}, pass request.query to QueryFromJson and get a query object back
func QueryFromJson(input interface{}) (*spec.Query, error) {
	b, err := json.Marshal(input)
	if err != nil {
		return nil, err
	}
	q := &spec.Query{}
	err = jsonpb.Unmarshal(bytes.NewReader(b), q)
	if err != nil {
		return nil, err
	}

	return q, nil
}

// Take spec.Query object and a makeTermQuery function and produce a parsed query
// Example:
//  return Parse(input, func(k, v string) iq.Query {
//  	kv := k + ":"+ v
//  	return iq.Term(0, kv, postings[kv])
//  })
func Parse(input *spec.Query, makeTermQuery func(string, string) iq.Query) (iq.Query, error) {
	if input == nil {
		return nil, fmt.Errorf("nil input")
	}
	if input.Type == spec.Query_TERM {
		if input.Not != nil || len(input.Queries) != 0 {
			return nil, fmt.Errorf("term queries can have only field and value, %v", input)
		}
		if input.Field == "" {
			return nil, fmt.Errorf("missing field, %v", input)
		}
		t := makeTermQuery(input.Field, input.Value)
		if input.Boost > 0 {
			t.SetBoost(input.Boost)
		}
		return t, nil
	}
	queries := []iq.Query{}
	for _, q := range input.Queries {
		p, err := Parse(q, makeTermQuery)
		if err != nil {
			return nil, err
		}
		queries = append(queries, p)
	}

	if input.Type == spec.Query_AND {
		and := iq.And(queries...)
		if input.Not != nil {
			p, err := Parse(input.Not, makeTermQuery)
			if err != nil {
				return nil, err
			}
			and.SetNot(p)

		} else {
			if len(queries) == 1 {
				return queries[0], nil
			}
		}
		if input.Boost > 0 {
			and.SetBoost(input.Boost)
		}

		return and, nil
	}

	if input.Type == spec.Query_OR {
		if input.Not != nil {
			return nil, fmt.Errorf("or queries cant have 'not' value, %v", input)
		}
		if len(queries) == 1 {
			return queries[0], nil
		}
		or := iq.Or(queries...)
		if input.Boost > 0 {
			or.SetBoost(input.Boost)
		}

		return or, nil
	}

	if input.Type == spec.Query_DISMAX {
		if input.Not != nil {
			return nil, fmt.Errorf("or queries cant have 'not' value, %v", input)
		}
		if len(queries) == 1 {
			return queries[0], nil
		}

		d := iq.DisMax(input.Tiebreaker, queries...)
		if input.Boost > 0 {
			d.SetBoost(input.Boost)
		}

		return d, nil
	}

	return nil, fmt.Errorf("unknown type %v", input)
}
