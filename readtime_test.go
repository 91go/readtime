package readtime

import (
	"fmt"
	"testing"
)

func TestReadTime(t *testing.T) {
	minutes := NewReadTime().ReadFile("./readtime.go").SetTranslation("ca").ToMap()
	j := NewReadTime().ReadFile("./readtime.go").SetTranslation("cn").ToJSON()

	fmt.Println(minutes, j)
}

func TestStat(t *testing.T) {
	tests := []struct {
		name  string
		input string
		total int
		words int
	}{
		{"en1", "hello,playground", 3, 2},
		{"en2", "hello, playground", 3, 2},
		{"cn1", "你好世界", 4, 4},
		{"encn1", "Hello你好世界", 5, 5},
		{"encn2", "Hello 你好世界", 5, 5},
		{"encn3", "Hello，你好世界。", 7, 5},
		{"link1", "Hello，你好世界。https://studygolang.com Go中文网", 11, 9},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			readTime := NewReadTime().ReadStr(tt.input)
			if readTime.WordsCount.Total != tt.total || readTime.WordsCount.Words != tt.words {
				t.Errorf("Total = %v, want %v; Words=%v, want %v",
					readTime.WordsCount.Total, tt.total, readTime.WordsCount.Words, tt.words)
			}
		})
	}
}
