package convertx

import (
	"strings"
	"unicode"
	"unsafe"

	"github.com/mozillazg/go-pinyin"
)

// Bytes2String 字节数组转字符串
func Bytes2String(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

// String2Bytes 字符串转字节数组
func String2Bytes(s string) []byte {
	return *(*[]byte)(unsafe.Pointer(&struct {
		string
		len int
	}{
		string: s,
		len:    len(s),
	}))
}

// Chinese2Spell 汉字转拼音
func Chinese2Spell(words string) string {
	b := strings.Builder{}
	a := pinyin.NewArgs()
	for i, c := range words {
		if unicode.Is(unicode.Han, c) {
			b.WriteString(pinyin.SinglePinyin(c, a)[0])
			continue
		}
		b.WriteByte(words[i])
	}
	return b.String()
}
