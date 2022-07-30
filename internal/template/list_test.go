package template

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUniq(t *testing.T) {
	vars := map[string]interface{}{
		"listS": []string{"one", "two", "three", "one", "one"},
		"listI": []int{1, 1, 1, 2, 2, 3, 4, 3, 5, 3, 4},
	}

	tests := map[string]string{
		`{{ .listS | uniq }}`: `[one two three]`,
		`{{ .listI | uniq }}`: `[1 2 3 4 5]`,
	}
	for tpl, expect := range tests {
		assert.NoError(t, runtv(tpl, expect, vars))
	}
}

func TestHas(t *testing.T) {
	vars := map[string]interface{}{
		"list": []string{"one", "two", "three"},
	}

	tests := map[string]string{
		`{{ .list  | has "one" }}`:  `true`,
		`{{ .list  | has "four" }}`: `false`,
		`{{ has "bar" nil }}`:       `false`,
	}
	for tpl, expect := range tests {
		assert.NoError(t, runtv(tpl, expect, vars))
	}
}
