package funcs

func Transform[V1 any, V2 any](collection []V1, transform func(V1) V2) []V2 {
	slice := make([]V2, len(collection))
	for i, elem := range collection {
		slice[i] = transform(elem)
	}
	return slice
}

func Collect[K comparable, V any](collection []K, relation map[K]V) []V {
	values := make([]V, 0, len(collection))
	for _, key := range collection {
		if value, exist := relation[key]; exist {
			values = append(values, value)
		}
	}
	return values
}

// 去重
func Distinct[K comparable, V any](collection []V, keyMap func(V) K) []V {
	keys := make(map[K]struct{}, len(collection))
	values := make([]V, 0, len(collection))
	for _, val := range collection {
		key := keyMap(val)
		if _, exist := keys[key]; exist {
			continue
		}
		keys[key] = struct{}{}
		values = append(values, val)
	}
	return values
}

// 过滤：保留 predicate 返回 true 的元素
func Filter[V any](collection []V, predicate func(V) bool) []V {
	values := make([]V, 0, len(collection))
	for _, val := range collection {
		if predicate(val) {
			values = append(values, val)
		}
	}
	return values
}
