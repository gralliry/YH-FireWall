package fc

func List2List[In any, Out any](s []In, f func(In) Out) []Out {
	output := make([]Out, len(s))
	for i, v := range s {
		output[i] = f(v)
	}
	return output
}

func Map2List[K comparable, V any, E any](m map[K]V, f func(K, V) E) []E {
	output := make([]E, len(m))
	i := 0
	for k, v := range m {
		output[i] = f(k, v)
		i++
	}
	return output
}

func List2Map[E any, K comparable, V any](l []E, f func(E) (K, V)) map[K]V {
	output := make(map[K]V, len(l))
	for _, e := range l {
		k, v := f(e)
		output[k] = v
	}
	return output
}

func Map2Map[K1 comparable, V1 any, K2 comparable, V2 any](m map[K1]V1, f func(K1, V1) (K2, V2)) map[K2]V2 {
	output := make(map[K2]V2, len(m))
	for k, v := range m {
		k2, v2 := f(k, v)
		output[k2] = v2
	}
	return output
}

func Filter[Elem any](s []Elem, f func(Elem) bool) []Elem {
	output := make([]Elem, 0)
	for _, v := range s {
		if f(v) {
			output = append(output, v)
		}
	}
	return output
}

func Distinct[K comparable, V any](l []V, f func(V) K) []V {
	output := make([]V, 0)
	seen := make(map[K]struct{})
	for _, v := range l {
		k := f(v)
		if _, e := seen[k]; e {
			continue
		}
		seen[k] = struct{}{}
		output = append(output, v)
	}
	return output
}

func Limit[T any](l []T, n int) []T {
	n = min(max(0, n), len(l))
	output := make([]T, n)
	copy(output, l)
	return output
}

func Reduce[T1, T2 any](arr []T1, zero T2, f func(T2, T1) T2) T2 {
	result := zero
	for _, v := range arr {
		result = f(result, v)
	}
	return result
}
