package Tools

func GetKey(k string, m map[int]string) bool {
	for i := 0; i < len(m); i++ {
		if m[i] == k {
			return true
		}
	}
	return false
}
