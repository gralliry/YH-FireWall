package multikeymap

type entry[K comparable, V any] struct {
	keys  []K
	value V
}

type Map[K comparable, V any] struct {
	l map[*entry[K, V]]struct{}
	m map[K]*entry[K, V]
}

func New[K comparable, V any]() *Map[K, V] {
	return &Map[K, V]{
		l: make(map[*entry[K, V]]struct{}),
		m: make(map[K]*entry[K, V]),
	}
}

func (m *Map[K, V]) Set(value V, keys ...K) {
	for _, k := range keys {
		if e, ok := m.m[k]; ok {
			m.unindex(e)
		}
	}
	e := &entry[K, V]{
		keys:  append([]K(nil), keys...),
		value: value,
	}
	for _, k := range keys {
		m.m[k] = e
	}
	m.l[e] = struct{}{}
}

func (m *Map[K, V]) Get(key K) (V, bool) {
	if e, ok := m.m[key]; ok {
		return e.value, true
	}
	var zero V
	return zero, false
}

func (m *Map[K, V]) Del(key K) bool {
	e, ok := m.m[key]
	if !ok {
		return false
	}
	m.unindex(e)
	return true
}

func (m *Map[K, V]) unindex(e *entry[K, V]) {
	for _, k := range e.keys {
		delete(m.m, k)
	}
	delete(m.l, e)
}

func (m *Map[K, V]) Len() int {
	return len(m.l)
}

func (m *Map[K, V]) Values() []V {
	values := make([]V, 0, len(m.l))
	for e := range m.l {
		values = append(values, e.value)
	}
	return values
}

func (m *Map[K, V]) Extract(f func(V) bool) []V {
	var removed []V
	for e := range m.l {
		if f(e.value) {
			removed = append(removed, e.value)
			m.unindex(e)
		}
	}
	return removed
}

func (m *Map[K, V]) Range(f func(V)) {
	for e := range m.l {
		f(e.value)
	}
}
