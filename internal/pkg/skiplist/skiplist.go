package skiplist

import "math/rand/v2"

const maxLevel = 16

type node[K comparable, V any] struct {
	key   K
	value V
	next  []*node[K, V]
}

type SkipList[K comparable, V any] struct {
	head    *node[K, V]
	level   int
	size    int
	compare func(a, b V) int
	keymap  map[K]*node[K, V]
}

func New[K comparable, V any](compare func(a, b V) int) *SkipList[K, V] {
	return &SkipList[K, V]{
		head:    &node[K, V]{next: make([]*node[K, V], maxLevel)},
		compare: compare,
		keymap:  make(map[K]*node[K, V]),
	}
}

func (sl *SkipList[K, V]) randomLevel() int {
	lv := 0
	for lv < maxLevel-1 && rand.Float64() < 0.25 {
		lv++
	}
	return lv
}

func (sl *SkipList[K, V]) Insert(key K, value V) {
	if _, ok := sl.keymap[key]; ok {
		sl.Delete(key)
	}

	update := make([]*node[K, V], maxLevel)
	cur := sl.head

	for i := sl.level; i >= 0; i-- {
		for cur.next[i] != nil && sl.compare(cur.next[i].value, value) < 0 {
			cur = cur.next[i]
		}
		update[i] = cur
	}

	lv := sl.randomLevel()
	if lv > sl.level {
		for i := sl.level + 1; i <= lv; i++ {
			update[i] = sl.head
		}
		sl.level = lv
	}

	newNode := &node[K, V]{
		key:   key,
		value: value,
		next:  make([]*node[K, V], lv+1),
	}

	for i := 0; i <= lv; i++ {
		newNode.next[i] = update[i].next[i]
		update[i].next[i] = newNode
	}

	sl.keymap[key] = newNode
	sl.size++
}

func (sl *SkipList[K, V]) Delete(key K) bool {
	n, ok := sl.keymap[key]
	if !ok {
		return false
	}

	update := make([]*node[K, V], maxLevel)
	cur := sl.head

	for i := sl.level; i >= 0; i-- {
		for cur.next[i] != nil && sl.compare(cur.next[i].value, n.value) < 0 {
			cur = cur.next[i]
		}
		update[i] = cur
	}

	for i := 0; i <= sl.level; i++ {
		if update[i].next[i] != n {
			break
		}
		update[i].next[i] = n.next[i]
	}

	for sl.level > 0 && sl.head.next[sl.level] == nil {
		sl.level--
	}

	delete(sl.keymap, key)
	sl.size--
	return true
}

func (sl *SkipList[K, V]) Search(key K) (V, bool) {
	node, ok := sl.keymap[key]
	if !ok {
		var zero V
		return zero, false
	}
	return node.value, true
}

func (sl *SkipList[K, V]) First(f func(V) bool) (K, V, bool) {
	cur := sl.head.next[0]
	for cur != nil {
		if f(cur.value) {
			return cur.key, cur.value, true
		}
		cur = cur.next[0]
	}
	var zeroK K
	var zeroV V
	return zeroK, zeroV, false
}

func (sl *SkipList[K, V]) Range(f func(key K, value V)) {
	cur := sl.head.next[0]
	for cur != nil {
		f(cur.key, cur.value)
		cur = cur.next[0]
	}
}

func (sl *SkipList[K, V]) Len() int {
	return sl.size
}
