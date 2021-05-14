package parser

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestValueGenerator(t *testing.T) {
	cases := []interface{}{
		1,
		1.5,
		nil,
		true,
		"a",
	}

	otherGens := []Generator{
		mkObjectGenerator(),
		&arrayGenerator{},
		mkValueGenerator(42),
	}

	for _, cas := range cases {
		v := mkValueGenerator(cas)

		require.Equal(t, cas, v.Get())

		for _, og := range otherGens {
			// Merge of a value wih another generator should always return the other
			require.Equal(t, og, v.Merge(og))
		}
	}
}

func TestArrayGenerator(t *testing.T) {
	values := []Any{6, true, " aloha"}

	g := &arrayGenerator{}
	for _, v := range values {
		g.add(mkValueGenerator(v))
	}

	require.Equal(t, values, g.Get())

	otherGens := []Generator{
		mkObjectGenerator(),
		&arrayGenerator{},
		mkValueGenerator(42),
	}

	for _, og := range otherGens {
		// Merge of an array wih another generator should always return the other
		require.Equal(t, og, g.Merge(og))
	}
}
