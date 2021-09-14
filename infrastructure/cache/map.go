package cache

import "sync"

type key interface{}
type value interface{}

type concurrentMap struct {
	mu sync.RWMutex

	m map[key]value
}

func newConcurrentMap() *concurrentMap {
	m := make(map[key]value)
	return &concurrentMap{
		m: m,
	}
}

func (m *concurrentMap) Get(k key) (v value, ok bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	v, ok = m.m[k]

	return
}

func (m *concurrentMap) Put(k key, v value) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.m[k] = v
}

func (m *concurrentMap) Delete(k key) {
	m.mu.Lock()
	defer m.mu.Unlock()

	delete(m.m, k)
}
