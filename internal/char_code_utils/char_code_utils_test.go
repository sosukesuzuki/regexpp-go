package char_code_utils_test

import (
	"testing"

	"github.com/sosukesuzuki/regexpp-go/internal/char_code_utils"
)

func TestLegacyWidth(t *testing.T) {
	tests := []struct {
		name       string
		input      int
		wantOutput int
	}{
		{
			name:       "どんなときも1を返す",
			input:      0x3046,
			wantOutput: 1,
		},
	}

	u := char_code_utils.Legacy{}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := u.Width(tt.input)
			if w != tt.wantOutput {
				t.Errorf("Unexpected width, expected %d, actual %d", tt.wantOutput, w)
			}
		})
	}
}

func TestLegacyAt(t *testing.T) {
	tests := []struct {
		name       string
		inputS     string
		inputI     int
		wantOutput int
	}{
		{
			name:       "`う`のコードユニットを返す",
			inputS:     "あいうえお",
			inputI:     2,
			wantOutput: 0x3046,
		},
		{
			name:       "`𠮟`のサロゲートペアの前半部のコードユニットを返す",
			inputS:     "あい𠮟えお",
			inputI:     2,
			wantOutput: 0xd842,
		},
	}

	u := char_code_utils.Legacy{}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := u.At(tt.inputS, tt.inputI)
			if w != tt.wantOutput {
				t.Errorf("Unexpected at, expected %d, actual %d", tt.wantOutput, w)
			}
		})
	}
}

func TestUnicodeWidth(t *testing.T) {
	tests := []struct {
		name       string
		input      int
		wantOutput int
	}{
		{
			name:       "サロゲートペア以外では1を返す",
			input:      0x3046,
			wantOutput: 1,
		},
		{
			name:       "サロゲートペアには2を返す",
			input:      0x20b9f,
			wantOutput: 2,
		},
	}

	u := char_code_utils.Unicode{}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := u.Width(tt.input)
			if w != tt.wantOutput {
				t.Errorf("Unexpected width, expected %d, actual %d", tt.wantOutput, w)
			}
		})
	}
}

func TestUnicodeAt(t *testing.T) {
	tests := []struct {
		name       string
		inputS     string
		inputI     int
		wantOutput int
	}{
		{
			name:       "`う`のコードポイントを返す",
			inputS:     "あいうえお",
			inputI:     2,
			wantOutput: 0x3046,
		},
		{
			name:       "`𠮟`のコードポイントを返す",
			inputS:     "あい𠮟えお",
			inputI:     2,
			wantOutput: 0x20b9f,
		},
	}

	u := char_code_utils.Unicode{}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := u.At(tt.inputS, tt.inputI)
			if w != tt.wantOutput {
				t.Errorf("Unexpected at, expected %d, actual %d", tt.wantOutput, w)
			}
		})
	}
}
