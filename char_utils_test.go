package regexpp_test

import (
	"testing"

	"github.com/sosukesuzuki/regexpp-go"
)

func TestLegacyWidth(t *testing.T) {
	tests := []struct {
		name       string
		input      uint
		wantOutput int
	}{
		{
			name:       "どんなときも1を返す",
			input:      0x3046,
			wantOutput: 1,
		},
	}

	u := regexpp.LegacyCharUtils{}

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
		wantOutput uint16
	}{
		{
			name:   "`う`のコードユニットを返す",
			inputS: "あいうえお",
			inputI: 2,
			// code point of "う"
			wantOutput: 0x3046,
		},
	}

	u := regexpp.LegacyCharUtils{}

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
		input      uint
		wantOutput int
	}{
		{
			name: "サロゲートペア以外では1を返す",
			// "う"のコードポイント
			input:      0x3046,
			wantOutput: 1,
		},
		{
			name: "サロゲートペアには2を返す",
			// "𠮟"のコードポイント
			input:      0x20b9f,
			wantOutput: 2,
		},
	}

	u := regexpp.UnicodeCharUtils{}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := u.Width(tt.input)
			if w != tt.wantOutput {
				t.Errorf("Unexpected width, expected %d, actual %d", tt.wantOutput, w)
			}
		})
	}
}
