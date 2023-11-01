package _import

import (
	"fmt"

	"github.com/auho/go-etl/v2/insight/model/excel/read"
	"github.com/auho/go-etl/v2/insight/model/table"
	"github.com/auho/go-etl/v2/insight/model/tool/slices"
	simpleDb "github.com/auho/go-simple-db/v2"
)

type Import struct {
	isRecreateTable      bool  // 是否 recreate table
	isAppendData         bool  // 是否 append data
	isShowSql            bool  // 是否显示 sql
	columnDropDuplicates []int // drop duplicates for column
	sourceor             sourceor
	db                   *simpleDb.SimpleDB
}

func RunImport(db *simpleDb.SimpleDB, sr sourceor) error {
	_i := &Import{
		isRecreateTable:      sr.GetIsRecreateTable(),
		isAppendData:         sr.GetIsAppendData(),
		isShowSql:            sr.GetIsShowSql(),
		columnDropDuplicates: sr.GetColumnDropDuplicates(),
		sourceor:             sr,
		db:                   db,
	}

	return _i.Import()
}

func (_i *Import) Import() error {
	_table := _i.sourceor.GetTable()

	err := _i.buildTable(_table)
	if err != nil {
		return fmt.Errorf("buildTable error; %w", err)
	}

	sheetData, err := _i.sourceor.GetSheetData()
	if err != nil {
		return fmt.Errorf("GetSheetData error; %w", err)
	}

	err = _i.importToTable(_table, sheetData)
	if err != nil {
		return fmt.Errorf("importToTable error; %w", err)
	}

	return nil
}

func (_i *Import) buildTable(table table.Tabler) error {
	if _i.isShowSql {
		fmt.Println(table.GetTable().SqlForCreate())
	}

	_, err := _i.db.GetTableColumns(table.GetTableName())
	if err != nil {
		_i.isRecreateTable = true
	} else {
		if _i.isRecreateTable {
			err = _i.db.Drop(table.GetTableName())
			if err != nil {
				return fmt.Errorf("drop; %w", err)
			}
		}
	}

	if _i.isRecreateTable {
		err = table.Build(_i.db)
		if err != nil {
			return fmt.Errorf("build error; %w", err)
		}
	}

	return nil
}

func (_i *Import) importToTable(table table.Tabler, sheetData read.SheetDataor) error {
	var err error

	if len(_i.columnDropDuplicates) > 0 {
		err = sheetData.HandlerRows(func(rows [][]string) ([][]string, error) {
			rows = slices.SliceSliceDropDuplicates(rows, _i.columnDropDuplicates)

			return rows, nil
		})
		if err != nil {
			return fmt.Errorf("drop duplicates error; %w", err)
		}
	}

	if !_i.isAppendData {
		err = _i.db.Truncate(table.GetTableName())
		if err != nil {
			return fmt.Errorf("truncate error; %w", err)
		}
	}

	err = _i.db.BulkInsertFromSliceSlice(table.GetTableName(), _i.sourceor.GetTitles(), sheetData.GetRowsWithAny(), 1000)
	if err != nil {
		return err
	}

	return nil
}
