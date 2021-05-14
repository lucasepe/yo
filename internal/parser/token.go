package parser

import "fmt"

// tokenType identifies the type of lex tokens.
type tokenType int

const (
	ttEof tokenType = -1

	ttError   tokenType = iota // error occurred; value is text of error
	ttComplex                  // complex constant (1+2i); imaginary is just a number
	ttAssign                   // equals ('=') introducing an assignment

	ttIdentifier   // alphanumeric identifier
	ttLeftBrace    // '{' object begin
	ttRightBrace   // '}' object end
	ttLeftBracket  // '[' array begin
	ttRightBracket // ']' array end
	ttNumber       // simple number, including imaginary
	ttString       // string (without quotes)

	// Keywords appear after all the rest.
	ttKeyword // used only to delimit the keywords
	ttBool    // boolean constant (true or false)
	ttDot     // the cursor, spelled '.'
	ttNil     // the untyped nil constant, easiest to treat as a keyword
)

// item represents a token or text string returned from the scanner.
type token struct {
	typ  tokenType // The type of this item.
	pos  int       // The starting position, in bytes, of this item in the input string.
	val  string    // The value of this item.
	line int       // The line number at the start of this item.
}

func (t token) String() string {
	switch {
	case t.typ == ttEof:
		return "EOF"
	case t.typ == ttError:
		return t.val
	case t.typ > ttKeyword:
		return fmt.Sprintf("<%s>", t.val)
	case len(t.val) > 10:
		return fmt.Sprintf("%.10q...", t.val)
	}
	return fmt.Sprintf("%q", t.val)
}

var key = map[string]tokenType{
	".":     ttDot,
	"true":  ttBool,
	"false": ttBool,
	"null":  ttNil,
	"nil":   ttNil,
}
