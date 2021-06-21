package index

import (
	"fmt"
	"reflect"
	"sync"
)

var (
	stringIndexes        = make(map[reflect.Type]map[string]int)
	stringOrderedIndexes = make(map[reflect.Type][]StringKey)
	stringIndexesMutex   sync.RWMutex
)

type StringKey struct {
	Key   string
	Index int
}

// GetString returns struct indexes of map encoding.
func GetString(t reflect.Type) map[string]int {

	stringIndexesMutex.RLock()
	if it, ok := stringIndexes[t]; ok {
		stringIndexesMutex.RUnlock()
		return it
	}
	stringIndexesMutex.RUnlock()

	stringIndexesMutex.Lock()
	defer stringIndexesMutex.Unlock()

	if it, ok := stringIndexes[t]; ok {
		return it
	}

	ordered, unordered := buildString(t)
	stringOrderedIndexes[t] = ordered
	stringIndexes[t] = unordered

	return unordered
}

// GetStringOrdered returns struct indexes of map encoding.
func GetStringOrdered(t reflect.Type) []StringKey {

	stringIndexesMutex.RLock()
	if it, ok := stringOrderedIndexes[t]; ok {
		stringIndexesMutex.RUnlock()
		return it
	}
	stringIndexesMutex.RUnlock()

	stringIndexesMutex.Lock()
	defer stringIndexesMutex.Unlock()

	if it, ok := stringOrderedIndexes[t]; ok {
		return it
	}

	ordered, unordered := buildString(t)
	stringOrderedIndexes[t] = ordered
	stringIndexes[t] = unordered

	return ordered
}

func buildString(t reflect.Type) ([]StringKey, map[string]int) {
	var (
		ordered   = make([]StringKey, 0)
		unordered = make(map[string]int)
	)
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		tag, ok := field.Tag.Lookup("msgpack")
		if ok && tag == "-" {
			continue
		}
		if tag == "" {
			tag = field.Name
		}

		if _, ok := unordered[tag]; ok {
			panic(
				fmt.Sprintf("encode key conflict: %s", tag),
			)
		}
		ordered = append(ordered, StringKey{
			Key:   tag,
			Index: i,
		})
		unordered[tag] = i
	}

	return ordered, unordered
}
