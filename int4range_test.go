package pgxtype

import (
	"testing"
)

func TestInt4range(t *testing.T) {
	conns := mustConnectAll(t)
	defer mustCloseAll(t, conns)

	tests := []struct {
		name   string
		sql    string
		args   []interface{}
		err    error
		result Int4range
	}{
		{
			name:   "Scan",
			sql:    "select int4range(1, 10)",
			args:   []interface{}{},
			err:    nil,
			result: Int4range{Lower: 1, Upper: 10, LowerType: Inclusive, UpperType: Exclusive},
		},
		{
			name:   "Scan Negative",
			sql:    "select int4range(-42, -5)",
			args:   []interface{}{},
			err:    nil,
			result: Int4range{Lower: -42, Upper: -5, LowerType: Inclusive, UpperType: Exclusive},
		},
		{
			name:   "Scan Normalized Bounds",
			sql:    "select int4range(1, 10, '(]')",
			args:   []interface{}{},
			err:    nil,
			result: Int4range{Lower: 2, Upper: 11, LowerType: Inclusive, UpperType: Exclusive},
		},
	}

	for _, conn := range conns {
		for _, tt := range tests {
			var r Int4range
			var s string
			err := conn.QueryRow(tt.sql, tt.args...).Scan(&s)
			if err != tt.err {
				t.Errorf("%s %s: %v", conn.DriverName(), tt.name, err)
			}

			err = r.ParseText(s)
			if err != nil {
				t.Errorf("%s %s: %v", conn.DriverName(), tt.name, err)
			}

			if r != tt.result {
				t.Errorf("%s %s: expected %#v, got %#v", conn.DriverName(), tt.name, tt.result, r)
			}
		}
	}
}
