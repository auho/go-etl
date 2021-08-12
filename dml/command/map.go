package command

import "sync"

type SortMap struct {
	mutex sync.Mutex
	m     map[string]string
	keys  []string
	index int
}

func NewSortMap() *SortMap {
	sm := &SortMap{}
	sm.m = make(map[string]string)
	sm.keys = make([]string, 0)
	sm.index = -1

	return sm
}

func (sm *SortMap) Store(k string, v string) {
	if _, ok := sm.m[k]; !ok {
		sm.keys = append(sm.keys, k)
	}

	sm.m[k] = v
}

func (sm *SortMap) Load(k string) string {
	return sm.m[k]
}

func (sm *SortMap) Next() bool {
	sm.mutex.Lock()
	if sm.index >= len(sm.keys) {
		return false
	}

	sm.index += 1

	return true
}

func (sm *SortMap) Scan() string {
	defer sm.mutex.Unlock()

	return sm.Load(sm.keys[sm.index])
}
