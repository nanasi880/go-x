package index

import "reflect"

type cacheKey struct {
	t   reflect.Type
	tag string
}
