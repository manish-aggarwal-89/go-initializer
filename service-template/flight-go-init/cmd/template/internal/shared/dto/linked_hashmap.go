package dto

type LinkedHashMap[V any] struct {
	keys   []int
	values map[int]V
}

func NewLinkedHashMap[V any]() *LinkedHashMap[V] {
	return &LinkedHashMap[V]{
		keys:   []int{},
		values: make(map[int]V),
	}
}

func (lhm *LinkedHashMap[V]) Put(key int, value V) {
	if _, exists := lhm.values[key]; !exists {
		lhm.keys = append(lhm.keys, key)
	}
	lhm.values[key] = value
}

func (lhm *LinkedHashMap[V]) Get(key int) (V, bool) {
	value, exists := lhm.values[key]
	return value, exists
}

func (lhm *LinkedHashMap[V]) GetOrDefault(key int, defaultVal V) V {
	value, exists := lhm.values[key]
	if !exists {
		return defaultVal
	}
	return value
}

func (lhm *LinkedHashMap[V]) Remove(key int) {
	if _, exists := lhm.values[key]; exists {
		delete(lhm.values, key)
		for i, k := range lhm.keys {
			if k == key {
				lhm.keys = append(lhm.keys[:i], lhm.keys[i+1:]...)
				break
			}
		}
	}
}

func (lhm *LinkedHashMap[V]) Keys() []int {
	return lhm.keys
}

func (lhm *LinkedHashMap[V]) Values() []V {
	values := make([]V, len(lhm.keys))
	for i, key := range lhm.keys {
		values[i] = lhm.values[key]
	}
	return values
}

func (lhm *LinkedHashMap[V]) GetMap() map[int]V {
	return lhm.values
}
