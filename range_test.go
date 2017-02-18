package pgxtype

import (
	"testing"
)

func TestUntypedTextRange(t *testing.T) {
	tests := []struct {
		src    string
		result UntypedTextRange
		err    error
	}{
		{
			src:    `[1,2)`,
			result: UntypedTextRange{Lower: "1", Upper: "2", LowerType: '[', UpperType: ')'},
			err:    nil,
		},
		{
			src:    ` [1,2) `,
			result: UntypedTextRange{Lower: "1", Upper: "2", LowerType: '[', UpperType: ')'},
			err:    nil,
		},
	}

	for i, tt := range tests {
		r, err := NewUntypedTextRange(tt.src)
		if err != tt.err {
			t.Errorf("%d. `%s`: expected err %v, got %v", i, tt.src, tt.err, err)
			continue
		}

		if *r != tt.result {
			t.Errorf("%d. `%s`: expected result %#v, got %#v", i, tt.src, tt.result, *r)
		}
	}
}
