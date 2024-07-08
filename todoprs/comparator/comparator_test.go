package comparator

import (
	"testing"
)

type pairTest[E comparable] struct {
	x              []E
	y              []E
	expectedResult bool
}

// CmpInsensitiveByteSlice
func TestCmpInsensitiveByteSlice(t *testing.T) {
	// case diff length , case identical, case last element diff, case first element diff. Casesensitive char
	tests := []pairTest[byte]{
		pairTest[byte]{[]byte("test"), []byte("test"), true},
		pairTest[byte]{[]byte("test"), []byte("testa"), false},
		pairTest[byte]{[]byte("test"), []byte("aest"), false},
		pairTest[byte]{[]byte("test"), []byte("tesa"), false},
		pairTest[byte]{[]byte("test"), []byte("Test"), true},
		pairTest[byte]{[]byte("test"), []byte("tesT"), true},
		pairTest[byte]{[]byte("Test"), []byte("test"), true},
		pairTest[byte]{[]byte("TesT"), []byte("test"), true},
		pairTest[byte]{[]byte("TEST"), []byte("test"), true},
		pairTest[byte]{[]byte("test"), []byte("TEST"), true},
		pairTest[byte]{[]byte("test"), []byte("TeSt"), true},
		pairTest[byte]{[]byte("tewt"), []byte("TeSt"), false},
	}
	for i, pt := range tests {
		result := CmpInsensitiveByteSlice(pt.x, pt.y)
		if result != pt.expectedResult {
			t.Fatalf("test[%d] with (%s,%s) : expecte=%t got=%t", i, pt.x, pt.y, pt.expectedResult, result)
		}
	}
}

func TestCmpSliceForBytes(t *testing.T) {
	tests := []pairTest[byte]{
		pairTest[byte]{[]byte("test"), []byte("test"), true},
		pairTest[byte]{[]byte("test"), []byte("testa"), false},
		pairTest[byte]{[]byte("test"), []byte("aest"), false},
		pairTest[byte]{[]byte("test"), []byte("tesa"), false},
		pairTest[byte]{[]byte("test"), []byte("Test"), false},
		pairTest[byte]{[]byte("test"), []byte("tesT"), false},
		pairTest[byte]{[]byte("Test"), []byte("test"), false},
		pairTest[byte]{[]byte("TesT"), []byte("test"), false},
		pairTest[byte]{[]byte("TEST"), []byte("test"), false},
		pairTest[byte]{[]byte("test"), []byte("TEST"), false},
		pairTest[byte]{[]byte("test"), []byte("TeSt"), false},
		pairTest[byte]{[]byte("tewt"), []byte("TeSt"), false},
	}
	for i, pt := range tests {
		result := CmpSlice(pt.x, pt.y)
		if result != pt.expectedResult {
			t.Fatalf("test[%d] with (%s,%s) : expecte=%t got=%t", i, pt.x, pt.y, pt.expectedResult, result)
		}
	}
}
