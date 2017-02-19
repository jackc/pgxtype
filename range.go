package pgxtype

import (
	"bytes"
	"fmt"
	"io"
	"unicode"
)

type BoundType byte

const (
	Inclusive = BoundType('[')
	Exclusive = BoundType(')')
	Unbounded = BoundType('U')
)

type UntypedTextRange struct {
	Lower     string
	Upper     string
	LowerType BoundType
	UpperType BoundType
}

func ParseUntypedTextRange(src string) (*UntypedTextRange, error) {
	buf := bytes.NewBufferString(src)
	utr := &UntypedTextRange{}

	skipWhitespace(buf)

	r, _, err := buf.ReadRune()
	if err != nil {
		return nil, fmt.Errorf("invalid lower bound: %v", err)
	}

	if r != '[' && r != '(' {
		return nil, fmt.Errorf("invalid lower bound %s", string(r))
	}

	utr.LowerType = BoundType(r)

	r, _, err = buf.ReadRune()
	if err != nil {
		return nil, fmt.Errorf("invalid lower value: %v", err)
	}
	buf.UnreadRune()

	if r == ',' {
		utr.LowerType = Unbounded
	} else {
		utr.Lower, err = rangeParseValue(buf)
		if err != nil {
			return nil, fmt.Errorf("invalid lower value: %v", err)
		}
	}

	r, _, err = buf.ReadRune()
	if err != nil {
		return nil, fmt.Errorf("missing range separator: %v", err)
	}
	if r != ',' {
		return nil, fmt.Errorf("missing range separator: %v", r)
	}

	r, _, err = buf.ReadRune()
	if err != nil {
		return nil, fmt.Errorf("invalid upper value: %v", err)
	}
	buf.UnreadRune()

	if r == ')' || r == ']' {
		utr.UpperType = Unbounded
	} else {
		utr.Upper, err = rangeParseValue(buf)
		if err != nil {
			return nil, fmt.Errorf("invalid upper value: %v", err)
		}
	}

	r, _, err = buf.ReadRune()
	if err != nil {
		return nil, fmt.Errorf("missing upper bound: %v", err)
	}
	if r != ')' && r != ']' {
		return nil, fmt.Errorf("missing upper bound, instead got: %v", string(r))
	}
	if utr.UpperType != Unbounded {
		utr.UpperType = BoundType(r)
	}

	skipWhitespace(buf)

	if buf.Len() > 0 {
		return nil, fmt.Errorf("unexpected trailing data: %v", buf.String())
	}

	return utr, nil
}

func skipWhitespace(buf *bytes.Buffer) {
	var r rune
	var err error
	for r, _, _ = buf.ReadRune(); unicode.IsSpace(r); r, _, _ = buf.ReadRune() {
	}

	if err != io.EOF {
		buf.UnreadRune()
	}
}

func rangeParseValue(buf *bytes.Buffer) (string, error) {
	r, _, err := buf.ReadRune()
	if err != nil {
		return "", err
	}
	if r == '"' {
		return rangeParseQuotedValue(buf)
	}
	buf.UnreadRune()

	s := &bytes.Buffer{}

	for {
		r, _, err := buf.ReadRune()
		if err != nil {
			return "", err
		}

		switch r {
		case '\\':
			r, _, err = buf.ReadRune()
			if err != nil {
				return "", err
			}
		case ',', '[', ']', '(', ')':
			buf.UnreadRune()
			return s.String(), nil
		}

		s.WriteRune(r)
	}
}

func rangeParseQuotedValue(buf *bytes.Buffer) (string, error) {
	s := &bytes.Buffer{}

	for {
		r, _, err := buf.ReadRune()
		if err != nil {
			return "", err
		}

		switch r {
		case '\\':
			r, _, err = buf.ReadRune()
			if err != nil {
				return "", err
			}
		case '"':
			r, _, err = buf.ReadRune()
			if err != nil {
				return "", err
			}
			if r != '"' {
				buf.UnreadRune()
				return s.String(), nil
			}
		}
		s.WriteRune(r)
	}
}
