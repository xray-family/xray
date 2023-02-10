package internal

func JoinPath(sep string, ss ...string) string {
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
