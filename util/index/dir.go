package index

import (
	"encoding/binary"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"sync"

	iq "github.com/rekki/go-query"
	"github.com/rekki/go-query/util/analyzer"
	"github.com/rekki/go-query/util/common"
	spec "github.com/rekki/go-query/util/go_query_dsl"
)

type FDCache struct {
	fdCache   map[string]*os.File
	maxOpenFD int
	sync.RWMutex
}

func NewFDCache(n int) *FDCache {
	return &FDCache{maxOpenFD: n, fdCache: map[string]*os.File{}}
}

func (x *FDCache) Close() {
	x.Lock()
	defer x.Unlock()

	for _, fd := range x.fdCache {
		_ = fd.Close()
	}
}

func (x *FDCache) Use(fn string, createFile func(fn string) (*os.File, error), cb func(*os.File) error) error {
	var err error
	var ok bool
	var f *os.File

	x.RLock()
	f, ok = x.fdCache[fn]
	if !ok {
		x.RUnlock()

		_ = os.MkdirAll(path.Dir(fn), 0700)

		f, err = createFile(fn)

		if err != nil {
			return err
		}

		x.Lock()

		overriden, ok := x.fdCache[fn]
		if ok {
			f.Close()
			f = overriden
		} else {
			if len(x.fdCache) > x.maxOpenFD {
				for _, fd := range x.fdCache {
					_ = fd.Close()
				}
				x.fdCache = map[string]*os.File{}
			}
			x.fdCache[fn] = f
		}

		err = cb(f)

		x.Unlock()
		return err
	}

	err = cb(f)

	x.RUnlock()
	return err
}

type FileDescriptorCache interface {
	Use(fn string, createFile func(fn string) (*os.File, error), cb func(*os.File) error) error
	Close()
}

type DirIndex struct {
	perField          map[string]*analyzer.Analyzer
	root              string
	fdCache           FileDescriptorCache
	TotalNumberOfDocs int
	Lazy              bool
	DirHash           func(s string) string
}

func NewDirIndex(root string, fdCache FileDescriptorCache, perField map[string]*analyzer.Analyzer) *DirIndex {
	if perField == nil {
		perField = map[string]*analyzer.Analyzer{}
	}

	dh := func(s string) string {
		return string(s[len(s)-1])
	}
	return &DirIndex{TotalNumberOfDocs: 1, root: root, fdCache: fdCache, perField: perField, DirHash: dh}
}

var DirIndexMaxTermLen = 64

func termCleanup(s string) string {
	x := common.ReplaceNonAlphanumericWith(s, '_')
	if len(x) > DirIndexMaxTermLen {
		return x[:DirIndexMaxTermLen]
	}
	return x
}

func (d *DirIndex) add(fn string, docs []int32) error {
	err := d.fdCache.Use(
		fn,
		func(_s string) (*os.File, error) {
			return os.OpenFile(fn, os.O_CREATE|os.O_WRONLY, 0600)
		}, func(f *os.File) error {
			return iq.AppendFileTerm(f, docs)
		})
	return err
}

type DocumentWithID interface {
	IndexableFields() map[string][]string
	DocumentID() int32
}

func (d *DirIndex) Index(docs ...DocumentWithID) error {
	var sb strings.Builder

	todo := map[string][]int32{}

	for _, doc := range docs {
		did := doc.DocumentID()

		fields := doc.IndexableFields()
		for field, value := range fields {
			field = termCleanup(field)
			if len(field) == 0 {
				continue
			}

			analyzer, ok := d.perField[field]
			if !ok {
				analyzer = DefaultAnalyzer
			}
			for _, v := range value {
				tokens := analyzer.AnalyzeIndex(v)
				for _, t := range tokens {
					t = termCleanup(t)
					if len(t) == 0 {
						continue
					}

					sb.WriteString(d.root)
					sb.WriteRune('/')
					sb.WriteString(field)
					sb.WriteRune('/')
					sb.WriteString(d.DirHash(t))
					sb.WriteRune('/')
					sb.WriteString(t)

					s := sb.String()
					todo[s] = append(todo[s], did)
					sb.Reset()
				}
			}
		}
	}

	for t, docs := range todo {
		err := d.add(t, docs)
		if err != nil {
			return err
		}
	}

	return nil
}

func (d *DirIndex) Parse(input *spec.Query) (iq.Query, error) {
	return Parse(input, func(k, v string) iq.Query {
		terms := d.Terms(k, v)
		if len(terms) == 1 {
			return terms[0]
		}
		return iq.Or(terms...)
	})
}

func (d *DirIndex) Terms(field string, term string) []iq.Query {
	analyzer, ok := d.perField[field]
	if !ok {
		analyzer = DefaultAnalyzer
	}
	tokens := analyzer.AnalyzeSearch(term)
	queries := []iq.Query{}
	for _, t := range tokens {
		queries = append(queries, d.newTermQuery(field, t))
	}
	return queries
}

func (d *DirIndex) newTermQuery(field string, term string) iq.Query {
	field = termCleanup(field)
	term = termCleanup(term)
	if len(field) == 0 || len(term) == 0 {
		return iq.Term(d.TotalNumberOfDocs, fmt.Sprintf("broken(%s:%s)", field, term), []int32{})
	}
	fn := path.Join(d.root, field, d.DirHash(term), term)

	if d.Lazy {
		return iq.FileTerm(d.TotalNumberOfDocs, fn)
	} else {
		data, err := ioutil.ReadFile(fn)
		if err != nil {
			return iq.Term(d.TotalNumberOfDocs, fn, []int32{})
		}
		postings := make([]int32, len(data)/4)
		for i := 0; i < len(postings); i++ {
			from := i * 4
			postings[i] = int32(binary.LittleEndian.Uint32(data[from : from+4]))
		}
		return iq.Term(d.TotalNumberOfDocs, fn, postings)
	}
}

func (d *DirIndex) Close() {
	d.fdCache.Close()
}

func (d *DirIndex) Foreach(query iq.Query, cb func(int32, float32)) {
	for query.Next() != iq.NO_MORE {
		did := query.GetDocId()
		score := query.Score()

		cb(did, score)
	}
}
