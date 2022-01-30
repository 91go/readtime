package readtime

import (
	"fmt"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

// func TestReadTime(t *testing.T) {
// 	minutes := NewReadTime().ReadFile("./readtime.go").SetTranslation("ca").ToMap()
// 	j := NewReadTime().ReadFile("./readtime.go").SetTranslation("cn").ToJSON()
// 	NewReadTime().ReadStr("").ToMap()
//
// 	fmt.Println(minutes, j)
// }

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
			readTime := NewReadTime().Read(tt.input)
			if readTime.WordsCount.Total != tt.total || readTime.WordsCount.Words != tt.words {
				t.Errorf("Total = %v, want %v; Words=%v, want %v",
					readTime.WordsCount.Total, tt.total, readTime.WordsCount.Words, tt.words)
			}
		})
	}
}

func TestReadTime_GetMinutes(t *testing.T) {
	Convey("minutes", t, func() {
		Convey("minutes < 1", func() {
			minutes := NewReadTime().Read("xxx").GetMinutes()
			So(minutes, ShouldBeLessThanOrEqualTo, 1)
		})
		Convey("minutes > 1", func() {
			minutes, _ := NewReadTime().ReadFile("./testdata/big-file.html")
			So(minutes.GetMinutes(), ShouldBeGreaterThan, 1)
		})
	})
}

func TestReadTime_ToMap(t *testing.T) {
	Convey("ToMap", t, func() {
		Convey("wpm", func() {
			minutes := NewReadTime().Read("设计模式").SetWordsPerMinute(1).GetMinutes()
			So(minutes, ShouldEqual, 4)
		})
		Convey("translation", func() {
			Convey("error", func() {
				Convey("lang is blank", func() {
					minutes := NewReadTime().Read("abc").SetTranslation("").GetMinutes()
					So(minutes, ShouldEqual, 1)
				})
				Convey("lang not exist", func() {
					minutes := NewReadTime().Read("abc").SetTranslation("yyds").GetMinutes()
					So(minutes, ShouldEqual, 1)
				})
			})
			Convey("success", func() {
				minutes := NewReadTime().Read("abc").SetTranslation("cn").GetMinutes()
				So(minutes, ShouldEqual, 1)
			})
		})
		Convey("read-file", func() {
			Convey("error", func() {
				_, err := NewReadTime().ReadFile("")
				So(err, ShouldBeError)
			})
		})
		Convey("map", func() {
			Convey("success", func() {
				toMap, err := NewReadTime().Read("abc").SetTranslation("cn").ToMap()
				fmt.Println(toMap)
				So(toMap, ShouldNotEqual, map[string]interface{}{})
				So(err, ShouldBeNil)
			})
		})
		Convey("json", func() {
			Convey("success", func() {
				toJSON, err := NewReadTime().Read("abc").SetTranslation("cn").ToJSON()
				fmt.Println(toJSON)
				So(toJSON, ShouldNotEqual, "")
				So(toJSON, ShouldEqual, `{"Minutes":1,"Translation":{"Min":"分","Minute":"分钟","Read":"阅读","Sec":"秒","Second":"秒钟"},"WordsCount":{"CodeLines":0,"Links":0,"Pics":0,"Puncts":0,"Total":1,"Words":1},"WordsPerMinute":300}`)
				So(err, ShouldBeNil)
			})
		})
	})
}
