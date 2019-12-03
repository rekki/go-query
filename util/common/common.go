package common

import (
	"strings"
	"unicode"
)

func Unique(s []string) []string {
	seen := map[string]bool{}
	out := []string{}
	for _, v := range s {
		_, ok := seen[v]
		if !ok {
			out = append(out, v)
			seen[v] = true
		}
	}
	return out
}

func IsASCII(s string) bool {
	for _, c := range s {
		if c > 127 {
			return false
		}
	}
	return true
}

func IsAZ(s string) bool {
	for _, c := range s {
		if !(('a' <= c && c <= 'z') || ('A' <= c && c <= 'Z')) {
			return false
		}
	}
	return true
}

func HasDigit(s string) bool {
	for _, c := range s {
		if unicode.IsDigit(c) {
			return true
		}
	}
	return false
}

func SpaceBetweenDigits(s string) string {
	if !HasDigit(s) {
		return s
	}
	digitMode := false
	var sb strings.Builder

	sb.Grow(len(s) * 2)

	for i, c := range s {
		isDigit := unicode.IsDigit(c) || c == '-'
		if i == 0 {
			digitMode = isDigit
			sb.WriteRune(c)
			continue
		}
		if c != ' ' {
			if isDigit {
				if !digitMode {
					digitMode = true
					if s[i-1] != ' ' {
						sb.WriteRune(' ')
					}
				}
			} else {
				if digitMode {
					digitMode = false
					if s[i-1] != ' ' {
						sb.WriteRune(' ')
					}
				}
			}
		}
		sb.WriteRune(c)
	}

	return sb.String()
}
