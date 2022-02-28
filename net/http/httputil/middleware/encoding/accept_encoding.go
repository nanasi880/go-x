package encoding

type acceptEncoding struct {
	encoding   string
	weight     int
	incomplete bool
}

type acceptEncodings []acceptEncoding

// Len is implementation of sort.Interface
func (s acceptEncodings) Len() int {
	return len(s)
}

// Less is implementation of sort.Interface
func (s acceptEncodings) Less(i, j int) bool {
	a := s[i]
	b := s[j]

	if a.weight != b.weight {
		return a.weight > b.weight
	}
	if a.incomplete == b.incomplete {
		return false
	}
	if !a.incomplete && b.incomplete {
		return true
	}
	return true
}

// Swap is implementation of sort.Interface
func (s acceptEncodings) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
