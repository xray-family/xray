package internal

import "strings"

// JoinPath 拼接请求路径
// 如果搬来就是有效路径, 直接返回
func JoinPath(sep string, ss ...string) string {
	switch len(ss) {
	case 0:
		return sep
	case 1:
		if ss[0] == "" || ss[0] == sep {
			return sep
		}
		if isValidPath(sep, ss[0]) {
			return ss[0]
		}
	}

	var ch = sep[0]
	var cursor = 0
	var b = make([]byte, 0, 32)
	b = append(b, ch)
	cursor++

	for _, v := range ss {
		var n = len(v)
		if n == 0 {
			continue
		}

		if b[cursor-1] != ch {
			b = append(b, ch)
			cursor++
		}
		for j := 0; j < n; j++ {
			if !(b[cursor-1] == ch && v[j] == ch) {
				b = append(b, v[j])
				cursor++
			}
		}
	}

	if cursor > 1 && b[cursor-1] == ch {
		return string(b[:cursor-1])
	}
	return string(b)
}

func isValidPath(sep string, path string) bool {
	var ch = sep[0]
	var n = len(path)
	if path[0] != ch || path[n-1] == ch {
		return false
	}
	for i := 1; i < n; i++ {
		if path[i] == ch && path[i-1] == ch {
			return false
		}
	}
	return true
}

// SelectString 三元操作
func SelectString(expression bool, a, b string) string {
	if expression {
		return a
	}
	return b
}

// Split 分割字符串(空值将会被过滤掉)
func Split(s string, sep string) []string {
	var list = strings.Split(s, sep)
	var j = 0
	for _, v := range list {
		if v = strings.TrimSpace(v); v != "" {
			list[j] = v
			j++
		}
	}
	return list[:j]
}

func GetMaxLength(args ...string) int {
	var x = 0
	for _, v := range args {
		if n := len(v); n > x {
			x = n
		}
	}
	return x
}

func Padding(s string, length int) string {
	var b = []byte(s)
	for len(b) < length {
		b = append(b, ' ')
	}
	return string(b)
}
