package helper

import (
	"hash"
	"strings"
)

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

// MustHash 哈希
func MustHash(data []byte, algo hash.Hash, encoding func([]byte) string) string {
	_, _ = algo.Write(data)
	return encoding(algo.Sum(nil))
}
