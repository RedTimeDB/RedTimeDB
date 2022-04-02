package numconvert

import (
	"testing"
)

func TestStringToFloat64(t *testing.T) {
	var b []byte
	b = []byte{49, 46, 50, 51}
	StringToFloat64(string(b))
}
