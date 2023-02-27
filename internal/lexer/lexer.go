package lexer

import "github.com/sosukesuzuki/regexpp-go/internal/char_code_utils"

type Lexer struct {
	/*
	 文字コードに関するユーティリティ
	 ユニコードモードかどうかによって実装が異なる
	*/
	cu *char_code_utils.CharCodeUtils

	// 現在見ている位置
	I int

	// 現在見ている文字の width
	w int

	// ソースの文字列
	s string

	/*
	 現在のコードポイント
	 ユニコードモードでない場合はコードユニットの場合もある
	*/
	CP uint
}

func NewLexer(s string, u bool) *Lexer {
	var cu char_code_utils.CharCodeUtils
	if u {
		cu = &char_code_utils.Unicode{}
	} else {
		cu = &char_code_utils.Legacy{}
	}

	i := 0
	cp := cu.At(s, i)
	w := cu.Width(cp)
	return &Lexer{
		cu: &cu,
		I:  i,
		w:  w,
		s:  s,
		CP: cp,
	}
}

func (t *Lexer) Next() {
	t.I = t.I + t.w
	t.CP = (*t.cu).At(t.s, t.I)
	t.w = (*t.cu).Width(t.CP)
}

func (t *Lexer) Eat(c uint) bool {
	if t.CP == c {
		t.Next()
		return true
	}
	return false
}

func (t *Lexer) Match(c uint) bool {
	return t.CP == c
}
