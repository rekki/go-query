package query

import (
	"encoding/binary"
	"os"
)

var ByteOrder = binary.LittleEndian

type FileTermData struct {
	cursor   int32
	postings *os.File
	n        int32
	docId    int32
	closed   bool
	boost    float32
	idf      float32
}

// Create new lazy term from stored ByteOrder (by default little
// endian) encoded array of integers
//
// The file will be closed automatically when the query is exhausted (reaches the end)
//
// WARNING: you must exhaust the query, otherwise you will leak file descriptors.
func FileTerm(totalDocumentsInIndex int, fn string) *FileTermData {
	file, err := os.OpenFile(fn, os.O_RDONLY, 0600)
	if err != nil {
		if os.IsNotExist(err) {
			return &FileTermData{
				cursor:   0,
				postings: nil,
				n:        0,
				docId:    NO_MORE,
				boost:    1,
				idf:      0,
				closed:   true,
			}
		}
		panic(err)
	}

	s, err := file.Stat()
	if err != nil {
		panic(err)
	}

	n := int32(s.Size() / 4)
	return &FileTermData{
		cursor:   0,
		postings: file,
		n:        n,
		docId:    NOT_READY,
		boost:    1,
		idf:      computeIDF(totalDocumentsInIndex, int(n)),
	}
}

func (t *FileTermData) GetDocId() int32 {
	return t.docId
}

func (t *FileTermData) SetBoost(b float32) Query {
	t.boost = b
	return t
}

func (t *FileTermData) Cost() int {
	return int(t.n)
}

func (t *FileTermData) String() string {
	s, err := t.postings.Stat()
	if err != nil {
		panic(err)
	}
	return s.Name()
}

func (t *FileTermData) Score() float32 {
	return t.idf * t.boost
}

func (t *FileTermData) getAt(idx int32) uint32 {
	b := []byte{0, 0, 0, 0}
	_, err := t.postings.ReadAt(b, int64(idx*4))
	if err != nil {
		panic(err)
	}
	return ByteOrder.Uint32(b)
}
func (t *FileTermData) Close() {
	if !t.closed {
		t.postings.Close()
		t.closed = true
	}
}
func (t *FileTermData) Advance(target int32) int32 {
	if t.docId == NO_MORE || t.docId == target || target == NO_MORE {
		t.docId = target
		t.Close()
		return t.docId
	}
	start := t.cursor
	end := t.n
	for start < end {
		mid := start + ((end - start) / 2)
		current := int32(t.getAt(mid))
		if current == target {
			t.cursor = mid
			t.docId = target
			return t.GetDocId()
		}

		if current < target {
			start = mid + 1
		} else {
			end = mid
		}
	}

	return t.move(start)
}

func (t *FileTermData) move(to int32) int32 {
	t.cursor = to
	if t.cursor >= t.n {
		t.Close()
		t.docId = NO_MORE
	} else {
		t.docId = int32(t.getAt(t.cursor))
	}
	return t.docId
}

func (t *FileTermData) Next() int32 {
	if t.docId != NOT_READY {
		t.cursor++
	}
	return t.move(t.cursor)
}

func (t *FileTermData) PayloadDecode(p Payload) {
	panic("unsupported")
}

func AppendFileNameTerm(fn string, docs []int32) error {
	f, err := os.OpenFile(fn, os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		return err
	}
	defer f.Close()
	return AppendFileTerm(f, docs)
}

func AppendFileTerm(f *os.File, docs []int32) error {
	b := make([]byte, 4*len(docs))
	for i, did := range docs {
		binary.LittleEndian.PutUint32(b[i*4:], uint32(did))
	}

	return AppendFilePayload(f, 4, b)
}

func AppendFilePayload(f *os.File, size int64, b []byte) error {
	off, err := f.Seek(0, os.SEEK_END)
	if err != nil {
		return err
	}

	// write at closest multiple of 4
	_, err = f.WriteAt(b, (off/size)*size)
	return err
}
