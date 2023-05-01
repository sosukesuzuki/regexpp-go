package unicode_consts

const (
	Eof                 = 0x1A
	Backspace = 0x08
	LatinSmallLetterA   = 0x61 // a
	LatinSmallLetterB = 0x62 // b
	LatinSmallLetterF   = 0x66 // f
	LatinCapitalLetterA = 0x41 // A
	LatinCapitalLetterF = 0x46 // F
	DigitZero           = 0x30 // 0
	DigitNine           = 0x39 // 9
	VerticalLine        = 0x7c // |
	CircumflexAccent    = 0x5e // ^
	DollarSign          = 0x24 // $
	ReverseSolidus      = 0x5c // \
	FullStop            = 0x2e // .
	Asterisk            = 0x2a // *
	PlusSign            = 0x2b // +
	QuestionMark        = 0x3f // ?
	LeftParenthesis     = 0x28 // (
	RightParenthesis    = 0x29 // )
	LeftSquareBracket   = 0x5b // [
	RightSquareBracket  = 0x5d // ]
	LeftCurlyBracket    = 0x7b // {
	RightCurlyBracket   = 0x7d // {
	Comma               = 0x2c // ,
	HyphenMinus = 0x2d // -
)

func IsDecimalDigit(code int) bool {
	return code >= DigitZero && code <= DigitNine
}

func DecimalToDigit(code int) int {
	if code >= LatinSmallLetterA && code <= LatinSmallLetterF {
		return code - LatinSmallLetterA + 10
	}
	if code >= LatinCapitalLetterA && code <= LatinCapitalLetterF {
		return code - LatinCapitalLetterA + 10
	}
	return code - DigitZero
}
