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

func Join1(a string, sep string) string {
	if len(a) == 0 {
		return sep
	}
	if a[0:1] != sep {
		a = sep + a
	}
	return trimRight(a, sep)
}

func Join2(a string, b string, sep string) string {
	b = trimRight(b, sep)
	var m = len(a)
	var n = len(b)
	if m == 0 && n == 0 {
		return sep
	}

	var f1 = m > 0 && a[m-1:m] == sep
	var f2 = n > 0 && b[0:1] == sep
	if f1 && f2 {
		return a + b[1:]
	} else if f1 || f2 {
		if m == 0 {
			return b
		}
		if n == 0 {
			return a[:m-1]
		}
		return a + b
	}
	return a + sep + b
}

func trimRight(path string, sep string) string {
	var n = len(path)
	if n == 0 {
		return path
	}
	if path[n-1:] == sep {
		return path[:n-1]
	}
	return path
}
