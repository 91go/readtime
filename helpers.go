package readtime

import (
	"math"
	"regexp"
	"strings"
	"unicode"
	"unicode/utf8"

	"mvdan.cc/xurls/v2"
)

// AutoSpace 自动给中英文之间加上空格
func AutoSpace(str string) string {
	out := ""

	for _, r := range str {
		out = addSpaceAtBoundary(out, r)
	}

	return out
}

func addSpaceAtBoundary(prefix string, nextChar rune) string {
	if prefix == "" {
		return string(nextChar)
	}

	r, size := utf8.DecodeLastRuneInString(prefix)
	if isLatin(size) != isLatin(utf8.RuneLen(nextChar)) &&
		isAllowSpace(nextChar) && isAllowSpace(r) {
		return prefix + " " + string(nextChar)
	}

	return prefix + string(nextChar)
}

var (
	rxStrict          = xurls.Strict()
	imgReg            = regexp.MustCompile(`<img [^>]*>`)
	stripHTMLReplacer = strings.NewReplacer("\n", " ", "</p>", "\n", "<br>", "\n", "<br />", "\n")
)

// StripHTML accepts a string, strips out all HTML tags and returns it.
func StripHTML(s string) string {
	// Shortcut strings with no tags in them
	if !strings.ContainsAny(s, "<>") {
		return s
	}
	s = stripHTMLReplacer.Replace(s)

	// Walk through the string removing all tags
	b := GetBuffer()
	defer PutBuffer(b)
	var inTag, isSpace, wasSpace bool
	for _, r := range s {
		if !inTag {
			isSpace = false
		}

		switch {
		case r == '<':
			inTag = true
		case r == '>':
			inTag = false
		case unicode.IsSpace(r):
			isSpace = true
			fallthrough
		default:
			if !inTag && (!isSpace || (isSpace && !wasSpace)) {
				b.WriteRune(r)
			}
		}

		wasSpace = isSpace
	}
	return b.String()
}

func isLatin(size int) bool {
	return size == 1
}

func isAllowSpace(r rune) bool {
	return !unicode.IsSpace(r) && !unicode.IsPunct(r)
}

// Round 四舍五入
func Round(x float64) int {
	return int(math.Floor(x + 0.5))
}
