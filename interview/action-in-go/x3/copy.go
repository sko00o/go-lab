package cp

// 关于只是值传递和指针传递的对比
// ref: https://medium.com/@blanchon.vincent/go-should-i-use-a-pointer-instead-of-a-copy-of-my-struct-44b43b104963

type S struct {
	a, b, c int64
	d, e, f string
	g, h, i float64
}

func byCopy() S {
	return S{
		1, 1, 1,
		"f", "o", "o",
		0.1, 0.1, 0.1,
	}
}

func byPointer() *S {
	return &S{
		1, 1, 1,
		"f", "o", "o",
		0.1, 0.1, 0.1,
	}
}

func (s S) stack(S)  {}
func (s *S) heap(*S) {}
