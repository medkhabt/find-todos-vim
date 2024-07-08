package comparator

type comparator[E comparable] func(E, E) bool

// TODO test this
func cmp[E comparable](x []E, y []E, f comparator[E]) bool {
	if x == nil && y == nil {
		return true
	}
	if x == nil || y == nil {
		return false
	}
	if len(x) != len(y) {
		return false
	}

	if len(x) == 0 {
		return true
	}
	for i := 0; i < len(x); i++ {
		if !f(x[i], y[i]) {
			return false
		}
	}
	return true
}

// TODO test this
func CmpInsensitiveByteSlice(x []byte, y []byte) bool {
	return cmp(x, y, func(u byte, w byte) bool {
		return u == w || u == w+byte(0x20) || u+byte(0x20) == w
	})
}

// TODO test this
func CmpSlice[E comparable](x []E, y []E) bool {
	return cmp(x, y, func(u E, w E) bool {
		return u == w
	})
}

func CmpSlicePointers[E comparable](x []*E, y []*E) bool {
	return cmp(x, y, func(u *E, w *E) bool {
		return *u == *w
	})
}
