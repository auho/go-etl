package task

import (
	"fmt"
	"testing"

	"github.com/auho/go-etl/v3/job/extract/tag"
	"github.com/auho/go-etl/v3/job/transform"
	"github.com/auho/go-etl/v3/job/transform/collector"
)

// Expected-counts mapping lives with the data: see setup_test.go (baseNames).
// Helpers below derive counts from _maxA/_maxB driven by that mapping.

func baseMultiplier() int64 {
	return int64(_maxA / 10)
}

func expectedRows(patternCount int) int64 {
	return int64(patternCount) * baseMultiplier() * int64(_maxB)
}

func totalRows() int64 {
	return int64(_maxA * _maxB)
}

func sourceConfig() SourceConfig {
	return SourceConfig{
		PageSize: _pageSize,
	}
}

func Test_Update(t *testing.T) {
	src := newTestSource(t, "update")
	m := transform.NewUpdate(collector.NewKeysAll([]string{_keyName}, tag.NewMostKey(_rule)), nil)
	ua := NewUpdate(src, []transform.UpdateOperator{m})

	err := RunProducer(src, []itemProducer{ua}, WithRunnerSourceConfig(sourceConfig()))
	if err != nil {
		t.Error(err)
	}

	tbl := src.TableName()
	total := totalRows()
	matched := expectedRows(2 + 2 + 1 + 1 + 1)
	noMatch := expectedRows(3)

	gotTotal := getAmount(tbl, t)
	if gotTotal != total {
		t.Errorf("total: got %d, want %d", gotTotal, total)
	}

	gotMatched := getFieldAmount(tbl, "a !=", "", t)
	if gotMatched != matched {
		t.Errorf("matched (a != ''): got %d, want %d", gotMatched, matched)
	}

	gotNoMatch := getFieldAmount(tbl, "a =", "", t)
	if gotNoMatch != noMatch {
		t.Errorf("unmatched (a == ''): got %d, want %d", gotNoMatch, noMatch)
	}

	keywordCases := []struct {
		keyword string
		amount  int
		aVal    string
		abVal   string
		base    int
	}{
		{"a", 2, "a", "a1", 2},
		{"中文", 3, "中文", "中文1", 2},
		{"b", 3, "a", "a1", 1},
		{"123", 2, "123", "123", 1},
		{"ab", 3, "ab", "ab1", 1},
	}

	for _, tc := range keywordCases {
		got := countByKeywordAmount(tbl, tc.keyword, tc.amount, t)
		want := expectedRows(tc.base)
		if got != want {
			t.Errorf("keyword=%s amount=%d: got %d, want %d", tc.keyword, tc.amount, got, want)
		}

		gotA := countByTriple(tbl, tc.keyword, tc.amount, "a", tc.aVal, t)
		if gotA != want {
			t.Errorf("keyword=%s a=%s: got %d, want %d", tc.keyword, tc.aVal, gotA, want)
		}

		gotAB := countByTriple(tbl, tc.keyword, tc.amount, "ab", tc.abVal, t)
		if gotAB != want {
			t.Errorf("keyword=%s ab=%s: got %d, want %d", tc.keyword, tc.abVal, gotAB, want)
		}
	}
}

func Test_UpdateTransfer(t *testing.T) {
	src := newTestSource(t, "ut_src")
	m := transform.NewUpdate(collector.NewKeysAll([]string{_keyName}, tag.NewMostKey(_rule)), nil)
	ut := NewUpdateTransfer(src, _targetUpdateTransfer, []transform.UpdateOperator{m})

	err := RunProducer(src, []itemProducer{ut}, WithRunnerSourceConfig(sourceConfig()))
	if err != nil {
		t.Error(err)
	}

	matched := expectedRows(2 + 2 + 1 + 1 + 1)

	gotTransfer := getAmount(_updateAndTransferTable, t)
	if gotTransfer != matched {
		t.Errorf("transferCount: got %d, want %d", gotTransfer, matched)
	}

	keywordCases := []struct {
		keyword string
		amount  int
		base    int
	}{
		{"a", 2, 2},
		{"中文", 3, 2},
		{"b", 3, 1},
		{"123", 2, 1},
		{"ab", 3, 1},
	}

	for _, tc := range keywordCases {
		got := countByKeywordAmount(_updateAndTransferTable, tc.keyword, tc.amount, t)
		want := expectedRows(tc.base)
		if got != want {
			t.Errorf("keyword=%s amount=%d: got %d, want %d", tc.keyword, tc.amount, got, want)
		}
	}
}

func Test_Insert(t *testing.T) {
	src := newTestSource(t, "insert")

	err := _gormDB.Exec(fmt.Sprintf("TRUNCATE TABLE `%s`", _targetTagA.TableName())).Error
	if err != nil {
		t.Error(err)
	}

	insertConfig := WithInsertConfig(InsertConfig{
		SkipTruncate: true,
		ExtraKeys:    []string{src.IDName()},
	})

	m := transform.NewInsert(collector.NewKeysAll([]string{_keyName}, tag.NewKey(_rule)), nil)
	ia := NewInsert(_targetTagA, m, insertConfig)

	err = RunProducer(src, []itemProducer{ia}, WithRunnerSourceConfig(sourceConfig()))
	if err != nil {
		t.Error(err)
	}

	// ScanKey rows per base pattern:
	//   #0,#1 (a一a一...): 5 keywords (a,b,ab,123,中文) x 2 rows = 10
	//   #2,#3 (中文...):   1 keyword  (中文)          x 2 rows = 2
	//   #7    (b一b...):   2 keywords (b,a)           x 1 row  = 2
	//   #8    (123...):    2 keywords (123,ab)        x 1 row  = 2
	//   #9    (ab...):     1 keyword  (ab)            x 1 row  = 1
	// Total ScanKey rows = 17 base rows per round
	scanKeyBaseTotal := int64(17)
	perTableExpected := scanKeyBaseTotal * baseMultiplier() * int64(_maxB)

	got := getAmount(_targetTagA.TableName(), t)
	if got != perTableExpected {
		t.Errorf("%s: got %d, want %d", _targetTagA.TableName(), got, perTableExpected)
	}

	keywordCases := []struct {
		keyword string
		amount  int
		base    int
	}{
		{"a", 2, 2},
		{"a", 1, 1},
		{"b", 1, 2},
		{"b", 3, 1},
		{"ab", 1, 3},
		{"ab", 3, 1},
		{"123", 1, 2},
		{"123", 2, 1},
		{"中文", 1, 2},
		{"中文", 3, 2},
	}

	for _, tc := range keywordCases {
		got = countByKeywordAmount(_targetTagA.TableName(), tc.keyword, tc.amount, t)
		want := expectedRows(tc.base)
		if got != want {
			t.Errorf("keyword=%s amount=%d: got %d, want %d", tc.keyword, tc.amount, got, want)
		}
	}

	aValCases := []struct {
		keyword string
		amount  int
		aVal    string
		abVal   string
		base    int
	}{
		{"a", 2, "a", "a1", 2},
		{"b", 3, "a", "a1", 1},
		{"ab", 3, "ab", "ab1", 1},
		{"123", 2, "123", "123", 1},
		{"中文", 3, "中文", "中文1", 2},
	}

	for _, tc := range aValCases {
		gotA := countByTriple(_targetTagA.TableName(), tc.keyword, tc.amount, "a", tc.aVal, t)
		want := expectedRows(tc.base)
		if gotA != want {
			t.Errorf("keyword=%s a=%s: got %d, want %d", tc.keyword, tc.aVal, gotA, want)
		}

		gotAB := countByTriple(_targetTagA.TableName(), tc.keyword, tc.amount, "ab", tc.abVal, t)
		if gotAB != want {
			t.Errorf("keyword=%s ab=%s: got %d, want %d", tc.keyword, tc.abVal, gotAB, want)
		}
	}

	gotDidNotNull := getFieldAmount(_targetTagA.TableName(), "did IS NOT", nil, t)
	if gotDidNotNull != perTableExpected {
		t.Errorf("did not null: got %d, want %d", gotDidNotNull, perTableExpected)
	}
}

func Test_Transfer(t *testing.T) {
	src := newTestSource(t, "tf_src")
	keys := []string{"did", "name", "a", "ab", "a_keyword", "a_keyword_num", "a_keyword_amount"}
	alias := map[string]string{
		"did":           "did",
		"name":          "name",
		"a":             "a1",
		"ab":            "ab1",
		"a_keyword":     "a_keyword",
		"a_keyword_num": "a_keyword_num",
	}

	m := transform.NewTransfer(keys, alias, map[string]any{"xyz": "xyz1"})
	tf := NewTransfer(_targetTransfer, m)

	err := RunProducer(src, []itemProducer{tf}, WithRunnerSourceConfig(sourceConfig()))
	if err != nil {
		t.Error(err)
	}

	total := totalRows()

	gotTransfer := getAmount(_targetTransfer.TableName(), t)
	if gotTransfer != total {
		t.Errorf("transferCount: got %d, want %d", gotTransfer, total)
	}

	gotXyz := getFieldAmount(_targetTransfer.TableName(), "xyz =", "xyz1", t)
	if gotXyz != total {
		t.Errorf("xyz=xyz1: got %d, want %d", gotXyz, total)
	}

	// alias mapping: a -> a1, ab -> ab1 (source a/ab are "" initially)
	gotA1Empty := getFieldAmount(_targetTransfer.TableName(), "a1 =", "", t)
	if gotA1Empty != total {
		t.Errorf("a1='' (alias of a): got %d, want %d", gotA1Empty, total)
	}

	gotAB1Empty := getFieldAmount(_targetTransfer.TableName(), "ab1 =", "", t)
	if gotAB1Empty != total {
		t.Errorf("ab1='' (alias of ab): got %d, want %d", gotAB1Empty, total)
	}

	gotKeywordEmpty := getFieldAmount(_targetTransfer.TableName(), "a_keyword =", "", t)
	if gotKeywordEmpty != total {
		t.Errorf("a_keyword='' (mapped through): got %d, want %d", gotKeywordEmpty, total)
	}

	gotDidNotNull := getFieldAmount(_targetTransfer.TableName(), "did IS NOT", nil, t)
	if gotDidNotNull != total {
		t.Errorf("did not null: got %d, want %d", gotDidNotNull, total)
	}
}

func Test_Clean(t *testing.T) {
	src := newTestSource(t, "clean")
	resource := &targetCleanTest{source: src}
	m := transform.NewUpdate(collector.NewKeysAll([]string{_keyName}, tag.NewMostKey(_rule)), nil)

	clean := NewClean(resource, []transform.UpdateOperator{m})
	err := RunConsumer(resource.Source(), []itemConsumer{clean}, WithRunnerSourceConfig(sourceConfig()))
	if err != nil {
		t.Error(err)
	}

	total := totalRows()
	matched := expectedRows(2 + 2 + 1 + 1 + 1)
	noMatch := expectedRows(3)

	gotClean := getAmount(resource.Data().TableName(), t)
	if gotClean != noMatch {
		t.Errorf("cleanData: got %d, want %d", gotClean, noMatch)
	}

	gotCleanEmpty := getFieldAmount(resource.Data().TableName(), "a =", "", t)
	if gotCleanEmpty != noMatch {
		t.Errorf("cleanData a='': got %d, want %d", gotCleanEmpty, noMatch)
	}

	gotDeleted := getAmount(resource.Deleted().TableName(), t)
	if gotDeleted != matched {
		t.Errorf("deletedData: got %d, want %d", gotDeleted, matched)
	}

	if gotClean+gotDeleted != total {
		t.Errorf("conservation: clean[%d] + deleted[%d] != total[%d]", gotClean, gotDeleted, total)
	}

	gotDeletedEmpty := getFieldAmount(resource.Deleted().TableName(), "a =", "", t)
	if gotDeletedEmpty != matched {
		t.Errorf("deletedData a='': got %d, want %d", gotDeletedEmpty, matched)
	}

	nameCases := []struct {
		name string
		base int
	}{
		{"a一a一b一ab一123一中文", 2},
		{"中文一中文一中文", 2},
		{"b一b一b一a", 1},
		{"123一123一ab", 1},
		{"ab一ab一ab", 1},
	}

	for _, tc := range nameCases {
		got := getFieldAmount(resource.Deleted().TableName(), "name =", tc.name, t)
		want := expectedRows(tc.base)
		if got != want {
			t.Errorf("deleted name=%s: got %d, want %d", tc.name, got, want)
		}
	}

	gotNoMatchName := getFieldAmount(resource.Data().TableName(), "name =", "xyz_no_match", t)
	if gotNoMatchName != noMatch {
		t.Errorf("cleanData name=xyz_no_match: got %d, want %d", gotNoMatchName, noMatch)
	}
}

// --- helpers ---

func getAmount(tableName string, t *testing.T) int64 {
	var count int64
	err := _gormDB.Table(tableName).Count(&count).Error
	if err != nil {
		t.Error(err)
	}

	return count
}

func getFieldAmount(tableName string, op string, value any, t *testing.T) int64 {
	var count int64
	var err error
	if value == nil {
		err = _gormDB.Table(tableName).Where(fmt.Sprintf("%s NULL", op)).Count(&count).Error
	} else {
		err = _gormDB.Table(tableName).Where(fmt.Sprintf("%s ?", op), value).Count(&count).Error
	}
	if err != nil {
		t.Error(err)
	}

	return count
}

func countByKeywordAmount(tableName string, keyword string, amount int, t *testing.T) int64 {
	var count int64
	err := _gormDB.Table(tableName).
		Where("a_keyword = ?", keyword).
		Where("a_keyword_amount = ?", amount).
		Count(&count).Error
	if err != nil {
		t.Error(err)
	}

	return count
}

func countByTriple(tableName string, keyword string, amount int, field string, value string, t *testing.T) int64 {
	var count int64
	err := _gormDB.Table(tableName).
		Where("a_keyword = ?", keyword).
		Where("a_keyword_amount = ?", amount).
		Where(fmt.Sprintf("%s = ?", field), value).
		Count(&count).Error
	if err != nil {
		t.Error(err)
	}

	return count
}
