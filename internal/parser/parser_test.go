package parser

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func oneFieldObjGen(field string, value interface{}) Generator {
	return mkObjectGenerator().add(field, mkValueGenerator(value))
}

func TestParseFieldGenerator(t *testing.T) {
	testCases := []struct {
		input    string
		expected Generator
	}{
		{
			input:    `a=v1`,
			expected: oneFieldObjGen("a", "v1"),
		},
		{
			input:    `a=b`,
			expected: oneFieldObjGen("a", "b"),
		},
		{
			input:    `a=42`,
			expected: oneFieldObjGen("a", int64(42)),
		},
		{
			input:    `a="42"`,
			expected: oneFieldObjGen("a", "42"),
		},
		{
			input:    `a="b c"`,
			expected: oneFieldObjGen("a", "b c"),
		},
		{
			input:    `a=true`,
			expected: oneFieldObjGen("a", true),
		},
		{
			input:    `a="true"`,
			expected: oneFieldObjGen("a", "true"),
		},
		{
			input:    `a=false`,
			expected: oneFieldObjGen("a", false),
		},
		{
			input:    `a=null`,
			expected: oneFieldObjGen("a", nil),
		},
		{
			input:    `a="just @ test . -"`,
			expected: oneFieldObjGen("a", "just @ test . -"),
		},
		{
			input:    `a="just.test"`,
			expected: oneFieldObjGen("a", "just.test"),
		},
		{
			input:    `a=(upper "pippo")`,
			expected: oneFieldObjGen("a", "PIPPO"),
		},
	}

	for _, cas := range testCases {
		t.Logf("Testing input: %s", cas.input)

		ast, err := ParseString(cas.input, nil)

		require.NoError(t, err)
		require.Equal(t, []Generator{cas.expected}, ast)
	}
}

func TestParseObjectGenerator(t *testing.T) {
	testCases := []struct {
		input    string
		expected Generator
	}{
		{
			input: `a={b=c d=2 e=true f= {g=8.8 i="l m @n"}}`,
			expected: &ObjectGenerator{
				fields: map[string]Generator{
					"a": &ObjectGenerator{
						fields: map[string]Generator{
							"b": mkValueGenerator("c"),
							"d": mkValueGenerator(int64(2)),
							"e": mkValueGenerator(true),
							"f": &ObjectGenerator{
								fields: map[string]Generator{
									"g": mkValueGenerator(float64(8.8)),
									"i": mkValueGenerator("l m @n"),
								},
							},
						},
					},
				},
			},
		},

		{
			input: `a={ b = (snakecase "FirstName") }`,
			expected: &ObjectGenerator{
				fields: map[string]Generator{
					"a": &ObjectGenerator{
						fields: map[string]Generator{
							"b": mkValueGenerator("first_name"),
						},
					},
				},
			},
		},

		{
			input: `apiVersion=v1 kind=Secret metadata.name=mysecret type=Opaque data={username=(b64enc "USER") password=(b64enc "PASS")}`,
			expected: &ObjectGenerator{
				fields: map[string]Generator{
					"apiVersion": mkValueGenerator("v1"),
					"kind":       mkValueGenerator("Secret"),
					"metadata": &ObjectGenerator{
						fields: map[string]Generator{
							"name": mkValueGenerator("mysecret"),
						},
					},
					"type": mkValueGenerator("Opaque"),
					"data": &ObjectGenerator{
						fields: map[string]Generator{
							"username": mkValueGenerator("VVNFUg=="),
							"password": mkValueGenerator("UEFTUw=="),
						},
					},
				},
			},
		},
	}

	for _, cas := range testCases {
		t.Logf("Testing input: %s", cas.input)

		ast, err := ParseString(cas.input, nil)

		require.NoError(t, err)
		require.Equal(t, []Generator{cas.expected}, ast)
	}
}

func TestParseDotObjectGenerator(t *testing.T) {
	testCases := []struct {
		input    string
		expected Generator
	}{
		{
			input: `a."b.b".c=d`,
			expected: &ObjectGenerator{
				fields: map[string]Generator{
					"a": &ObjectGenerator{
						fields: map[string]Generator{
							"b.b": &ObjectGenerator{
								fields: map[string]Generator{
									"c": mkValueGenerator("d"),
								},
							},
						},
					},
				},
			},
		},
		{
			input: `parent.child1=value1 parent.child2=value2`,
			expected: &ObjectGenerator{
				fields: map[string]Generator{
					"parent": &ObjectGenerator{
						fields: map[string]Generator{
							"child1": mkValueGenerator("value1"),
							"child2": mkValueGenerator("value2"),
						},
					},
				},
			},
		},
	}

	for _, cas := range testCases {
		t.Logf("Testing input: %s", cas.input)

		ast, err := ParseString(cas.input, nil)

		require.NoError(t, err)
		require.Equal(t, []Generator{cas.expected}, ast)
	}
}

func TestComplexParse(t *testing.T) {
	expected := &ObjectGenerator{
		fields: map[string]Generator{
			"id":      mkValueGenerator(int64(42)),
			"enabled": mkValueGenerator(true),
			"score":   mkValueGenerator(float64(8.171)),
			"caller": &ObjectGenerator{
				fields: map[string]Generator{
					"gender": &ObjectGenerator{
						fields: map[string]Generator{
							"code": mkValueGenerator("MTIz"),
						},
					},
				},
			},
			"customer": &ObjectGenerator{
				fields: map[string]Generator{
					"name": mkValueGenerator("Geralt of Rivia"),
					"age":  mkValueGenerator(int64(86)),
					"address": &ObjectGenerator{
						fields: map[string]Generator{
							"zip": mkValueGenerator("75018"),
						},
					},
				},
			},
		},
	}

	ast, err := ParseString(`id=42 score=8.171 caller.gender.code=(b64enc "123") customer={name="Geralt of Rivia" age=86 address.zip="75018"} enabled=true`, nil)

	require.NoError(t, err)
	require.Equal(t, []Generator{expected}, ast)
}

func TestArrayParse(t *testing.T) {
	expected := &ObjectGenerator{
		fields: map[string]Generator{
			"tags": &arrayGenerator{
				mkValueGenerator("d1"),
				mkValueGenerator("pluto"),
			},
		},
	}

	ast, err := ParseString(`tags = [ (regexFind "[a-zA-Z][1-9]" "abcd1234") pluto ]`, nil)

	require.NoError(t, err)
	require.Equal(t, []Generator{expected}, ast)
}
