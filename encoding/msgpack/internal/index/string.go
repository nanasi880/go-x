package index

import (
	"fmt"
	"reflect"
	"sync"
)

var (
	stringIndexes        = make(map[cacheKey]map[string]int)
	stringOrderedIndexes = make(map[cacheKey][]StringKey)
	stringIndexesMutex   sync.RWMutex
)

type StringKey struct {
	Key   string
	Index int
}

// GetString returns struct indexes of map encoding.
func GetString(t reflect.Type, tagName string) (map[string]int, error) {
	key := cacheKey{
		t:   t,
		tag: tagName,
	}

	if indexes := lookupString(key, false); indexes != nil {
		return indexes, nil
	}

	stringIndexesMutex.Lock()
	defer stringIndexesMutex.Unlock()

	if indexes := lookupString(key, true); indexes != nil {
		return indexes, nil
	}

	ordered, unordered, err := buildString(t, tagName)
	if err != nil {
		return nil, err
	}
	setString(key, ordered, unordered)

	return unordered, nil
}

// GetStringOrdered returns struct indexes of map encoding.
func GetStringOrdered(t reflect.Type, tagName string) ([]StringKey, error) {
	key := cacheKey{
		t:   t,
		tag: tagName,
	}

	if indexes := lookupOrderedString(key, false); indexes != nil {
		return indexes, nil
	}

	stringIndexesMutex.Lock()
	defer stringIndexesMutex.Unlock()

	if indexes := lookupOrderedString(key, true); indexes != nil {
		return indexes, nil
	}

	ordered, unordered, err := buildString(t, tagName)
	if err != nil {
		return nil, err
	}
	setString(key, ordered, unordered)

	return ordered, nil
}

func lookupString(key cacheKey, locked bool) map[string]int {
	if !locked {
		stringIndexesMutex.RLock()
		defer stringIndexesMutex.RUnlock()
	}

	indexes, ok := stringIndexes[key]
	if !ok {
		return nil
	}

	return indexes
}

func lookupOrderedString(key cacheKey, locked bool) []StringKey {
	if !locked {
		stringIndexesMutex.RLock()
		defer stringIndexesMutex.RUnlock()
	}

	indexes, ok := stringOrderedIndexes[key]
	if !ok {
		return nil
	}

	return indexes
}

func buildString(t reflect.Type, tagName string) ([]StringKey, map[string]int, error) {
	var (
		ordered   = make([]StringKey, 0)
		unordered = make(map[string]int)
	)
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		tag, ok := field.Tag.Lookup(tagName)
		if ok && tag == "-" {
			continue
		}
		if tag == "" {
			tag = field.Name
		}

		if _, ok := unordered[tag]; ok {
			return nil, nil, fmt.Errorf("encode key conflict: %s", tag)
		}
		ordered = append(ordered, StringKey{
			Key:   tag,
			Index: i,
		})
		unordered[tag] = i
	}

	return ordered, unordered, nil
}

func setString(key cacheKey, ordered []StringKey, unordered map[string]int) {
	stringIndexes[key] = unordered
	stringOrderedIndexes[key] = ordered
}
