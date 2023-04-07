package lexer_test

import (
	"testing"

	"github.com/sosukesuzuki/regexpp-go/internal/lexer"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name      string
		inputS    string
		inputU    bool
		inputLoop int
		outputCPs []int
	}{
		{
			name:      "ユニコードモードで、2回 Next を呼び出したときに `あ`, `い`, `う` のコードポイントを CP で参照できる",
			inputS:    "あいう",
			inputU:    true,
			inputLoop: 2,
			outputCPs: []int{0x3042, 0x3044, 0x3046},
		},
		{
			name:      "非ユニコードモードで、2回 Next を呼び出したときに `あ`, `い`, `う` のコードポイントを CP で参照できる",
			inputS:    "あいう",
			inputU:    false,
			inputLoop: 2,
			outputCPs: []int{0x3042, 0x3044, 0x3046},
		},
		{
			name:      "ユニコードモードで、2回 Next を呼び出したときに `あ`, `い`, `𠮟` のコードポイントを CP で参照できる",
			inputS:    "あい𠮟",
			inputU:    true,
			inputLoop: 2,
			outputCPs: []int{0x3042, 0x3044, 0x20b9f},
		},
		{
			name:      "非ユニコードモードで、3回 Next を呼び出したときに `あ`, `い`, `𠮟`の前半部, `𠮟`の後半部 のコードポイントを CP で参照できる",
			inputS:    "あい𠮟",
			inputU:    false,
			inputLoop: 3,
			outputCPs: []int{0x3042, 0x3044, 0xd842, 0xdf9f},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tok := lexer.NewLexer(tt.inputS, tt.inputU)
			for i := 0; i < tt.inputLoop; i++ {
				if tok.CP != tt.outputCPs[i] {
					t.Errorf("Unexpected CP, expected %d, actual %d", tt.outputCPs[i], tok.CP)
				}
				tok.Next()
			}
		})
	}
}

func TestEat(t *testing.T) {
	tests := []struct {
		name         string
		inputS       string
		inputU       bool
		inputCs      []int
		inputLoop    int
		outputEatens []bool
	}{
		{
			name:         "ユニコードモードで、`あいう`に対してそれぞれの文字コードがマッチし、Eat が true を返す",
			inputS:       "あいう",
			inputU:       true,
			inputCs:      []int{0x3042, 0x3044, 0x3046},
			inputLoop:    2,
			outputEatens: []bool{true, true, true},
		},
		{
			name:         "非ユニコードモードで、`あいう`に対してそれぞれの文字コードがマッチし、Eat が true を返す",
			inputS:       "あいう",
			inputU:       false,
			inputCs:      []int{0x3042, 0x3044, 0x3046},
			inputLoop:    2,
			outputEatens: []bool{true, true, true},
		},
		{
			name:         "ユニコードモードで、`あい𠮟`に対してそれぞれの文字コードがマッチし、Eat が true を返す",
			inputS:       "あい𠮟",
			inputU:       true,
			inputCs:      []int{0x3042, 0x3044, 0x20b9f},
			inputLoop:    2,
			outputEatens: []bool{true, true, true},
		},
		{
			name:         "非ユニコードモードで、`あい𠮟`に対してそれぞれの文字コードがマッチし、Eat が true を返す",
			inputS:       "あい𠮟",
			inputU:       true,
			inputCs:      []int{0x3042, 0x3044, 0xd842, 0xdf9f},
			inputLoop:    2,
			outputEatens: []bool{true, true, true, true},
		},
		{
			name:         "ユニコードモードで、`あいう`に対してそれぞれの文字コードに対してのみ Eat が true を返す",
			inputS:       "あいう",
			inputU:       true,
			inputCs:      []int{0x3042, 0x3044, 0x3048, 0x3046},
			inputLoop:    3,
			outputEatens: []bool{true, true, false, true},
		},
		{
			name:         "非ユニコードモードで、`あいう`に対してそれぞれの文字コードに対してのみ Eat が true を返す",
			inputS:       "あいう",
			inputU:       false,
			inputCs:      []int{0x3042, 0x3044, 0x3048, 0x3046},
			inputLoop:    3,
			outputEatens: []bool{true, true, false, true},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tok := lexer.NewLexer(tt.inputS, tt.inputU)
			for i := 0; i < tt.inputLoop; i++ {
				c := tt.inputCs[i]
				e := tok.Eat(c)
				if e != tt.outputEatens[i] {
					t.Errorf("Unexpected CP, expected %t, actual %t", tt.outputEatens[i], e)
				}
			}
		})
	}
}
