package importtodb

import (
	"fmt"

	"github.com/auho/go-etl/v2/insight/model/excel/read"
	"github.com/auho/go-etl/v2/insight/model/table"
	"github.com/auho/go-etl/v2/insight/model/tool/slices"
	simpleDb "github.com/auho/go-simple-db/v2"
)

type ImportToDb struct {
	isRecreateTable      bool  // 是否 recreate table
	isAppendData         bool  // 是否 append data
	isShowSql            bool  // 是否显示 sql
	columnDropDuplicates []int // drop duplicates for column
	resource             resourcer
	db                   *simpleDb.SimpleDB
}

func RunImportToDb(db *simpleDb.SimpleDB, sr resourcer) error {
	e := &ImportToDb{
		isRecreateTable:      sr.GetIsRecreateTable(),
		isAppendData:         sr.GetIsAppendData(),
		isShowSql:            sr.GetIsShowSql(),
		columnDropDuplicates: sr.GetColumnDropDuplicates(),
		resource:             sr,
		db:                   db,
	}

	return e.Import()
}

func (e *ImportToDb) Import() error {
	_table := e.resource.GetTable()

	err := e.buildTable(_table)
	if err != nil {
		return fmt.Errorf("buildTable error; %w", err)
	}

	sheetData, err := e.resource.GetSheetData()
	if err != nil {
		return fmt.Errorf("GetSheetData error; %w", err)
	}

	err = e.importToTable(_table, sheetData)
	if err != nil {
		return fmt.Errorf("importToTable error; %w", err)
	}

	return nil
}

func (e *ImportToDb) buildTable(table table.Tabler) error {
	if e.isShowSql {
		fmt.Println(table.Sql())
	}

	_, err := e.db.GetTableColumns(table.GetTableName())
	if err != nil {
		e.isRecreateTable = true
	} else {
		if e.isRecreateTable {
			err = e.db.Drop(table.GetTableName())
			if err != nil {
				return fmt.Errorf("drop; %w", err)
			}
		}
	}

	if e.isRecreateTable {
		err = table.Build(e.db)
		if err != nil {
			return fmt.Errorf("build error; %w", err)
		}
	}

	return nil
}

func (e *ImportToDb) importToTable(table table.Tabler, sheetData read.SheetDataor) error {
	var err error

	if len(e.columnDropDuplicates) > 0 {
		err = sheetData.HandlerRows(func(rows [][]string) ([][]string, error) {
			rows = slices.SliceSliceDropDuplicates(rows, e.columnDropDuplicates)

			return rows, nil
		})
		if err != nil {
			return fmt.Errorf("drop duplicates error; %w", err)
		}
	}

	if !e.isAppendData {
		err = e.db.Truncate(table.GetTableName())
		if err != nil {
			return fmt.Errorf("truncate error; %w", err)
		}
	}

	err = e.db.BulkInsertFromSliceSlice(table.GetTableName(), e.resource.GetTitles(), sheetData.GetRowsWithAny(), 1000)
	if err != nil {
		return err
	}

	return nil
}
