package regexpp_test

import (
	"testing"

	"github.com/sosukesuzuki/regexpp-go"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name      string
		inputS    string
		inputU    bool
		inputLoop int
		outputCPs []uint
	}{
		{
			name:      "ユニコードモードで、2回 Next を呼び出したときに `あ`, `い`, `う` のコードポイントを CP で参照できる",
			inputS:    "あいう",
			inputU:    true,
			inputLoop: 2,
			outputCPs: []uint{0x3042, 0x3044, 0x3046},
		},
		{
			name:      "非ユニコードモードで、2回 Next を呼び出したときに `あ`, `い`, `う` のコードポイントを CP で参照できる",
			inputS:    "あいう",
			inputU:    false,
			inputLoop: 2,
			outputCPs: []uint{0x3042, 0x3044, 0x3046},
		},
		{
			name:      "ユニコードモードで、2回 Next を呼び出したときに `あ`, `い`, `𠮟` のコードポイントを CP で参照できる",
			inputS:    "あい𠮟",
			inputU:    true,
			inputLoop: 2,
			outputCPs: []uint{0x3042, 0x3044, 0x20b9f},
		},
		{
			name:      "非ユニコードモードで、3回 Next を呼び出したときに `あ`, `い`, `𠮟`の前半部, `𠮟`の後半部 のコードポイントを CP で参照できる",
			inputS:    "あい𠮟",
			inputU:    false,
			inputLoop: 3,
			outputCPs: []uint{0x3042, 0x3044, 0xd842, 0xdf9f},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tok := regexpp.NewTokenizer(tt.inputS, tt.inputU)
			for i := 0; i < tt.inputLoop; i++ {
				if tok.CP != tt.outputCPs[i] {
					t.Errorf("Unexpected CP, expected %d, actual %d", tt.outputCPs[i], tok.CP)
				}
				tok.Next()
			}
		})
	}
}
