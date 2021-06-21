package set

type Int struct {
	m map[int]struct{}
}

func (i *Int) Insert(v int) *Int {
	if i.m == nil {
		i.m = make(map[int]struct{})
	}
	i.m[v] = struct{}{}
	return i
}

func (i *Int) Erase(v int) *Int {
	delete(i.m, v)
	return i
}

func (i *Int) Contains(v int) bool {
	_, ok := i.m[v]
	return ok
}
