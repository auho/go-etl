package query

import (
	"fmt"
	"time"

	"github.com/auho/go-toolkit/v3/time/stopwatch"
)

type state struct {
	sourceDuration  time.Duration
	datasetDuration time.Duration
	toSheetDuration time.Duration

	queriesDuration time.Duration
	saveDuration    time.Duration
	totalDuration   time.Duration
	amount          int
}

func (s *state) add(ss sqlState) {
	s.sourceDuration += ss.sourceDuration
	s.datasetDuration += ss.datasetDuration
	s.toSheetDuration += ss.toSheetDuration
	s.amount += ss.amount
}

func (s *state) overview() string {
	return fmt.Sprintf("source: %s, dataset: %s, toSheet: %s <= queries: %s, save: %s, total: %s, amount: %d",
		stopwatch.PrettyDuration(s.sourceDuration),
		stopwatch.PrettyDuration(s.datasetDuration),
		stopwatch.PrettyDuration(s.toSheetDuration),
		stopwatch.PrettyDuration(s.queriesDuration),
		stopwatch.PrettyDuration(s.saveDuration),
		stopwatch.PrettyDuration(s.totalDuration),
		s.amount,
	)
}

type sqlState struct {
	sourceDuration  time.Duration
	datasetDuration time.Duration
	toSheetDuration time.Duration
	totalDuration   time.Duration
	amount          int
}

func (ss *sqlState) overview() string {
	return fmt.Sprintf("source: %s, dataset: %s, toSheet: %s, total: %s, amount: %d",
		stopwatch.PrettyDuration(ss.sourceDuration),
		stopwatch.PrettyDuration(ss.datasetDuration),
		stopwatch.PrettyDuration(ss.toSheetDuration),
		stopwatch.PrettyDuration(ss.totalDuration),
		ss.amount,
	)
}
