package dynasty

// Layers keep layers ranks, and is encoded path
type Layers struct {
	data map[string][]byte
}

// NewLayers returns a new Layers
func NewLayers() *Layers {
	return &Layers{
		data: make(map[string][]byte),
	}
}

func (l *Layers) layer(layer string) []byte {
	ll, ok := l.data[layer]
	if ok {
		return ll
	}
	l.data[layer] = encode(len(l.data), 3)
	return l.data[layer]
}

// encode a number : n is converted to base 26, with fixed number of digits
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
