package query

import (
	"fmt"
	"math"
	"sort"
)

type block struct {
	maxDoc int32
	maxIdx int
}

type termQuery struct {
	docId             int32
	cursor            int
	postings          []int32
	blocks            []block
	currentBlockIndex int
	currentBlock      block
	term              string
	idf               float32 // XXX: unnormalized idf
	boost             float32
}

func computeIDF(N, d int) float32 {
	// idf is log(1 + N/D)
	// N = total documents in the index
	// d = documents matching (len(postings))
	return float32(math.Log1p(float64(N) / float64(d)))
}

// splits the postings list into chunks that are binary searched and inside each chunk linearly searching for next advance()
var TERM_CHUNK_SIZE = 4096

// Basic []int32{} that the whole interface works on top
// score is IDF (not tf*idf, just 1*idf, since we dont store the term frequency for now)
// if you dont know totalDocumentsInIndex, which could be the case sometimes, pass any constant > 0
// WARNING: the query *can not* be reused
// WARNING: the query it not thread safe
func Term(totalDocumentsInIndex int, t string, postings []int32) *termQuery {
	q := &termQuery{
		term:         t,
		cursor:       -1,
		postings:     postings,
		docId:        NOT_READY,
		currentBlock: block{maxIdx: 0, maxDoc: NOT_READY},
		idf:          computeIDF(totalDocumentsInIndex, len(postings)),
		boost:        1,
	}
	if len(postings) == 0 {
		q.idf = 0
		return q
	}

	chunkSize := TERM_CHUNK_SIZE

	q.blocks = make([]block, ((len(postings) + chunkSize - 1) / chunkSize)) // ceil
	blockIndex := 0

	for i := 0; i < len(postings); i += chunkSize {
		minIdx := i
		maxIdx := (minIdx + chunkSize) - 1
		if maxIdx >= len(postings)-1 {
			maxIdx = len(postings) - 1
		}
		q.blocks[blockIndex] = block{
			maxDoc: postings[maxIdx],
			maxIdx: maxIdx,
		}
		blockIndex++
	}

	return q
}

func (t *termQuery) GetDocId() int32 {
	return t.docId
}

func (t *termQuery) Cost() int {
	return len(t.postings) - t.cursor
}

func (t *termQuery) String() string {
	return fmt.Sprintf("%s/%.2f", t.term, t.idf)
}

func (t *termQuery) Score() float32 {
	return t.idf * t.boost
}

func (t *termQuery) findBlock(target int32) int32 {
	if len(t.blocks) == 0 {
		return NO_MORE
	}

	if len(t.blocks)-t.currentBlockIndex < 32 {
		for i := t.currentBlockIndex; i < len(t.blocks); i++ {
			current := t.blocks[i]
			if target <= current.maxDoc {
				t.currentBlockIndex = i
				t.currentBlock = current
				return target
			}
		}
		return NO_MORE
	}

	found := sort.Search(len(t.blocks)-t.currentBlockIndex, func(i int) bool {
		current := t.blocks[i+t.currentBlockIndex]
		return target <= current.maxDoc
	}) + t.currentBlockIndex

	if found < len(t.blocks) {
		t.currentBlockIndex = found
		t.currentBlock = t.blocks[found]
		return target
	}
	return NO_MORE
}

func (t *termQuery) Advance(target int32) int32 {
	if target > t.currentBlock.maxDoc {
		if t.findBlock(target) == NO_MORE {
			t.docId = NO_MORE
			return NO_MORE
		}
	}

	if t.cursor < 0 {
		t.cursor = 0
	}

	t.docId = NO_MORE

	for i := t.cursor; i <= t.currentBlock.maxIdx; i++ {
		x := t.postings[i]
		if x >= target {
			t.cursor = i
			t.docId = x
			return x
		}
	}
	// invariant, can not happen because at this point we will be within the block
	// panic?
	return t.docId
}

func (t *termQuery) Next() int32 {
	t.cursor++
	if t.cursor >= len(t.postings) {
		t.docId = NO_MORE
	} else {
		t.docId = t.postings[t.cursor]
	}
	return t.docId
}

func (t *termQuery) SetBoost(b float32) Query {
	t.boost = b
	return t
}
