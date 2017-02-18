package pgxtype

import (
	"io"
)

type Int4range struct {
	lower     int32
	upper     int32
	lowerType BoundType
	upperType BoundType
}

func (r *Int4range) ParseText(src string) error {
	return nil
}

func (r *Int4range) ParseBinary(src []byte) error {
	return nil
}

func (r *Int4range) FormatText(w io.Writer) error {
	return nil
}

func (r *Int4range) FormatBinary(w io.Writer) error {
	return nil
}
