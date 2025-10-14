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
