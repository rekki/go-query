package query

import "math"

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
}

// Basic []int32{} that the whole interface works on top
// score is unnormalized IDF, there is no term frequency
func Term(t string, postings []int32) *termQuery {
	q := &termQuery{
		term:         t,
		cursor:       -1,
		postings:     postings,
		docId:        NOT_READY,
		currentBlock: block{maxIdx: 0, maxDoc: NOT_READY},
		idf:          float32(1) / float32(math.Log1p(float64(len(postings)))),
	}
	if len(postings) == 0 {
		q.idf = 0
		return q
	}

	chunkSize := 4096

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

func (t *termQuery) cost() int {
	return len(t.postings) - t.cursor
}

func (t *termQuery) String() string {
	return t.term
}

func (t *termQuery) Score() float32 {
	return t.idf
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

	start := t.currentBlockIndex
	end := len(t.blocks)

	for start < end {
		mid := start + ((end - start) >> 1)
		current := t.blocks[mid]
		if target <= current.maxDoc {
			t.currentBlockIndex = mid
			t.currentBlock = current
			return target
		}
		if current.maxDoc > target {
			start = mid - 1
		} else {
			end = mid
		}
	}
	return NO_MORE
}

func (t *termQuery) advance(target int32) int32 {
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
