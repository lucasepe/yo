package template

import (
	"encoding/base64"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSubstr(t *testing.T) {
	tpl := `{{"fooo" | substr 0 3 }}`
	if err := runt(tpl, "foo"); err != nil {
		t.Error(err)
	}
}

func TestSubstr_shorterString(t *testing.T) {
	tpl := `{{"foo" | substr 0 10 }}`
	if err := runt(tpl, "foo"); err != nil {
		t.Error(err)
	}
}

func TestContains(t *testing.T) {
	// Mainly, we're just verifying the paramater order swap.
	tests := []string{
		`{{if contains "cat" "fair catch"}}1{{end}}`,
		`{{if hasPrefix "cat" "catch"}}1{{end}}`,
		`{{if hasSuffix "cat" "ducat"}}1{{end}}`,
	}
	for _, tt := range tests {
		if err := runt(tt, "1"); err != nil {
			t.Error(err)
		}
	}
}

func TestSplit(t *testing.T) {
	tpl := `{{$v := "foo$bar$baz" | split "$"}}{{index $v 0}}`
	if err := runt(tpl, "foo"); err != nil {
		t.Error(err)
	}
}

func TestJoin(t *testing.T) {
	assert.NoError(t, runtv(`{{ join "-" .V }}`, "a-b-c", map[string]interface{}{"V": []string{"a", "b", "c"}}))
	assert.NoError(t, runtv(`{{ join "-" .V }}`, "abc", map[string]interface{}{"V": "abc"}))
	assert.NoError(t, runtv(`{{ join "-" .V }}`, "1-2-3", map[string]interface{}{"V": []int{1, 2, 3}}))
	assert.NoError(t, runtv(`{{ join "-" .value }}`, "1-2", map[string]interface{}{"value": []interface{}{"1", nil, "2"}}))
}

func TestBase64EncodeDecode(t *testing.T) {
	magicWord := "coffee"
	expect := base64.StdEncoding.EncodeToString([]byte(magicWord))

	if expect == magicWord {
		t.Fatal("Encoder doesn't work.")
	}

	tpl := `{{b64enc "coffee"}}`
	if err := runt(tpl, expect); err != nil {
		t.Error(err)
	}
	tpl = fmt.Sprintf("{{b64dec %q}}", expect)
	if err := runt(tpl, magicWord); err != nil {
		t.Error(err)
	}
}

func TestReplace(t *testing.T) {
	tpl := `{{"I Am Henry VIII" | replace " " "-"}}`
	if err := runt(tpl, "I-Am-Henry-VIII"); err != nil {
		t.Error(err)
	}
}
