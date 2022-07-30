package template

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDefault(t *testing.T) {
	tpl := `{{"" | default "foo"}}`
	if err := runt(tpl, "foo"); err != nil {
		t.Error(err)
	}
	tpl = `{{default "foo" 234}}`
	if err := runt(tpl, "234"); err != nil {
		t.Error(err)
	}
	tpl = `{{default "foo" 2.34}}`
	if err := runt(tpl, "2.34"); err != nil {
		t.Error(err)
	}

	tpl = `{{ .Nothing | default "123" }}`
	if err := runt(tpl, "123"); err != nil {
		t.Error(err)
	}
	tpl = `{{ default "123" }}`
	if err := runt(tpl, "123"); err != nil {
		t.Error(err)
	}
}

func TestCoalesce(t *testing.T) {
	tests := map[string]string{
		`{{ coalesce 1 }}`:                            "1",
		`{{ coalesce "" 0 nil 2 }}`:                   "2",
		`{{ $two := 2 }}{{ coalesce "" 0 nil $two }}`: "2",
		`{{ $two := 2 }}{{ coalesce "" $two 0 0 0 }}`: "2",
		`{{ $two := 2 }}{{ coalesce "" $two 3 4 5 }}`: "2",
		`{{ coalesce }}`:                              "<no value>",
	}
	for tpl, expect := range tests {
		assert.NoError(t, runt(tpl, expect))
	}

	dict := map[string]interface{}{"top": map[string]interface{}{}}
	tpl := `{{ coalesce .top.NoSuchThing .bottom .bottom.dollar "airplane"}}`
	if err := runtv(tpl, "airplane", dict); err != nil {
		t.Error(err)
	}
}

func TestAll(t *testing.T) {
	tests := map[string]string{
		`{{ all 1 }}`:                            "true",
		`{{ all "" 0 nil 2 }}`:                   "false",
		`{{ $two := 2 }}{{ all "" 0 nil $two }}`: "false",
		`{{ $two := 2 }}{{ all "" $two 0 0 0 }}`: "false",
		`{{ $two := 2 }}{{ all "" $two 3 4 5 }}`: "false",
		`{{ all }}`:                              "true",
	}
	for tpl, expect := range tests {
		assert.NoError(t, runt(tpl, expect))
	}

	dict := map[string]interface{}{"top": map[string]interface{}{}}
	tpl := `{{ all .top.NoSuchThing .bottom .bottom.dollar "airplane"}}`
	if err := runtv(tpl, "false", dict); err != nil {
		t.Error(err)
	}
}

func TestAny(t *testing.T) {
	tests := map[string]string{
		`{{ any 1 }}`:                              "true",
		`{{ any "" 0 nil 2 }}`:                     "true",
		`{{ $two := 2 }}{{ any "" 0 nil $two }}`:   "true",
		`{{ $two := 2 }}{{ any "" $two 3 4 5 }}`:   "true",
		`{{ $zero := 0 }}{{ any "" $zero 0 0 0 }}`: "false",
		`{{ any }}`: "false",
	}
	for tpl, expect := range tests {
		assert.NoError(t, runt(tpl, expect))
	}

	dict := map[string]interface{}{"top": map[string]interface{}{}}
	tpl := `{{ any .top.NoSuchThing .bottom .bottom.dollar "airplane"}}`
	if err := runtv(tpl, "true", dict); err != nil {
		t.Error(err)
	}
}
