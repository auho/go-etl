package command

import (
	"strings"
	"sync"
)

type Entity struct {
	key   string
	value string
	flag  string
}

func NewEntity(k string, v string) *Entity {
	return &Entity{k, v, ""}
}

func NewExpressionEntity(k string, v string) *Entity {
	e := NewEntity(k, v)
	e.flag = FlagExpression

	return e
}

func (e *Entity) Get() (string, string) {
	return e.key, e.value
}

func (e *Entity) GetKey() string {
	return e.key
}

func (e *Entity) GetValue() string {
	return e.value
}

func (e *Entity) IsExpression() bool {
	return e.flag == FlagExpression
}

type Entities struct {
	keys    map[string]int // map[entity name]index
	entries []*Entity
}

func NewEntries() *Entities {
	es := &Entities{}
	es.keys = make(map[string]int)
	es.entries = make([]*Entity, 0)

	return es
}

func (es *Entities) Get() []*Entity {
	return es.entries
}

func (es *Entities) AddEntry(k string, v string) {
	es.Add(NewEntity(k, v))
}

func (es *Entities) AddEntryExpression(k string, v string) {
	es.Add(NewExpressionEntity(k, v))
}

func (es *Entities) Add(e *Entity) {
	// 前面 add 覆盖后面 add 的
	// 已存在的值，不会被覆盖
	if k, ok := es.keys[e.key]; ok {
		_ = k
	} else {
		es.keys[e.key] = len(es.entries)
		es.entries = append(es.entries, e)
	}

	//sort.Sort(sortEntries(es.entries))
}

func (es *Entities) Len() int {
	return len(es.entries)
}

type sortEntries []*Entity

func (se sortEntries) Len() int {
	return len(se)
}

func (se sortEntries) Less(i, j int) bool {
	return strings.ToLower(se[i].key) < strings.ToLower(se[j].key)
}

func (se sortEntries) Swap(i, j int) {
	se[j], se[i] = se[i], se[j]
}

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
