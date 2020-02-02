package dynasty

type Layers struct {
	data map[string]string
}

func encode(n, size int) []byte {
	buf := make([]byte, size)
	nums := num(n, []int{})
	for i, a := range nums {
		buf[i] = byte(rune(97 + a))
	}
	for i := len(nums); i < size; i++ {
		buf[i] = byte('_')
	}
	return buf
}

// base 26
func num(n int, r []int) []int {
	d := n / 26
	r = append(r, n%26)
	if d > 0 {
		return num(d, r)
	}
	return r
}
