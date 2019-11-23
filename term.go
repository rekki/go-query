package query

type TermQuery struct {
	docId    int32
	cursor   int
	postings []int32
	term     string
}

func Term(t string, postings []int32) *TermQuery {
	return &TermQuery{
		term:     t,
		cursor:   -1,
		postings: postings,
		docId:    NOT_READY,
	}
}

func (t *TermQuery) GetDocId() int32 {
	return t.docId
}

func (t *TermQuery) String() string {
	return t.term
}

func (t *TermQuery) Score() float32 {
	return float32(1)
}

func (t *TermQuery) advance(target int32) int32 {
	if t.docId == NO_MORE || t.docId == target || target == NO_MORE {
		t.docId = target
		return t.docId
	}
	if t.cursor < 0 {
		t.cursor = 0
	}

	start := t.cursor
	end := len(t.postings)

	for start < end {
		mid := start + ((end - start) >> 1)
		current := t.postings[mid]
		if current == target {
			t.cursor = mid
			t.docId = target
			return target
		}

		if current < target {
			start = mid + 1
		} else {
			end = mid
		}
	}
	if start >= len(t.postings) {
		t.docId = NO_MORE
		return NO_MORE
	}
	t.cursor = start
	t.docId = t.postings[start]
	return t.docId
}

func (t *TermQuery) Next() int32 {
	t.cursor++
	if t.cursor >= len(t.postings) {
		t.docId = NO_MORE
	} else {
		t.docId = t.postings[t.cursor]
	}
	return t.docId
}
