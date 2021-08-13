package command

import "sync"

const flagAggregation = "aggregation"

type SortMap struct {
	RWMutex sync.RWMutex
	m       map[string]string
	keys    []string
	index   int
}

func NewSortMap() *SortMap {
	sm := &SortMap{}
	sm.init()

	return sm
}

func (sm *SortMap) init() {
	sm.m = make(map[string]string)
	sm.keys = make([]string, 0)
	sm.index = -1
}

func (sm *SortMap) Store(k string, v string) {
	sm.RWMutex.Lock()
	defer sm.RWMutex.Unlock()

	if _, ok := sm.m[k]; !ok {
		sm.keys = append(sm.keys, k)
	}

	sm.m[k] = v
}

func (sm *SortMap) Load(k string) string {
	sm.RWMutex.RLock()
	defer sm.RWMutex.RUnlock()

	return sm.m[k]
}

func (sm *SortMap) Next() bool {
	sm.RWMutex.RLock()
	if sm.index+1 >= len(sm.keys) {
		sm.index = -1
		return false
	}

	sm.index += 1

	return true
}

func (sm *SortMap) Scan() (string, string) {
	// TODO
	defer sm.RWMutex.RUnlock()
	key := sm.keys[sm.index]

	return key, sm.Load(key)
}

func (sm *SortMap) Len() int {
	return len(sm.m)
}

type Entry struct {
	key   string
	value string
	flag  string
}

func NewEntry(k string, v string) *Entry {
	return &Entry{k, v, ""}
}

func NewAggregationEntry(k string, v string) *Entry {
	e := NewEntry(k, v)
	e.flag = flagAggregation

	return e
}

func (e *Entry) Get() (string, string) {
	return e.key, e.value
}

func (e *Entry) GetKey() string {
	return e.key
}

func (e *Entry) GetValue() string {
	return e.value
}

func (e *Entry) IsAggregation() bool {
	return e.flag == flagAggregation
}

type Entries struct {
	keys    map[string]int
	entries []*Entry
}

func NewEntries() *Entries {
	es := &Entries{}
	es.keys = make(map[string]int)
	es.entries = make([]*Entry, 0)

	return es
}

func (es *Entries) Get() []*Entry {
	return es.entries
}

func (es *Entries) AddEntry(k string, v string) {
	es.Add(NewEntry(k, v))
}

func (es *Entries) Add(e *Entry) {
	if k, ok := es.keys[e.key]; ok {
		es.entries[k] = e
	} else {
		es.keys[e.key] = len(es.entries)
		es.entries = append(es.entries, e)
	}
}

func (es *Entries) Len() int {
	return len(es.entries)
}
