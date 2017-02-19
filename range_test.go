package pgxtype

import (
	"testing"
)

func TestParseUntypedTextRange(t *testing.T) {
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
			src:    `[1,2]`,
			result: UntypedTextRange{Lower: "1", Upper: "2", LowerType: '[', UpperType: ']'},
			err:    nil,
		},
		{
			src:    `(1,3)`,
			result: UntypedTextRange{Lower: "1", Upper: "3", LowerType: '(', UpperType: ')'},
			err:    nil,
		},
		{
			src:    ` [1,2) `,
			result: UntypedTextRange{Lower: "1", Upper: "2", LowerType: '[', UpperType: ')'},
			err:    nil,
		},
		{
			src:    `[ foo , bar )`,
			result: UntypedTextRange{Lower: " foo ", Upper: " bar ", LowerType: '[', UpperType: ')'},
			err:    nil,
		},
		{
			src:    `["foo","bar")`,
			result: UntypedTextRange{Lower: "foo", Upper: "bar", LowerType: '[', UpperType: ')'},
			err:    nil,
		},
		{
			src:    `["f""oo","b""ar")`,
			result: UntypedTextRange{Lower: `f"oo`, Upper: `b"ar`, LowerType: '[', UpperType: ')'},
			err:    nil,
		},
		{
			src:    `["f""oo","b""ar")`,
			result: UntypedTextRange{Lower: `f"oo`, Upper: `b"ar`, LowerType: '[', UpperType: ')'},
			err:    nil,
		},
		{
			src:    `["","bar")`,
			result: UntypedTextRange{Lower: ``, Upper: `bar`, LowerType: '[', UpperType: ')'},
			err:    nil,
		},
		{
			src:    `[f\"oo\,,b\\ar\))`,
			result: UntypedTextRange{Lower: `f"oo,`, Upper: `b\ar)`, LowerType: '[', UpperType: ')'},
			err:    nil,
		},
		{
			src:    `empty`,
			result: UntypedTextRange{Lower: "", Upper: "", LowerType: 'E', UpperType: 'E'},
			err:    nil,
		},
	}

	for i, tt := range tests {
		r, err := ParseUntypedTextRange(tt.src)
		if err != tt.err {
			t.Errorf("%d. `%s`: expected err %v, got %v", i, tt.src, tt.err, err)
			continue
		}

		if *r != tt.result {
			t.Errorf("%d. `%s`: expected result %#v, got %#v", i, tt.src, tt.result, *r)
		}
	}
}
