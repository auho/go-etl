package query

import (
	"fmt"
	"time"

	"github.com/auho/go-toolkit/time/timing"
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
	return fmt.Sprintf("source: %s, dateset: %s, toSheet: %s <= queries: %s, save: %s, total: %s, amount: %d",
		timing.PrettyDuration(s.sourceDuration),
		timing.PrettyDuration(s.datasetDuration),
		timing.PrettyDuration(s.toSheetDuration),
		timing.PrettyDuration(s.queriesDuration),
		timing.PrettyDuration(s.saveDuration),
		timing.PrettyDuration(s.totalDuration),
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
	return fmt.Sprintf("source: %s, dateset: %s, toSheet: %s, total: %s, amount: %d",
		timing.PrettyDuration(ss.sourceDuration),
		timing.PrettyDuration(ss.datasetDuration),
		timing.PrettyDuration(ss.toSheetDuration),
		timing.PrettyDuration(ss.totalDuration),
		ss.amount,
	)
}
