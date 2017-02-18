package pgxtype

import (
	"errors"
	"unicode"
	"unicode/utf8"
)

var errParseRange = errors.New("unable to parse range")

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

func NewUntypedTextRange(src string) (*UntypedTextRange, error) {
	lex := &rangeLex{src: src,
		tr:    &UntypedTextRange{},
		state: lexLeadingWhitespace,
	}

	var err error
	for lex.state != nil {
		lex.state, err = lex.state(lex)
		if err != nil {
			return nil, err
		}
	}

	return lex.tr, nil
}

type stateFn func(*rangeLex) (stateFn, error)

type rangeLex struct {
	src   string
	start int
	pos   int
	width int
	state stateFn
	tr    *UntypedTextRange
}

func (l *rangeLex) next() (r rune) {
	if l.pos >= len(l.src) {
		l.width = 0 // because backing up from having read eof should read eof again
		return 0
	}

	r, l.width = utf8.DecodeRuneInString(l.src[l.pos:])
	l.pos += l.width

	return r
}

func (l *rangeLex) unnext() {
	l.pos -= l.width
}

func (l *rangeLex) ignore() {
	l.start = l.pos
}

func (l *rangeLex) acceptRunFunc(f func(rune) bool) {
	for f(l.next()) {
	}
	l.unnext()
}

func lexLeadingWhitespace(l *rangeLex) (stateFn, error) {
	switch r := l.next(); {
	case r == '[' || r == '(':
		return lexLowerBound, nil
	case isWhitespace(r):
		l.skipWhitespace()
		return lexLeadingWhitespace, nil
	default:
		return nil, errParseRange
	}
}

func lexLowerBound(l *rangeLex) (stateFn, error) {
	l.tr.LowerType = BoundType(l.src[l.start])
	l.start = l.pos
	return lexLowerValue, nil
}

func lexLowerValue(l *rangeLex) (stateFn, error) {
	for {
		switch r := l.next(); {
		case r == 0:
			return nil, errParseRange
		case r == '"':
			return lexLowerQuotedValue, nil
		case r == ',':
			l.tr.Lower = l.src[l.start : l.pos-1]
			l.start = l.pos
			return lexUpperValue, nil
		default:
		}
	}
}

func lexLowerQuotedValue(l *rangeLex) (stateFn, error) {
	return nil, errors.New("lexLowerQuotedValue not implemented")
}

func lexUpperValue(l *rangeLex) (stateFn, error) {
	for {
		switch r := l.next(); {
		case r == 0:
			return nil, errParseRange
		case r == '"':
			return lexUpperQuotedValue, nil
		case r == ')' || r == ']':
			l.unnext()
			l.tr.Upper = l.src[l.start:l.pos]
			l.start = l.pos
			return lexUpperBound, nil
		default:
		}
	}
}

func lexUpperQuotedValue(l *rangeLex) (stateFn, error) {
	return nil, errors.New("lexUpperQuotedValue not implemented")
}

func lexUpperBound(l *rangeLex) (stateFn, error) {
	switch r := l.next(); {
	case BoundType(r) == Inclusive:
		l.tr.UpperType = Inclusive
		return lexTrailingWhitespace, nil
	case BoundType(r) == Exclusive:
		l.tr.UpperType = Exclusive
		return lexTrailingWhitespace, nil
	default:
		return nil, errParseRange
	}
}

func lexTrailingWhitespace(l *rangeLex) (stateFn, error) {
	switch r := l.next(); {
	case r == 0:
		return nil, nil
	case isWhitespace(r):
		l.skipWhitespace()
		return lexTrailingWhitespace, nil
	default:
		return nil, errParseRange
	}
}

func (l *rangeLex) skipWhitespace() {
	var r rune
	for r = l.next(); isWhitespace(r); r = l.next() {
	}

	if r != 0 {
		l.unnext()
	}

	l.ignore()
}

func isWhitespace(r rune) bool {
	return unicode.IsSpace(r)
}
