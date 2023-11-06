package importtodb

import (
	"fmt"

	"github.com/auho/go-etl/v2/insight/assistant/excel/read"
	"github.com/auho/go-etl/v2/insight/assistant/tablestructure/buildtable"
	"github.com/auho/go-etl/v2/tool/slices"
)

type ImportToDb struct {
	xlsxPath string
	resource []resourcer
	excel    *read.Excel
}

func RunImportToDb(xlsxPath string, sr ...resourcer) error {
	e := &ImportToDb{
		xlsxPath: xlsxPath,
		resource: sr,
	}

	return e.Import()
}

func (it *ImportToDb) Import() error {
	var err error
	it.excel, err = read.NewExcel(it.xlsxPath)
	if err != nil {
		return fmt.Errorf("NewExcel error; %w", err)
	}

	for _, resource := range it.resource {
		err = it.importResource(resource)
		if err != nil {
			return fmt.Errorf("importResource error; %w", err)
		}
	}

	return nil
}

func (it *ImportToDb) importResource(resource resourcer) error {
	err := resource.Prepare()
	if err != nil {
		return fmt.Errorf("prepare error; %w", err)
	}

	_table := resource.GetTable()

	resource.CommandExec(_table.GetCommand())

	err = it.buildResourceTable(resource, _table)
	if err != nil {
		return fmt.Errorf("buildResourceTable error; %w", err)
	}

	sheetData, err := resource.GetSheetData(it.excel)
	if err != nil {
		return fmt.Errorf("GetSheetData error; %w", err)
	}

	err = it.importResourceToTable(resource, _table, sheetData)
	if err != nil {
		return fmt.Errorf("importResourceToTable error; %w", err)
	}

	return nil
}

func (it *ImportToDb) buildResourceTable(resource resourcer, table buildtable.Tabler) error {
	if resource.GetIsShowSql() {
		fmt.Println(table.Sql())
	}

	isRecreateTable := resource.GetIsRecreateTable()
	_, err := resource.GetDB().GetTableColumns(table.GetTableName())
	if err != nil {
		isRecreateTable = true
	} else {
		if isRecreateTable {
			err = resource.GetDB().Drop(table.GetTableName())
			if err != nil {
				return fmt.Errorf("drop; %w", err)
			}
		}
	}

	if isRecreateTable {
		err = table.Build(resource.GetDB())
		if err != nil {
			return fmt.Errorf("build error; %w", err)
		}
	}

	return nil
}

func (it *ImportToDb) importResourceToTable(resource resourcer, table buildtable.Tabler, sheetData read.SheetDataor) error {
	var err error

	if len(resource.GetColumnDropDuplicates()) > 0 {
		err = sheetData.HandlerRows(func(rows [][]string) ([][]string, error) {
			rows = slices.SliceSliceDropDuplicates(rows, resource.GetColumnDropDuplicates())

			return rows, nil
		})
		if err != nil {
			return fmt.Errorf("drop duplicates error; %w", err)
		}
	}

	if !resource.GetIsAppendData() {
		err = resource.GetDB().Truncate(table.GetTableName())
		if err != nil {
			return fmt.Errorf("truncate error; %w", err)
		}
	}

	err = resource.GetDB().BulkInsertFromSliceSlice(table.GetTableName(), resource.GetTitlesKey(), sheetData.GetRowsWithAny(), 1000)
	if err != nil {
		return err
	}

	return nil
}
