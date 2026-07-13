package pairmap

type entry[K comparable, V any] struct {
	left  K
	right K
	value V
}

type PairMap[K comparable, V any] struct {
	lmap map[K]*entry[K, V]
	rmap map[K]*entry[K, V]
}

func New[K comparable, V any]() *PairMap[K, V] {
	return &PairMap[K, V]{
		lmap: make(map[K]*entry[K, V]),
		rmap: make(map[K]*entry[K, V]),
	}
}

func (pm *PairMap[K, V]) Set(left, right K, value V) {
	if e, ok := pm.lmap[left]; ok {
		delete(pm.rmap, e.right)
	}
	if e, ok := pm.rmap[right]; ok {
		delete(pm.lmap, e.left)
	}
	e := &entry[K, V]{left: left, right: right, value: value}
	pm.lmap[left] = e
	pm.rmap[right] = e
}

func (pm *PairMap[K, V]) Get(key K) (V, bool) {
	if e, ok := pm.lmap[key]; ok {
		return e.value, true
	}
	if e, ok := pm.rmap[key]; ok {
		return e.value, true
	}
	var zero V
	return zero, false
}

func (pm *PairMap[K, V]) Delete(key K) bool {
	if e, ok := pm.lmap[key]; ok {
		delete(pm.lmap, key)
		delete(pm.rmap, e.right)
		return true
	}
	if e, ok := pm.rmap[key]; ok {
		delete(pm.rmap, key)
		delete(pm.lmap, e.left)
		return true
	}
	return false
}

func (pm *PairMap[K, V]) Len() int {
	return len(pm.lmap)
}
