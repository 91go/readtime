package readtime

import (
	"encoding/json"
	"os"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/mitchellh/mapstructure"
)

// ReadTime 阅读时间struct
type ReadTime struct {
	WordsCount  WordsCount
	Translation Translation
	// 阅读时间
	Minutes int
	// 平均每分钟阅读数
	WordsPerMinute int
}

// WordsCount 字数统计
type WordsCount struct {
	Total     int // 总字数 = Words + Puncts
	Words     int // 只包含字符数
	Puncts    int // 标点数
	Links     int // 链接数
	Pics      int // 图片数
	CodeLines int // 代码行数
}

// Translation 语言包
type Translation struct {
	Min    string
	Minute string
	Sec    string
	Second string
	Read   string
}

const (
	// DefaultWordsPerMinute 默认每分钟阅读字数
	DefaultWordsPerMinute = 300
)

// NewReadTime 实例化ReadTime
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

// SetMinutes 设置时间
func (rt *ReadTime) SetMinutes() *ReadTime {
	rt.Minutes = rt.GetMinutes()
	return rt
}

// SetWordsPerMinute 设置每分钟阅读文字数
func (rt *ReadTime) SetWordsPerMinute(words int) *ReadTime {
	rt.WordsPerMinute = words
	return rt
}

// SetTranslation 设置语言
func (rt *ReadTime) SetTranslation(key string) *ReadTime {
	if _, ok := Trans[key]; !ok || key == "" {
		return rt
	}
	rt.Translation = Trans[key]
	return rt
}

// ReadFile 直接读取文件
func (rt *ReadTime) ReadFile(filename string) (*ReadTime, error) {
	bytes, err := os.ReadFile(filename)
	if err != nil {
		return rt, err
	}
	rt.Read(string(bytes))
	return rt, nil
}

// Read 读取字符串，获取总字数
func (rt *ReadTime) Read(str string) *ReadTime {
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

// GetMinutes 计算阅读时间
func (rt *ReadTime) GetMinutes() int {
	x := float64(rt.WordsCount.Total / rt.WordsPerMinute)
	minutes := Round(x)
	if minutes <= 1 {
		return 1
	}
	return minutes
}

// ToMap 获取ReadTime的map数据
func (rt *ReadTime) ToMap() (map[string]interface{}, error) {
	rt.SetMinutes()
	ret := make(map[string]interface{})

	err := mapstructure.Decode(rt, &ret)
	if err != nil {
		return nil, err
	}
	return ret, nil
}

// ToJSON 获取ReadTime的JSON串
func (rt *ReadTime) ToJSON() (string, error) {
	toMap, err := rt.ToMap()
	if err != nil {
		return "", err
	}
	bytes, err := json.Marshal(toMap)
	if err != nil {
		return "", err
	}

	return string(bytes), nil
}
