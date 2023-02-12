package regexpp

type Tokenizer struct {
	// 文字コードに関するユーティリティ
	// ユニコードモードかどうかによって実装が異なる
	cu *CharCodeUtils
	// 現在見ている位置
	i int
	// ソースの文字列
	s string
	// 現在のコードポイント
	// ユニコードモードでない場合はコードユニットの場合もある
	CP uint
}

func NewTokenizer(s string, u bool) Tokenizer {
	var cu CharCodeUtils
	if u {
		cu = &UnicodeCharUtils{}
	} else {
		cu = &LegacyCharUtils{}
	}
	i := 0
	return Tokenizer{
		cu: &cu,
		i:  i,
		s:  s,
		CP: cu.At(s, i),
	}
}
