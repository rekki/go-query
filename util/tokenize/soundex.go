package tokenize

import (
	"unicode"
)

const BASE_SOUNDEX = "0000"

var soundexLookup = map[rune]rune{
	'B': '1', 'F': '1', 'P': '1', 'V': '1',
	'C': '2', 'G': '2', 'J': '2', 'K': '2', 'Q': '2', 'S': '2', 'X': '2', 'Z': '2',
	'D': '3', 'T': '3',
	'L': '4',
	'M': '5', 'N': '5',
	'R': '6',
}

func EncodeSoundex(word string) string {
	if len(word) == 0 {
		return BASE_SOUNDEX
	}

	out := make([]rune, 4)
	pos := 0
	prev := rune(0)
	for i, r := range word {
		upper := unicode.ToUpper(r)
		if i == 0 {
			out[0] = upper
			prev = r
			pos++
		} else {
			letterCode, ok := soundexLookup[upper]
			if ok && upper != prev {
				out[pos] = letterCode
				pos++
				if pos == 4 {
					break
				}
				prev = upper
			}
		}
	}

	for i := 0; i < len(out); i++ {
		if out[i] == 0 {
			out[i] = '0'
		}
	}
	return string(out)
}

type Soundex struct {
}

func NewSoundex() *Soundex {
	return &Soundex{}
}

func (w *Soundex) Apply(current []string) []string {
	if len(current) == 0 {
		return current
	}
	out := make([]string, len(current))
	for i := 0; i < len(current); i++ {
		out[i] = EncodeSoundex(current[i])
	}
	return out
}
