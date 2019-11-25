package main

import (
	"bufio"
	"io"
	"os"
	"strings"

	"github.com/RoaringBitmap/roaring"
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
func DoIndex(fn string) (map[string][]int32, map[string]*roaring.Bitmap) {
	f, err := os.Open(fn)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	r := bufio.NewReader(f)
	i32 := map[string][]int32{}
	ir := map[string]*roaring.Bitmap{}
	lineNo := int32(0)
	for {
		data, err := r.ReadBytes('\n')
		if err == io.EOF {
			break
		}
		if len(data) == 0 {
			continue
		}
		for _, word := range unique(strings.Split(string(data), " ")) {
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
