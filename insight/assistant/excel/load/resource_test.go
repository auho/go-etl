package load

import (
	"errors"
	"testing"

	"github.com/auho/go-etl/v3/insight/assistant/schema"
)

func TestGetBatchInsertSize_Default(t *testing.T) {
	b := &baseResource{BatchInsertSize: 0}
	if got := b.GetBatchInsertSize(); got != 2000 {
		t.Errorf("GetBatchInsertSize() = %d, want 2000", got)
	}
}

func TestGetBatchInsertSize_Negative(t *testing.T) {
	b := &baseResource{BatchInsertSize: -1}
	if got := b.GetBatchInsertSize(); got != 2000 {
		t.Errorf("GetBatchInsertSize() = %d, want 2000", got)
	}
}

func TestGetBatchInsertSize_Custom(t *testing.T) {
	b := &baseResource{BatchInsertSize: 500}
	if got := b.GetBatchInsertSize(); got != 500 {
		t.Errorf("GetBatchInsertSize() = %d, want 500", got)
	}
}

func TestGetIsRecreateTable(t *testing.T) {
	b := &baseResource{IsRecreateTable: true}
	if !b.GetIsRecreateTable() {
		t.Error("GetIsRecreateTable() = false, want true")
	}

	b.IsRecreateTable = false
	if b.GetIsRecreateTable() {
		t.Error("GetIsRecreateTable() = true, want false")
	}
}

func TestGetIsAppendData(t *testing.T) {
	b := &baseResource{IsAppendData: true}
	if !b.GetIsAppendData() {
		t.Error("GetIsAppendData() = false, want true")
	}

	b.IsAppendData = false
	if b.GetIsAppendData() {
		t.Error("GetIsAppendData() = true, want false")
	}
}

func TestGetIsShowSql(t *testing.T) {
	b := &baseResource{IsShowSql: true}
	if !b.GetIsShowSql() {
		t.Error("GetIsShowSql() = false, want true")
	}

	b.IsShowSql = false
	if b.GetIsShowSql() {
		t.Error("GetIsShowSql() = true, want false")
	}
}

func TestGetColumnDropDuplicates(t *testing.T) {
	cols := []int{1, 3, 5}
	b := &baseResource{columnDropDuplicates: cols}
	got := b.GetColumnDropDuplicates()
	if len(got) != len(cols) {
		t.Fatalf("GetColumnDropDuplicates() len = %d, want %d", len(got), len(cols))
	}
	for i, v := range cols {
		if got[i] != v {
			t.Errorf("GetColumnDropDuplicates()[%d] = %d, want %d", i, got[i], v)
		}
	}
}

func TestGetColumnDropDuplicates_Empty(t *testing.T) {
	b := &baseResource{}
	if got := b.GetColumnDropDuplicates(); len(got) != 0 {
		t.Errorf("GetColumnDropDuplicates() = %v, want empty", got)
	}
}

func TestExecCommand_Nil(t *testing.T) {
	b := &baseResource{}
	b.ExecCommand(schema.NewCommandMysql())
}

func TestExecCommand_WithFun(t *testing.T) {
	called := false
	b := &baseResource{
		CommandFun: func(c *schema.Command) {
			called = true
			c.AddString("test_col")
		},
	}
	b.ExecCommand(schema.NewCommandMysql())
	if !called {
		t.Error("CommandFun was not called")
	}
}

func TestAfterDo_Nil(t *testing.T) {
	b := &baseResource{}
	err := b.AfterDo(nil)
	if err != nil {
		t.Errorf("AfterDo() error = %v, want nil", err)
	}
}

func TestAfterDo_WithFun(t *testing.T) {
	called := false
	b := &baseResource{
		AfterFun: func(r Resource) error {
			called = true
			return nil
		},
	}
	err := b.AfterDo(nil)
	if err != nil {
		t.Errorf("AfterDo() error = %v, want nil", err)
	}
	if !called {
		t.Error("AfterFun was not called")
	}
}

func TestAfterDo_WithError(t *testing.T) {
	expectedErr := errors.New("after error")
	b := &baseResource{
		AfterFun: func(r Resource) error {
			return expectedErr
		},
	}
	err := b.AfterDo(nil)
	if !errors.Is(err, expectedErr) {
		t.Errorf("AfterDo() error = %v, want %v", err, expectedErr)
	}
}

func TestBuildSheetConfig(t *testing.T) {
	b := &baseResource{
		SheetName:  "Sheet1",
		SheetIndex: 2,
		StartRow:   3,
		EndRow:     10,
	}
	cfg := b.buildSheetConfig()
	if cfg.SheetName != "Sheet1" {
		t.Errorf("SheetName = %q, want %q", cfg.SheetName, "Sheet1")
	}
	if cfg.SheetIndex != 2 {
		t.Errorf("SheetIndex = %d, want %d", cfg.SheetIndex, 2)
	}
	if cfg.StartRow != 3 {
		t.Errorf("StartRow = %d, want %d", cfg.StartRow, 3)
	}
	if cfg.EndRow != 10 {
		t.Errorf("EndRow = %d, want %d", cfg.EndRow, 10)
	}
}

func TestBuildSheetConfig_Empty(t *testing.T) {
	b := &baseResource{}
	cfg := b.buildSheetConfig()
	if cfg.SheetName != "" {
		t.Errorf("SheetName = %q, want empty", cfg.SheetName)
	}
	if cfg.SheetIndex != 0 {
		t.Errorf("SheetIndex = %d, want 0", cfg.SheetIndex)
	}
	if cfg.StartRow != 0 {
		t.Errorf("StartRow = %d, want 0", cfg.StartRow)
	}
	if cfg.EndRow != 0 {
		t.Errorf("EndRow = %d, want 0", cfg.EndRow)
	}
}
