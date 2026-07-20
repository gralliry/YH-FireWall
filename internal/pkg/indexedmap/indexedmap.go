package indexedmap

import "math/rand/v2"

const maxLevel = 32

type node[V any] struct {
	value V

	prev []*node[V]
	next []*node[V]
}

type IndexedMap[K comparable, V any] struct {
	head node[V]

	level int
	size  int

	compare func(a, b V) int

	keymap map[K]*node[V]
}

func New[K comparable, V any](compare func(a, b V) int) *IndexedMap[K, V] {

	return &IndexedMap[K, V]{
		head: node[V]{
			next: make([]*node[V], maxLevel),
			prev: make([]*node[V], maxLevel),
		},
		level:   1,
		compare: compare,
		keymap:  make(map[K]*node[V]),
	}
}

// randomLevel 返回随机的节点层数（1～maxLevel）
func (m *IndexedMap[K, V]) randomLevel() int {
	level := 1
	for level < maxLevel && rand.IntN(2) == 0 {
		level++
	}
	return level
}

// Insert 插入元素，若 key 已存在则覆盖
func (m *IndexedMap[K, V]) Insert(key K, value V) {
	if old, ok := m.keymap[key]; ok {
		m.removeNode(old)
		delete(m.keymap, key)
	}

	n := &node[V]{value: value}
	level := m.randomLevel()
	n.next = make([]*node[V], level)
	n.prev = make([]*node[V], level)

	update := make([]*node[V], maxLevel)

	// 从最高层向下，查找每层的前驱
	cur := &m.head
	for i := m.level - 1; i >= 0; i-- {
		for cur.next[i] != nil && m.compare(cur.next[i].value, value) < 0 {
			cur = cur.next[i]
		}
		update[i] = cur
	}

	// 若新节点层数超过当前最高层，补全前驱为 head
	if level > m.level {
		for i := m.level; i < level; i++ {
			update[i] = &m.head
		}
		m.level = level
	}

	// 将节点插入各层链表
	for i := range level {
		next := update[i].next[i]
		n.next[i] = next
		n.prev[i] = update[i]
		update[i].next[i] = n
		if next != nil {
			next.prev[i] = n
		}
	}

	m.keymap[key] = n
	m.size++
}

// Search 根据 key 查找元素
func (m *IndexedMap[K, V]) Search(key K) (V, bool) {
	n, ok := m.keymap[key]
	if !ok {
		var zero V
		return zero, false
	}
	return n.value, true
}

// Delete 根据 key 删除元素，返回被删除的值和是否存在
func (m *IndexedMap[K, V]) Delete(key K) (V, bool) {
	n, ok := m.keymap[key]
	if !ok {
		var zero V
		return zero, false
	}
	m.removeNode(n)
	delete(m.keymap, key)
	return n.value, true
}

// removeNode 从跳表中移除节点，并收缩最高层
func (m *IndexedMap[K, V]) removeNode(n *node[V]) {
	for i := range len(n.next) {
		prev := n.prev[i]
		next := n.next[i]
		prev.next[i] = next
		if next != nil {
			next.prev[i] = prev
		}
	}
	m.size--

	// 收缩最高层：删除的节点可能是某一层的唯一节点
	for m.level > 1 && m.head.next[m.level-1] == nil {
		m.level--
	}

	n.next = nil
	n.prev = nil
}

// First 返回第一个满足条件 f 的元素，未找到时返回零值和 false
func (m *IndexedMap[K, V]) First(f func(V) bool) (V, bool) {
	for n := m.head.next[0]; n != nil; n = n.next[0] {
		if f(n.value) {
			return n.value, true
		}
	}
	var zero V
	return zero, false
}

// Range 按排序顺序遍历所有元素
func (m *IndexedMap[K, V]) Range(fn func(V)) {
	for n := m.head.next[0]; n != nil; n = n.next[0] {
		fn(n.value)
	}
}

// Values 返回排序后的所有值
func (m *IndexedMap[K, V]) Values() []V {
	values := make([]V, 0, m.size)
	for n := m.head.next[0]; n != nil; n = n.next[0] {
		values = append(values, n.value)
	}
	return values
}
