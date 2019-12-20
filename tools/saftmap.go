package tools

import (
	"sync"
)

type SaftMap struct {
	lock *sync.RWMutex
	bm   map[interface{}]interface{}
}

func NewSaftMap() *SaftMap {
	return &SaftMap{
		lock: new(sync.RWMutex),
		bm:   make(map[interface{}]interface{}),
	}
}

func (m *SaftMap) Get(k interface{}) interface{} {
	m.lock.RLock()
	defer m.lock.RUnlock()
	if val, ok := m.bm[k]; ok {
		return val
	}
	return nil
}

func (m *SaftMap) Set(k interface{}, v interface{}) bool {
	m.lock.Lock()
	defer m.lock.Unlock()
	if val, ok := m.bm[k]; !ok {
		m.bm[k] = v
	} else if val != v {
		m.bm[k] = v
	} else {
		return false
	}
	return true
}

func (m *SaftMap) Check(k interface{}) bool {
	m.lock.RLock()
	defer m.lock.RUnlock()
	_, ok := m.bm[k]
	return ok
}

func (m *SaftMap) Delete(k interface{}) {
	m.lock.Lock()
	defer m.lock.Unlock()
	delete(m.bm, k)
}

func (m *SaftMap) Items() map[interface{}]interface{} {
	m.lock.RLock()
	defer m.lock.RUnlock()
	r := make(map[interface{}]interface{})
	for k, v := range m.bm {
		r[k] = v
	}
	return r
}

func (m *SaftMap) Count() int {
	m.lock.RLock()
	defer m.lock.RUnlock()
	return len(m.bm)
}
