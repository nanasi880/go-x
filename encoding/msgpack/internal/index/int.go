package index

import (
	"fmt"
	"math"
	"reflect"
	"strconv"
	"sync"
)

var (
	intIndexes      = make(map[cacheKey][]int)
	intIndexesMutex sync.RWMutex
)

// GetInt returns struct indexes of array encoding.
func GetInt(t reflect.Type, tagName string) ([]int, error) {
	key := cacheKey{
		t:   t,
		tag: tagName,
	}

	if indexes := lookupInt(key, false); indexes != nil {
		return indexes, nil
	}

	intIndexesMutex.Lock()
	defer intIndexesMutex.Unlock()

	if indexes := lookupInt(key, true); indexes != nil {
		return indexes, nil
	}

	indexes, err := buildInt(t, tagName)
	if err != nil {
		return nil, err
	}
	intIndexes[key] = indexes

	return indexes, nil
}

func lookupInt(key cacheKey, locked bool) []int {
	if !locked {
		intIndexesMutex.RLock()
		defer intIndexesMutex.RUnlock()
	}

	indexes, ok := intIndexes[key]
	if !ok {
		return nil
	}

	return indexes
}

func buildInt(t reflect.Type, tagName string) ([]int, error) {

	var (
		keyToIndex = make(map[int]int)
		maxKey     = -1
	)

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		tag, ok := field.Tag.Lookup(tagName)
		if !ok {
			return nil, fmt.Errorf("struct tag `%s` is not found: %v", tagName, t)
		}
		if tag == "-" {
			continue
		}

		i64, err := strconv.ParseInt(tag, 10, 32)
		if err != nil {
			return nil, fmt.Errorf("struct tag `%s` parse error: %v: %v", tagName, t, err)
		}
		if i64 < 0 {
			return nil, fmt.Errorf("encode key cannot use negative value: %d", i64)
		}
		if i64 > math.MaxInt32 {
			return nil, fmt.Errorf("encode key cannot use greater or equal math.MaxInt32: %d", i64)
		}

		key := int(i64)
		if _, ok := keyToIndex[key]; ok {
			return nil, fmt.Errorf("encode key conflict: %d", key)
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

	return indexes, nil
}
