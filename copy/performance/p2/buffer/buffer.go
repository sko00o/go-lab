package buffer

func HandRange(s1, s2 []uint16) (n uint32) {
	if len(s1) == 0 || len(s2) == 0 {
		return 0
	}

	i2, l2, v2 := 0, len(s2), s2[0]
	var v1 uint16
	for i1 := 0; i1 < len(s1); i1++ {
		v1 = s1[i1]
		for v2 < v1 {
			i2++
			if i2 >= l2 {
				return
			}
			v2 = s2[i2]
		}
		if v2 == v1 {
			n++
		}
	}

	return
}

func IterRange(s1, s2 []uint16) (n uint32) {
	if len(s1) == 0 || len(s2) == 0 {
		return 0
	}

	i2, l2, v2 := 0, len(s2), s2[0]
	for _, v1 := range s1 {
		for v2 < v1 {
			i2++
			if i2 >= l2 {
				return
			}
			v2 = s2[i2]
		}
		if v2 == v1 {
			n++
		}
	}

	return
}
