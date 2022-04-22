package query

import (
	"fmt"
	"sort"
)

type TermTFQuery struct {
	docId             int32
	cursor            int
	postings          []int32
	blocks            []block
	currentBlockIndex int
	currentBlock      block
	term              string
	idf               float32 // XXX: unnormalized idf
	boost             float32
	freqBits          int32
	freqMask          int32
}

// Splits the postings list into chunks that are binary searched and inside each
// chunk linearly searching for next advance() Basic []int32{} that the whole
// interface works on top. The Score is TF*IDF you have to specify how many bits
// from the docID are actually term frequency e.g if you want to store the
// frequency in 4 bits then document id 999 with term frequency 2 for this
// specific term could be stored as (999 << 4) | 2, usually you just store the
// floored sqrt(frequency), so 3-4 bits are enough. it is zero based, so 0 is
// frequency 1
//
// if you dont know totalDocumentsInIndex, which could be the case sometimes, pass any constant > 0
// WARNING: the query *can not* be reused
// WARNING: the query it not thread safe
func TermTF(totalDocumentsInIndex int, freqBits int32, t string, postings []int32) *TermTFQuery {
	q := &TermTFQuery{
		term:         t,
		cursor:       -1,
		postings:     postings,
		docId:        NOT_READY,
		currentBlock: block{maxIdx: 0, maxDoc: NOT_READY},
		idf:          computeIDF(totalDocumentsInIndex, len(postings)),
		boost:        1,
		freqBits:     freqBits,
		freqMask:     (1 << freqBits) - 1,
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
			maxDoc: postings[maxIdx] >> q.freqBits,
			maxIdx: maxIdx,
		}
		blockIndex++
	}
	return q
}

func (t *TermTFQuery) GetDocId() int32 {
	return t.docId
}

func (t *TermTFQuery) Cost() int {
	return len(t.postings) - t.cursor
}

func (t *TermTFQuery) String() string {
	return fmt.Sprintf("%s/%.2f", t.term, t.idf)
}

func (t *TermTFQuery) Score() float32 {
	if t.docId == NO_MORE {
		return 0
	}

	tf := float32(1 + (t.postings[t.cursor] & t.freqMask))
	return tf * t.idf * t.boost
}

func (t *TermTFQuery) findBlock(target int32) int32 {
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

func (t *TermTFQuery) Advance(target int32) int32 {
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
		x := t.postings[i] >> t.freqBits
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

func (t *TermTFQuery) Next() int32 {
	t.cursor++
	if t.cursor >= len(t.postings) {
		t.docId = NO_MORE
	} else {
		t.docId = t.postings[t.cursor] >> t.freqBits

	}
	return t.docId
}

func (t *TermTFQuery) SetBoost(b float32) Query {
	t.boost = b
	return t
}

func (t *TermTFQuery) PayloadDecode(p Payload) {
	panic("unsupported")
}

func (t *TermTFQuery) AddSubQuery(Query) Query {
	panic("unsupported")
}
