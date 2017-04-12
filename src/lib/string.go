package lib

// byte[] to string
func ByteToString(p []byte) string {
	for i := 0; i < len(p); i++ {
		if p[i] == 0 {
			return string(p[0:i])
		}
	}
	return string(p)
}

// 判断大写字母
func IsContainUpper(s string) (flag bool) {
	if s == "" {
		flag = false
		return
	}
	for i := 0; i < len(s); i++ {
		if 65 <= s[i] && s[i] <= 90 {
			flag = true
			break
		}

	}
	return
}
