package template

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestToInt64(t *testing.T) {
	target := int64(102)
	if target != toInt64(int8(102)) {
		t.Errorf("Expected 102")
	}
	if target != toInt64(int(102)) {
		t.Errorf("Expected 102")
	}
	if target != toInt64(int32(102)) {
		t.Errorf("Expected 102")
	}
	if target != toInt64(int16(102)) {
		t.Errorf("Expected 102")
	}
	if target != toInt64(int64(102)) {
		t.Errorf("Expected 102")
	}
	if target != toInt64("102") {
		t.Errorf("Expected 102")
	}
	if toInt64("frankie") != 0 {
		t.Errorf("Expected 0")
	}
	if target != toInt64(uint16(102)) {
		t.Errorf("Expected 102")
	}
	if target != toInt64(uint64(102)) {
		t.Errorf("Expected 102")
	}
	if target != toInt64(float64(102.1234)) {
		t.Errorf("Expected 102")
	}
	if toInt64(true) != 1 {
		t.Errorf("Expected 102")
	}
}

func TestIncr(t *testing.T) {
	tpl := `{{ incr 7 }}`
	if err := runt(tpl, `8`); err != nil {
		t.Error(err)
	}
}

func TestDecr(t *testing.T) {
	tpl := `{{ decr 9 }}`
	if err := runt(tpl, `8`); err != nil {
		t.Error(err)
	}
}

func TestRand(t *testing.T) {
	var tests = []struct {
		min int
		max int
	}{
		{10, 11},
		{10, 13},
		{0, 1},
		{5, 50},
	}
	for _, v := range tests {
		x, _ := runRaw(fmt.Sprintf(`{{ rand %d %d }}`, v.min, v.max), nil)
		r, err := strconv.Atoi(x)
		assert.NoError(t, err)
		assert.True(t, func(min, max, r int) bool {
			return r >= v.min && r < v.max
		}(v.min, v.max, r))
	}
}
