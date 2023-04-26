package char_code_utils

import (
	"unicode/utf16"
)

type CharCodeUtils interface {
	At(s string, i int) int
	Width(c int) int
}

type Legacy struct{}

func (u *Legacy) At(s string, i int) int {
	runes := utf16.Encode([]rune(s))
	if i >= 0 && i < len(runes) {
		return int(utf16.Encode([]rune(s))[i])
	}
	return -1
}
func (u *Legacy) Width(c int) int {
	return 1
}

type Unicode struct{}

func (u *Unicode) At(s string, i int) int {
	runes := []rune(s)
	if i >= 0 && i < len(runes) {
		return int([]rune(s)[i])
	}
	return -1
}
func (u *Unicode) Width(c int) int {
	if c > 0xffff {
		return 2
	} else {
		return 1
	}
}
