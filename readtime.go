package readtime

import (
	"math"
)

type ReadTimeBuilder struct {
	ReadTime
}

type ReadTime struct {
	// 是否要缩写，minutes/seconds是否要缩写为min/sec
	Abbreviated bool
	// 文章正文内容
	Content       string
	Estimate      map[string]string
	IsLeftToRight bool
	// 是否秒数
	OmitSeconds bool
	// 是否只展示时间
	TimeOnly     bool
	Translations map[string]string
	// 文章正文总字数
	WordsCount int
	// 平均每分钟阅读数
	WordsPerMinute int
}

func NewReadTimeBuilder() *ReadTimeBuilder {
	rt := new(ReadTimeBuilder)
	rt.ReadTime.TimeOnly = true

	return rt
}

func (rt *ReadTimeBuilder) WithIsLeftToRight(isLeftToRight bool) *ReadTimeBuilder {
	rt.ReadTime.IsLeftToRight = isLeftToRight
	return rt
}

func (this *ReadTime) calculateMinutes() int {
	minutes := math.Floor(float64(this.WordsCount / this.WordsPerMinute))
	if minutes < 1 {
		return 1
	}
	return int(minutes)
}
