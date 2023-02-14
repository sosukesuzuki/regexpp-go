package regexpp

type Lexer struct {
	/*
	 文字コードに関するユーティリティ
	 ユニコードモードかどうかによって実装が異なる
	*/
	cu *CharCodeUtils

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
	var cu CharCodeUtils
	if u {
		cu = &UnicodeCharUtils{}
	} else {
		cu = &LegacyCharUtils{}
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
