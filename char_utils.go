package regexpp

import (
	"unicode/utf16"
)

type CharUtils interface {
	At(s string, i int) uint
	Width(c uint) int
}

type LegacyCharUtils struct{}

func (u *LegacyCharUtils) At(s string, i int) uint {
	return uint(utf16.Encode([]rune(s))[i])
}
func (u *LegacyCharUtils) Width(c uint) int {
	return 1
}

type UnicodeCharUtils struct{}

func (u *UnicodeCharUtils) At(s string, i int) uint {
	return uint([]rune(s)[i])
}
func (u *UnicodeCharUtils) Width(c uint) int {
	if c > 0xffff {
		return 2
	} else {
		return 1
	}
}
