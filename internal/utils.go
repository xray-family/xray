package internal

func Join1(a string, sep string) string {
	if len(a) == 0 {
		return sep
	}
	if a[0:1] == sep {
		return a
	}
	return sep + a
}

func Join2(a string, b string, sep string) string {
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

func SelectString(expression bool, a, b string) string {
	if expression {
		return a
	}
	return b
}