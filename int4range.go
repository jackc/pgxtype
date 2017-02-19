package pgxtype

import (
	"io"
	"strconv"
)

type Int4range struct {
	Lower     int32
	Upper     int32
	LowerType BoundType
	UpperType BoundType
}

func (r *Int4range) ParseText(src string) error {
	utr, err := ParseUntypedTextRange(src)
	if err != nil {
		return err
	}

	n, err := strconv.ParseInt(utr.Lower, 10, 32)
	if err != nil {
		return err
	}
	r.Lower = int32(n)

	n, err = strconv.ParseInt(utr.Upper, 10, 32)
	if err != nil {
		return err
	}
	r.Upper = int32(n)

	r.LowerType = utr.LowerType
	r.UpperType = utr.UpperType

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
