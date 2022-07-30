package parser

import (
	"fmt"
	"testing"
)

// Make the types prettyprint.
var tokenName = map[tokenType]string{
	ttEof:   "EOF",
	ttError: "error",

	ttBool:       "bool",
	ttComplex:    "complex",
	ttNumber:     "number",
	ttString:     "string",
	ttIdentifier: "identifier",

	// keywords
	ttDot: ".",
	ttNil: "null",
}

func (tt tokenType) String() string {
	s := tokenName[tt]
	if s == "" {
		return fmt.Sprintf("token%d", int(tt))
	}

	return s
}

func mkToken(typ tokenType, text string) token {
	return token{
		typ: typ,
		val: text,
	}
}

var (
	tDot = mkToken(ttDot, ".")
	tEof = mkToken(ttEof, "")
)

type lexTest struct {
	name   string
	input  string
	tokens []token
}

var lexTests = []lexTest{
	{"empty", "", []token{tEof}},
	{"spaces", " \t\n", []token{tEof}},
	{"bools", "true false", []token{
		mkToken(ttBool, "true"),
		mkToken(ttBool, "false"),
		tEof,
	}},
	{"null", "null", []token{
		mkToken(ttNil, "null"),
		tEof,
	}},
	{"numbers", "1 7.82 -6.76 1e3 +1.2e-4 0x14", []token{
		mkToken(ttNumber, "1"),
		mkToken(ttNumber, "7.82"),
		mkToken(ttNumber, "-6.76"),
		mkToken(ttNumber, "1e3"),
		mkToken(ttNumber, "+1.2e-4"),
		mkToken(ttNumber, "0x14"),
		tEof,
	}},
	{"dots", ". .", []token{tDot, tDot, tEof}},
	{"quote", `"abc @ pin.eu "`, []token{
		mkToken(ttString, "abc @ pin.eu "),
		tEof,
	}},
	{"json", `user = { name=Luca active=true score=8.71 likes=10 }`, []token{
		mkToken(ttIdentifier, "user"),
		mkToken(ttAssign, "="),
		mkToken(ttLeftBrace, "{"),
		mkToken(ttIdentifier, "name"),
		mkToken(ttAssign, "="),
		mkToken(ttString, "Luca"),
		mkToken(ttIdentifier, "active"),
		mkToken(ttAssign, "="),
		mkToken(ttBool, "true"),
		mkToken(ttIdentifier, "score"),
		mkToken(ttAssign, "="),
		mkToken(ttNumber, "8.71"),
		mkToken(ttIdentifier, "likes"),
		mkToken(ttAssign, "="),
		mkToken(ttNumber, "10"),
		mkToken(ttRightBrace, "}"),
		tEof,
	}},
	{"tmpl", `user = { name=Luca }`, []token{
		mkToken(ttIdentifier, "user"),
		mkToken(ttAssign, "="),
		mkToken(ttLeftBrace, "{"),
		mkToken(ttIdentifier, "name"),
		mkToken(ttAssign, "="),
		mkToken(ttString, "Luca"),
		mkToken(ttRightBrace, "}"),
		tEof,
	}},
	{"rb", `user = { id=(uuid) }`, []token{
		mkToken(ttIdentifier, "user"),
		mkToken(ttAssign, "="),
		mkToken(ttLeftBrace, "{"),
		mkToken(ttIdentifier, "id"),
		mkToken(ttAssign, "="),
		mkToken(ttExpression, "uuid"),
		mkToken(ttRightBrace, "}"),
		tEof,
	}},
}

func TestLex(t *testing.T) {
	for _, test := range lexTests {
		items := collect(&test)

		if !equal(t, items, test.tokens, false) {
			t.Errorf("%s: got\n\t%+v\nexpected\n\t%+v", test.name, items, test.tokens)
		}
	}
}

// collect gathers the emitted items into a slice.
func collect(t *lexTest) (tokens []token) {
	l := newLexer(t.input)

	for {
		el := l.nextToken()
		tokens = append(tokens, el)

		if el.typ == ttEof || el.typ == ttError {
			break
		}
	}

	return
}

func equal(t *testing.T, i1, i2 []token, checkPos bool) bool {
	if len(i1) != len(i2) {
		return false
	}

	for k := range i1 {
		t.Logf("i1: %v->%s, i2: %v->%s", i1[k].typ, i1[k].val, i2[k].typ, i2[k].val)
		if i1[k].typ != i2[k].typ {
			return false
		}

		if i1[k].val != i2[k].val {
			return false
		}

		if checkPos && i1[k].pos != i2[k].pos {
			return false
		}

		if checkPos && i1[k].line != i2[k].line {
			return false
		}
	}

	return true
}
