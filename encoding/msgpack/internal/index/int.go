package index

import (
	"fmt"
	"math"
	"reflect"
	"strconv"
	"sync"
)

var (
	intIndexes      = make(map[reflect.Type][]int)
	intIndexesMutex sync.RWMutex
)

// GetInt returns struct indexes of array encoding.
func GetInt(t reflect.Type) []int {

	intIndexesMutex.RLock()
	if it, ok := intIndexes[t]; ok {
		intIndexesMutex.RUnlock()
		return it
	}
	intIndexesMutex.RUnlock()

	intIndexesMutex.Lock()
	defer intIndexesMutex.Unlock()

	if it, ok := intIndexes[t]; ok {
		return it
	}

	it := buildInt(t)
	intIndexes[t] = it

	return it
}

func buildInt(t reflect.Type) []int {

	var (
		keyToIndex = make(map[int]int)
		maxKey     = -1
	)

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		tag, ok := field.Tag.Lookup("msgpack")
		if !ok {
			panic(
				fmt.Sprintf("struct tag `msgpack` is not found: %v", t),
			)
		}
		if tag == "-" {
			continue
		}

		i64, err := strconv.ParseInt(tag, 10, 32)
		if err != nil {
			panic(
				fmt.Sprintf("struct tag `msgpack` parse error: %v: %v", t, err),
			)
		}
		if i64 < 0 {
			panic(
				fmt.Sprintf("encode key cannot use negative value: %d", i64),
			)
		}
		if i64 > math.MaxInt32 {
			panic(
				fmt.Sprintf("encode key cannot use greater or equal math.MaxInt32: %d", i64),
			)
		}

		key := int(i64)
		if _, ok := keyToIndex[key]; ok {
			panic(
				fmt.Sprintf("encode key conflict: %d", key),
			)
		}
		keyToIndex[key] = i

		if maxKey < key {
			maxKey = key
		}
	}

	indexes := make([]int, 0, maxKey+1)

	for i := 0; i <= maxKey; i++ {
		index, ok := keyToIndex[i]
		if !ok {
			indexes = append(indexes, -1)
		} else {
			indexes = append(indexes, index)
		}
	}

	return indexes
}
