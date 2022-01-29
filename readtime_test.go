package readtime

import (
	"fmt"
	"testing"
)

func TestReadTime(t *testing.T) {
	builder := NewReadTimeBuilder()
	right := builder.WithIsLeftToRight(false)
	fmt.Println(right)
}
