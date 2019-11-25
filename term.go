package query

import (
	"fmt"
	"sort"
)

type block struct {
	minDoc int32
	maxDoc int32
	minIdx int
	maxIdx int
}

func (b block) String() string {
	return fmt.Sprintf("{min doc: %v, max doc: %v}", b.minDoc, b.maxDoc)
}

type termQuery struct {
	docId             int32
	cursor            int
	postings          []int32
	blocks            []block
	currentBlockIndex int
	currentBlock      block
	term              string
}

// Basic []int32{} that the whole interface works on top
func Term(t string, postings []int32) *termQuery {
	q := &termQuery{
		term:     t,
		cursor:   -1,
		postings: postings,
		docId:    NOT_READY,
	}
	if len(postings) == 0 {
		return q
	}

	chunkSize := 1024
	q.blocks = make([]block, ((len(postings) + chunkSize - 1) / chunkSize)) // ceil

	blockIndex := 0
	for i := 0; i < len(postings); i += chunkSize {
		minIdx := i
		maxIdx := (minIdx + chunkSize) - 1
		if maxIdx >= len(postings)-1 {
			maxIdx = len(postings) - 1
		}
		q.blocks[blockIndex] = block{
			minDoc: postings[minIdx],
			maxDoc: postings[maxIdx],
			minIdx: minIdx,
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
	return float32(1)
}

func (t *termQuery) findBlock(target int32) int32 {
	if len(t.blocks) == 0 {
		return NO_MORE
	}

	found := sort.Search(len(t.blocks)-t.currentBlockIndex, func(i int) bool {
		current := t.blocks[i+t.currentBlockIndex]
		if target <= current.maxDoc {
			return true
		}
		return false
	}) + t.currentBlockIndex

	if found < len(t.blocks) {
		t.currentBlockIndex = found
		t.currentBlock = t.blocks[found]
		return target
	}
	return NO_MORE
}

func (t *termQuery) advance(target int32) int32 {
	if t.docId == NO_MORE || t.docId == target || target == NO_MORE {
		t.docId = target
		return t.docId
	}

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
