package template

import (
	crypto_rand "crypto/rand"
	"encoding/binary"
	math_rand "math/rand"
	"time"

	"github.com/lucasepe/yo/internal/cast"
)

func init() {
	var b [8]byte
	_, err := crypto_rand.Read(b[:])
	if err != nil {
		math_rand.Seed(time.Now().UnixNano())
	} else {
		math_rand.Seed(int64(binary.LittleEndian.Uint64(b[:])))
	}
}

func randInt(min, max int) int {
	return math_rand.Intn(max-min) + min
}

// toInt64 converts integer types to 64-bit integers
func toInt64(v interface{}) int64 {
	return cast.ToInt64(v)
}
