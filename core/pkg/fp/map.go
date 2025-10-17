package fp

func Map[In any, Out any](s []In, f func(In) Out) []Out {
	output := make([]Out, len(s))
	for i, v := range s {
		output[i] = f(v)
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

func Map2List[K comparable, V any, R any](m map[K]V, f func(K, V) R) []R {
	output := make([]R, len(m))
	i := 0
	for k, v := range m {
		output[i] = f(k, v)
		i++
	}
	return output
}

func List2Map[K comparable, V any](l []V, f func(V) K) map[K]V {
	output := make(map[K]V, len(l))
	for _, v := range l {
		output[f(v)] = v
	}
	return output
}
