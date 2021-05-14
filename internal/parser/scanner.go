package parser

import (
	"bufio"
	"io"
	"strings"
)

const (
	scannerBuffer = 128 * 1024
)

// ParseString accepts an input string and
// Returns either a slice of Generators on success or else an error.
func ParseString(input string) ([]Generator, error) {
	lexer := newLexer(input)
	return (&parser{
		lexer: lexer,
		next:  lexer.nextToken(),
	}).parse()
}

// ParseTextLines parse a slice of lines.
// Returns either a slice of Generators on success or else an error.
func ParseTextLines(lines []string) ([]Generator, error) {
	spec := strings.Join(lines, " ")
	return ParseString(spec)
}

// ParseTextLines parse a reader.
// Returns either a slice of Generators on success or else an error.
func ParseReader(reader io.Reader) ([]Generator, error) {
	buffer := make([]byte, scannerBuffer)

	scanner := bufio.NewScanner(reader)
	scanner.Buffer(buffer, scannerBuffer)

	res := []string{}

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		res = append(res, line)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return ParseTextLines(res)
}
