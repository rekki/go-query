package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/RoaringBitmap/roaring"
	"github.com/blevesearch/bleve"
)

func unique(s []string) []string {
	seen := map[string]bool{}
	out := []string{}
	for _, v := range s {
		if !seen[v] {
			seen[v] = true
			out = append(out, v)
		}
	}
	return out
}

func DoBleveIndex(fn string) bleve.Index {
	f, err := os.Open(fn)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	mapping := bleve.NewIndexMapping()
	index, err := bleve.NewMemOnly(mapping)
	if err != nil {
		panic(err)
	}

	r := bufio.NewReader(f)
	lineNo := int32(1)
	for {
		data, err := r.ReadBytes('\n')
		if err == io.EOF {
			break
		}
		if len(data) == 0 {
			continue
		}

		err = index.Index(fmt.Sprintf("%d", lineNo), map[string]interface{}{"line": string(data)})
		if err != nil {
			panic(err)
		}

		lineNo++
		if lineNo%1000 == 0 {
			log.Printf("bleve %v ...", lineNo)
		}
	}

	return index
}

func DoIndex(fn string) (map[string][]int32, map[string]*roaring.Bitmap) {
	f, err := os.Open(fn)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	r := bufio.NewReader(f)
	i32 := map[string][]int32{}
	ir := map[string]*roaring.Bitmap{}
	lineNo := int32(1)
	for {
		data, err := r.ReadBytes('\n')
		if err == io.EOF {
			break
		}
		if len(data) == 0 {
			continue
		}
		d := strings.Trim(string(data), "\n ")
		for _, word := range unique(strings.Split(d, " ")) {
			i32[word] = append(i32[word], lineNo)
			r, ok := ir[word]
			if !ok {
				r = roaring.New()
				ir[word] = r
			}
			r.Add(uint32(lineNo))
		}
		lineNo++
	}
	return i32, ir
}

func main() {
}
