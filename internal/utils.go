package internal

import "strings"

const (
	Separator     = "/"
	SeparatorByte = '/'
)

// JoinPath 拼接请求路径
func JoinPath(ss ...string) string {
	var cursor = 0
	var b = make([]byte, 0, 32)
	b = append(b, SeparatorByte)
	cursor++

	for _, v := range ss {
		var n = len(v)
		if n == 0 {
			continue
		}

		if b[cursor-1] != SeparatorByte {
			b = append(b, SeparatorByte)
			cursor++
		}
		for j := 0; j < n; j++ {
			if !(b[cursor-1] == SeparatorByte && v[j] == SeparatorByte) {
				b = append(b, v[j])
				cursor++
			}
		}
	}

	if cursor > 1 && b[cursor-1] == SeparatorByte {
		return string(b[:cursor-1])
	}
	return string(b)
}

// Split 分割字符串(空值将会被过滤掉)
func Split(s string) []string {
	var list = strings.Split(s, Separator)
	var j = 0
	for _, v := range list {
		if v = strings.TrimSpace(v); v != "" {
			list[j] = v
			j++
		}
	}
	return list[:j]
}

// TrimPath 去除路径两边多余的斜杠
func TrimPath(path string) string {
	path = strings.TrimSpace(path)
	n := len(path)
	if n == 0 {
		return Separator
	}

	if path[0] != SeparatorByte {
		path = Separator + path
	}
	if n >= 2 && path[0]+path[1] == 94 {
		return TrimPath(path[1:])
	}
	if path[n-1] == SeparatorByte {
		return TrimPath(path[:n-1])
	}
	return path
}

// SelectString 三元操作
func SelectString(expression bool, a, b string) string {
	if expression {
		return a
	}
	return b
}

// GetMaxLength 获取数组中最长字符串的长度
func GetMaxLength(args ...string) int {
	var x = 0
	for _, v := range args {
		if n := len(v); n > x {
			x = n
		}
	}
	return x
}

// Padding 填充空格, 使字符串到达指定长度
func Padding(s string, length int) string {
	var b = []byte(s)
	for len(b) < length {
		b = append(b, ' ')
	}
	return string(b)
}

// FastSplit 快速分割字符串, 0 alloc
// str是预先格式化好的, 必须以斜杠开头
func FastSplit(str string, f func(segment string) bool) {
	var n = len(str)
	var i = 1
	var j = i
	for k := i; k < n; k++ {
		if str[k] == SeparatorByte || k == n-1 {
			if k == n-1 {
				j++
			}
			if !f(str[i:j]) {
				return
			}
			i = k + 1
			j = i
		} else {
			j++
		}
	}
}
