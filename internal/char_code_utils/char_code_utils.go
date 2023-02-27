package char_code_utils

import (
	"unicode/utf16"
)

type CharCodeUtils interface {
	At(s string, i int) uint
	Width(c uint) int
}

type Legacy struct{}

func (u *Legacy) At(s string, i int) uint {
	return uint(utf16.Encode([]rune(s))[i])
}
func (u *Legacy) Width(c uint) int {
	return 1
}

type Unicode struct{}

func (u *Unicode) At(s string, i int) uint {
	return uint([]rune(s)[i])
}
func (u *Unicode) Width(c uint) int {
	if c > 0xffff {
		return 2
	} else {
		return 1
	}
}
