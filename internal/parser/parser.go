package parser

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/lucasepe/yo/internal/template"
)

// parseError is returned if the input cannot be successfuly parsed
type parseError struct {
	// The original query
	input string
	// The position where the parsing fails
	pos int
	// The error message
	message string
}

func (e parseError) Error() string {
	return fmt.Sprintf("parse error: %s\n%s\n%s^", e.message, e.input, strings.Repeat(" ", e.pos))
}

type parser struct {
	lexer   *lexer
	matched token
	next    token
	ds      map[string]interface{}
}

func newParser(lex *lexer, data map[string]interface{}) *parser {
	return (&parser{
		lexer: lex,
		next:  lex.nextToken(),
		ds:    data,
	})
}

func (p *parser) parse() (gen []Generator, err error) {
	defer func() {
		if r := recover(); r != nil {
			gen = nil
			err = parseError{
				input:   p.lexer.input,
				pos:     p.matched.pos,
				message: fmt.Sprintf("%v", r),
			}
		}
	}()
	gen = p.run()
	if !p.found(ttEof) {
		p.advance()
		panic("unexpected input")
	}
	return
}

func (p *parser) run() []Generator {
	p.peek()
	if p.peek(ttLeftBrace) || p.peek(ttLeftBracket) {
		res := []Generator{}
		for {
			switch {
			case p.found(ttLeftBrace):
				res = append(res, p.object())
			case p.found(ttLeftBracket):
				res = append(res, p.array())
			default:
				return res
			}
		}
	}

	objGen := mkObjectGenerator()
	for p.found(ttIdentifier) {
		if p.peek(ttAssign) || p.peek(ttDot) {
			field := p.matched.val
			value := p.field(field)
			objGen.add(field, value)
		}
	}

	return []Generator{objGen}
}

func (p *parser) object() Generator {
	res := mkObjectGenerator()
	for p.found(ttIdentifier) {
		if p.peek(ttAssign) || p.peek(ttDot) {
			field := p.matched.val
			value := p.field(field)
			res.add(field, value)
		}
	}

	if err := p.expect(ttRightBrace); err != nil {
		panic(err)
	}

	return res
}

func (p *parser) array() Generator {
	res := &arrayGenerator{}

	for {
		switch {
		case p.found(ttExpression):
			src, err := template.ExecuteInline(p.ds, fmt.Sprintf("{{%s}}", p.matched.val))
			if err != nil {
				panic(err)
			}
			v := mkValueGenerator(string(src))
			res.add(v)

		case p.found(ttString):
			//fmt.Println(" - found string", p.matched.val)
			v := mkValueGenerator(p.matched.val)
			res.add(v)

		case p.found(ttNil):
			return mkValueGenerator(nil)

		case p.found(ttNumber):
			src, err := parseNumber(p.matched.val)
			if err != nil {
				panic(err)
			}
			v := mkValueGenerator(src)
			res.add(v)

		case p.found(ttComplex):
			src, err := parseComplex(p.matched.val)
			if err != nil {
				panic(err)
			}
			v := mkValueGenerator(src)
			res.add(v)

		case p.found(ttBool):
			src, err := parseBool(p.matched.val)
			if err != nil {
				panic(err)
			}
			v := mkValueGenerator(src)
			res.add(v)

		case p.found(ttIdentifier):
			if p.peek(ttAssign) || p.peek(ttDot) {
				field := p.matched.val
				value := p.field(field)
				// Add 1-field obj to array
				obj := mkObjectGenerator()
				obj.add(field, value)
				res.add(obj)
			} else {
				res.add(mkValueGenerator(p.matched.val))
			}
		case p.found(ttLeftBrace):
			res.add(p.object())
			//Add obj as array elem
		case p.found(ttLeftBracket):
			res.add(p.array())
			// Add array as arr elem
		case p.found(ttRightBracket):
			// return, the array is complete
			return res

		case p.found(ttEof):
			panic("unclosed array")
		default:
			p.advance()
			panic("unexpected input")
		}
	}
}

func (p *parser) field(field string) Generator {
	switch {
	case p.found(ttAssign):
		return p.value()
	case p.found(ttDot):
		//nolint:errcheck
		p.expect(ttIdentifier)

		field := p.matched.val
		value := p.field(field)
		return mkObjectGenerator().add(field, value)
	case p.found(ttEof):
		panic("unexpected end of input")
	default:
		p.advance()
		panic("unexpected input")
	}
}

func (p *parser) value() Generator {
	switch {
	case p.found(ttExpression):
		res, err := template.ExecuteInline(p.ds, fmt.Sprintf("{{%s}}", p.matched.val))
		if err != nil {
			panic(err)
		}
		return mkValueGenerator(string(res))

	case p.found(ttString):
		return mkValueGenerator(p.matched.val)

	case p.found(ttNil):
		return mkValueGenerator(nil)

	case p.found(ttNumber):
		res, err := parseNumber(p.matched.val)
		if err != nil {
			panic(err)
		}
		return mkValueGenerator(res)

	case p.found(ttComplex):
		res, err := parseComplex(p.matched.val)
		if err != nil {
			panic(err)
		}
		return mkValueGenerator(res)

	case p.found(ttBool):
		res, err := parseBool(p.matched.val)
		if err != nil {
			panic(err)
		}
		return mkValueGenerator(res)

	case p.found(ttLeftBrace):
		return p.object()

	case p.found(ttLeftBracket):
		return p.array()

	case p.found(ttEof):
		panic("unexpected end of input")

	case p.found(ttError):
		panic(p.matched.val)

	default:
		p.advance()
		panic("unexpected input")
	}
}

func (p *parser) peek(tts ...tokenType) bool {
	for _, v := range tts {
		if p.next.typ == v {
			return true
		}
	}

	return false
}

func (p *parser) found(tts ...tokenType) bool {
	if p.peek(tts...) {
		p.advance()
		return true
	}
	return false
}

func (p *parser) expect(tts ...tokenType) error {
	if !p.found(tts...) {
		p.advance()
		return fmt.Errorf("was expecting %v", tts)
	}
	return nil
}

func (p *parser) advance() {
	p.matched = p.next
	p.next = p.lexer.nextToken()
}

func parseComplex(value string) (Any, error) {
	res, err := strconv.ParseComplex(value, 128)
	if err != nil {
		return nil, fmt.Errorf("invalid literal %q: is not a complex number", value)
	}
	return res, nil
}

func parseBool(value string) (Any, error) {
	switch {
	case value == "true":
		return true, nil
	case value == "false":
		return false, nil
	default:
		return nil, fmt.Errorf("invalid literal %q: is not a boolean", value)
	}
}

func parseNumber(value string) (Any, error) {
	var v Any
	v, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		v, err = strconv.ParseFloat(value, 64)
	}
	if err != nil {
		return nil, fmt.Errorf("invalid literal %q: is not a integer or a float number", value)
	}
	return v, nil
}
