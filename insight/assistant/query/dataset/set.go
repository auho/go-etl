package dataset

import (
	"time"
)

type Query struct {
	Amount   int
	Duration time.Duration
	Name     string
	Sql      string
}

type Set struct {
	Amount   int
	Duration time.Duration
	Name     string
	Rows     [][]any
	Queries  []Query
}

func NewSetWithQuery(name string, sql string, d time.Duration, rows [][]any) Set {
	q := Query{
		Amount:   len(rows),
		Duration: d,
		Name:     name,
		Sql:      sql,
	}

	s := Set{Name: name}
	s.AddQuery(q)
	s.Rows = rows

	return s
}

func NewSetWithSets(name string, ss []Set) Set {
	s := Set{Name: name}

	for _, _s := range ss {
		s.AddSet(_s)
	}

	return s
}

func (s *Set) AddSet(_s Set) {
	s.Rows = append(s.Rows, _s.Rows...)

	s.AddQueries(_s.Queries)
}

func (s *Set) AddQueries(qs []Query) {
	for _, q := range qs {
		s.AddQuery(q)
	}
}

func (s *Set) AddQuery(q Query) {
	s.Amount += q.Amount
	s.Duration += q.Duration
	s.Queries = append(s.Queries, q)
}
