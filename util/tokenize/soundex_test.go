package tokenize

import (
	"testing"
)

func TestSoundexEncode(t *testing.T) {
	cases := map[string]string{
		"":            "0000",
		"A":           "A000",
		"AB":          "A100",
		"ABC":         "A120",
		"ABCD":        "A123",
		"Ashcraft":    "A226",
		"Ashcroft":    "A226",
		"Burroughs":   "B622",
		"Burrows":     "B620",
		"Ciondecks":   "C532",
		"Ellery":      "E460",
		"Euler":       "E460",
		"Example":     "E251",
		"Gauss":       "G200",
		"Ghosh":       "G200",
		"Heilbronn":   "H416",
		"Hilbert":     "H416",
		"Kant":        "K530",
		"Knuth":       "K530",
		"Ladd":        "L300",
		"Lissajous":   "L222",
		"Lloyd":       "L300",
		"Lukasiewicz": "L222",
		"O'Hara":      "O600",
		"Robert":      "R163",
		"Rubin":       "R150",
		"Rupert":      "R163",
		"Soundex":     "S532",
		"Tymczak":     "T522",
		"Wheaton":     "W350",
	}
	for example, s := range cases {
		v := EncodeSoundex(example)
		if s != v {
			t.Fatalf("%s, expected: %s got %s", example, s, v)
		}
	}
}
