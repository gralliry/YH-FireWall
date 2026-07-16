package bimap

// Map 是双向映射，A 和 B 完全对等，任一方均可查询另一方。
type Map[A, B comparable] struct {
	a2b map[A]B
	b2a map[B]A
}

// New 创建一个空的双向映射。
func New[A, B comparable]() *Map[A, B] {
	return &Map[A, B]{
		a2b: make(map[A]B),
		b2a: make(map[B]A),
	}
}

// Insert 插入一对映射关系。如果 A 或 B 已存在，旧关系会被覆盖。
func (m *Map[A, B]) Insert(a A, b B) {
	if oldB, ok := m.a2b[a]; ok {
		delete(m.b2a, oldB)
	}
	if oldA, ok := m.b2a[b]; ok {
		delete(m.a2b, oldA)
	}
	m.a2b[a] = b
	m.b2a[b] = a
}

// GetByA 通过 A 查询对应的 B。
func (m *Map[A, B]) GetByA(a A) (B, bool) {
	b, ok := m.a2b[a]
	return b, ok
}

// GetByB 通过 B 查询对应的 A。
func (m *Map[A, B]) GetByB(b B) (A, bool) {
	a, ok := m.b2a[b]
	return a, ok
}

// DeleteByA 通过 A 删除映射关系。
func (m *Map[A, B]) DeleteByA(a A) {
	if b, ok := m.a2b[a]; ok {
		delete(m.b2a, b)
		delete(m.a2b, a)
	}
}

// DeleteByB 通过 B 删除映射关系。
func (m *Map[A, B]) DeleteByB(b B) {
	if a, ok := m.b2a[b]; ok {
		delete(m.a2b, a)
		delete(m.b2a, b)
	}
}

// Len 返回映射关系数量。
func (m *Map[A, B]) Len() int {
	return len(m.a2b)
}

// As 返回所有 A 的切片（顺序不保证）。
func (m *Map[A, B]) As() []A {
	as := make([]A, 0, len(m.a2b))
	for a := range m.a2b {
		as = append(as, a)
	}
	return as
}

// Bs 返回所有 B 的切片（顺序不保证）。
func (m *Map[A, B]) Bs() []B {
	bs := make([]B, 0, len(m.b2a))
	for b := range m.b2a {
		bs = append(bs, b)
	}
	return bs
}
