package dataset

import (
	"time"
)

// QueryResult represents the execution result of a single SQL query.
type QueryResult struct {
	Amount   int
	Duration time.Duration
	Name     string
	SQL      string
}

// Subset represents a subset of data within a Dataset.
type Subset struct {
	Amount   int
	Duration time.Duration
	Name     string
	Rows     [][]any
	Queries  []QueryResult
}

// NewSubsetWithQuery creates a Subset from a single query result.
func NewSubsetWithQuery(name string, sql string, d time.Duration, rows [][]any) Subset {
	q := QueryResult{
		Amount:   len(rows),
		Duration: d,
		Name:     name,
		SQL:      sql,
	}

	s := Subset{Name: name}
	s.AddQuery(q)
	s.Rows = rows

	return s
}

// NewSubsetFromSubsets creates a Subset by merging multiple Subsets.
func NewSubsetFromSubsets(name string, ss []Subset) Subset {
	s := Subset{Name: name}

	for _, sub := range ss {
		s.AddSubset(sub)
	}

	return s
}

func (s *Subset) AddSubset(sub Subset) {
	s.Rows = append(s.Rows, sub.Rows...)

	s.AddQueries(sub.Queries)
}

func (s *Subset) AddQueries(qs []QueryResult) {
	for _, q := range qs {
		s.AddQuery(q)
	}
}

func (s *Subset) AddQuery(q QueryResult) {
	s.Amount += q.Amount
	s.Duration += q.Duration
	s.Queries = append(s.Queries, q)
}
