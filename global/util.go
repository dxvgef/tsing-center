package global

import (
	"encoding/base64"
	"strings"
	"sync"
	"unsafe"
)

func BytesToStr(value []byte) string {
	return *(*string)(unsafe.Pointer(&value)) // nolint
}

// 编码键名
func EncodeKey(value string) string {
	return base64.RawURLEncoding.EncodeToString(StrToBytes(value))
}

// 解码键名
func DecodeKey(value string) (string, error) {
	keyBytes, err := base64.RawURLEncoding.DecodeString(value)
	if err != nil {
		return "", err
	}
	return BytesToStr(keyBytes), nil
}

func StrToBytes(s string) []byte {
	x := (*[2]uintptr)(unsafe.Pointer(&s)) // nolint
	h := [3]uintptr{x[0], x[1], x[1]}
	return *(*[]byte)(unsafe.Pointer(&h)) // nolint
}

func FormatTime(str string) string {
	str = strings.Replace(str, "y", "2006", -1)
	str = strings.Replace(str, "m", "01", -1)
	str = strings.Replace(str, "d", "02", -1)
	str = strings.Replace(str, "h", "15", -1)
	str = strings.Replace(str, "i", "04", -1)
	str = strings.Replace(str, "s", "05", -1)
	return str
}

// 计算sync.Map的长度
func SyncMapLen(m *sync.Map) (count int) {
	if m == nil {
		return
	}
	m.Range(func(key, value interface{}) bool {
		count++
		return true
	})
	return
}
