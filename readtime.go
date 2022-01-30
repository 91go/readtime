package readtime

import (
	"encoding/json"
	"os"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/mitchellh/mapstructure"
)

type ReadTime struct {
	WordsCount  WordsCount
	Translation Translation
	// 阅读时间
	Minutes int
	// 平均每分钟阅读数
	WordsPerMinute int
}

type WordsCount struct {
	Total     int // 总字数 = Words + Puncts
	Words     int // 只包含字符数
	Puncts    int // 标点数
	Links     int // 链接数
	Pics      int // 图片数
	CodeLines int // 代码行数
}

type Translation struct {
	Min    string
	Minute string
	Sec    string
	Second string
	Read   string
}

const (
	DefaultWordsPerMinute = 300
)

func NewReadTime() *ReadTime {
	return &ReadTime{
		WordsCount: WordsCount{
			Total:     0,
			Words:     0,
			Puncts:    0,
			Links:     0,
			Pics:      0,
			CodeLines: 0,
		},
		Translation:    Trans["en"],
		WordsPerMinute: DefaultWordsPerMinute,
		Minutes:        0,
	}
}

func (rt *ReadTime) SetMinutes() *ReadTime {
	rt.Minutes = rt.GetMinutes()
	return rt
}

func (rt *ReadTime) SetWordsPerMinute(words int) *ReadTime {
	rt.WordsPerMinute = words
	return rt
}

func (rt *ReadTime) SetTranslation(key string) *ReadTime {
	if _, ok := Trans[key]; !ok || key == "" {
		return rt
	}
	rt.Translation = Trans[key]
	return rt
}

func (rt *ReadTime) ReadFile(filename string) *ReadTime {
	bytes, err := os.ReadFile(filename)
	if err != nil {
		return nil
	}
	rt.ReadStr(string(bytes))
	return rt
}

// Stat
func (rt *ReadTime) ReadStr(str string) *ReadTime {
	rt.WordsCount.Links = len(rxStrict.FindAllString(str, -1))
	rt.WordsCount.Pics = len(imgReg.FindAllString(str, -1))

	// 剔除 HTML
	str = StripHTML(str)
	str = AutoSpace(str)

	// 普通的链接去除（非 HTML 标签链接）
	str = rxStrict.ReplaceAllString(str, " ")
	plainWords := strings.Fields(str)

	for _, plainWord := range plainWords {
		words := strings.FieldsFunc(plainWord, func(r rune) bool {
			if unicode.IsPunct(r) {
				rt.WordsCount.Puncts++
				return true
			}
			return false
		})

		for _, word := range words {
			runeCount := utf8.RuneCountInString(word)
			if len(word) == runeCount {
				rt.WordsCount.Words++
			} else {
				rt.WordsCount.Words += runeCount
			}
		}
	}

	rt.WordsCount.Total = rt.WordsCount.Words + rt.WordsCount.Puncts
	return rt
}

// 计算阅读时间
func (rt *ReadTime) GetMinutes() int {
	x := float64(rt.WordsCount.Total / rt.WordsPerMinute)
	minutes := Round(x)
	if minutes < 1 {
		return 1
	}
	return minutes
}

func (rt *ReadTime) ToMap() map[string]interface{} {
	rt.SetMinutes()
	ret := make(map[string]interface{})

	err := mapstructure.Decode(rt, &ret)
	if err != nil {
		return nil
	}
	return ret
}

func (rt *ReadTime) ToJSON() string {
	toMap := rt.ToMap()
	bytes, err := json.Marshal(toMap)
	if err != nil {
		return ""
	}

	return string(bytes)
}
