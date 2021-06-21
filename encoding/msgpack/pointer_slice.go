package msgpack

type pointerSlice []uintptr

func (s pointerSlice) contains(p uintptr) bool {
	for _, ptr := range s {
		if ptr == p {
			return true
		}
	}
	return false
}

func (s *pointerSlice) pop() {
	*s = (*s)[:len(*s)-1]
}
