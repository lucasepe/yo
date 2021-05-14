package parser

import (
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"
)

const eof = -1

type lexer struct {
	input     string    // the string being scanned
	pos       int       // current position in the input
	start     int       // start position of this item
	width     int       // width of last rune read from input
	line      int       // 1+number of newlines seen
	startLine int       // start line of this item
	lastSeen  tokenType // the last seen token type
}

// newLexer creates a new scanner for the input string.
func newLexer(input string) *lexer {
	l := &lexer{
		input:     input,
		line:      1,
		startLine: 1,
	}

	return l
}

func (l *lexer) nextToken() token {
	for {
		l.skipSpaces()

		switch r := l.pop(); {
		case r == eof:
			return l.emit(ttEof)
		case r == '{':
			return l.emit(ttLeftBrace)
		case r == '}':
			return l.emit(ttRightBrace)
		case r == '[':
			return l.emit(ttLeftBracket)
		case r == ']':
			return l.emit(ttRightBracket)
		case r == '=':
			return l.emit(ttAssign)
		case r == '"':
			return l.lexQuotedString()
		case r == '.':
			x := l.peek()
			if x < '0' || '9' < x {
				return l.emit(ttDot)
			}
			fallthrough // '.' can start a number.
		case r == '+' || r == '-' || ('0' <= r && r <= '9'):
			l.push()
			return l.lexNumber()
		case isAlphaNumeric(r):
			l.push()
			return l.lexIdentifier()
		default:
			l.push()
			return l.errorf("unrecognized character: %#U", r)
		}
	}
}

// lexIdentifier scans an alphanumeric.
func (l *lexer) lexIdentifier() token {
	for {
		switch r := l.pop(); {
		case isAlphaNumeric(r):
			// absorb.
		default:
			l.push()

			if l.lastSeen != ttAssign && !l.atTerminator() {
				return l.errorf("bad character %#U", r)
			}

			word := l.input[l.start:l.pos]
			if key[word] > ttKeyword {
				return l.emit(key[word])
			}

			if l.lastSeen == ttAssign {
				return l.emit(ttString)
			}

			return l.emit(ttIdentifier)
		}
	}
}

// lexQuotedString scans a quoted string.
func (l *lexer) lexQuotedString() token {
Loop:
	for {
		switch l.pop() {
		case '\\':
			if r := l.pop(); r != eof && r != '\n' {
				break
			}
			fallthrough
		case eof, '\n':
			return l.errorf("unterminated quoted string")
		case '"':
			break Loop
		}
	}

	val := l.input[l.start+1 : l.pos-1]
	return l.emitV(ttString, val)
}

// atTerminator reports whether the input is at valid termination character to
// appear after an identifier.
func (l *lexer) atTerminator() bool {
	r := l.peek()
	if isSpace(r) {
		return true
	}

	switch r {
	case eof, '=', '.', '[', ']', '{', '}':
		return true
	}

	return false
}

func (l *lexer) lexNumber() token {
	if !l.scanNumber() {
		return l.errorf("bad number syntax: %q", l.input[l.start:l.pos])
	}

	if sign := l.peek(); sign == '+' || sign == '-' {
		// Complex: 1+2i. No spaces, must end in 'i'.
		if !l.scanNumber() || l.input[l.pos-1] != 'i' {
			return l.errorf("bad number syntax: %q", l.input[l.start:l.pos])
		}

		return l.emit(ttComplex)
	}

	return l.emit(ttNumber)
}

func (l *lexer) scanNumber() bool {
	// Optional leading sign.
	l.accept("+-")

	// Is it hex?
	digits := "0123456789_"
	if l.accept("0") {
		// Note: Leading 0 does not mean octal in floats.
		if l.accept("xX") {
			digits = "0123456789abcdefABCDEF_"
		} else if l.accept("oO") {
			digits = "01234567_"
		} else if l.accept("bB") {
			digits = "01_"
		}
	}

	l.acceptRun(digits)
	if l.accept(".") {
		l.acceptRun(digits)
	}

	if len(digits) == 10+1 && l.accept("eE") {
		l.accept("+-")
		l.acceptRun("0123456789_")
	}

	if len(digits) == 16+6+1 && l.accept("pP") {
		l.accept("+-")
		l.acceptRun("0123456789_")
	}

	// Is it imaginary?
	l.accept("i")

	// Next thing mustn't be alphanumeric.
	if isAlphaNumeric(l.peek()) {
		l.pop()
		return false
	}

	return true
}

// skipSpaces eats all spaces.
func (l *lexer) skipSpaces() {
	for isSpace(l.peek()) {
		l.pop()
	}
	l.drop()
}

// pop returns the next rune in the input.
func (l *lexer) pop() rune {
	if l.pos >= len(l.input) {
		l.width = 0
		return eof
	}

	r, w := utf8.DecodeRuneInString(l.input[l.pos:])
	l.width = w
	l.pos += l.width

	if r == '\n' {
		l.line++
	}

	return r
}

// peek returns but does not consume the next rune in the input.
func (l *lexer) peek() rune {
	r := l.pop()
	l.push()
	return r
}

// push steps back one rune. Can only be called once per call of next.
func (l *lexer) push() {
	l.pos -= l.width
	// Correct newline count.
	if l.width == 1 && l.input[l.pos] == '\n' {
		l.line--
	}
}

func (l *lexer) drop() {
	l.start = l.pos
}

// accept consumes the next rune if it's from the valid set.
func (l *lexer) accept(valid string) bool {
	if strings.ContainsRune(valid, l.pop()) {
		return true
	}
	l.push()

	return false
}

// acceptRun consumes a run of runes from the valid set.
func (l *lexer) acceptRun(valid string) {
	for strings.ContainsRune(valid, l.pop()) {

	}
	l.push()
}

func (l *lexer) emit(typ tokenType) token {
	val := l.input[l.start:l.pos]
	return l.emitV(typ, val)
}

func (l *lexer) emitV(typ tokenType, v string) token {
	res := token{
		typ:  typ,
		val:  v,
		pos:  l.start,
		line: l.line,
	}
	l.start = l.pos
	l.lastSeen = typ
	return res
}

// errorf returns an error token and terminates the scan by passing
// back a nil pointer that will be the next state, terminating l.nextItem.
func (l *lexer) errorf(format string, args ...interface{}) token {
	return l.emitV(ttError, fmt.Sprintf(format, args...))
}

// isSpace reports whether r is a space character.
func isSpace(r rune) bool {
	return r == ' ' || r == '\t' || r == '\r' || r == '\n'
}

// isAlphaNumeric reports whether r is an alphabetic, digit, underscore or dash.
func isAlphaNumeric(r rune) bool {
	return r == '_' ||
		r == '-' ||
		unicode.IsLetter(r) ||
		unicode.IsDigit(r)
}
