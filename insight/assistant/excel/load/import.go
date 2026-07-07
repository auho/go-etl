package load

import (
	"context"
	"fmt"

	"github.com/auho/go-etl/v3/insight/assistant/excel/reader"
	"github.com/auho/go-etl/v3/insight/assistant/schema/create"
	"github.com/auho/go-etl/v3/tool/slicex"
)

type Loader struct {
	xlsxPath string
	resource []Resource
	excel    *reader.Excel
}

func RunLoad(xlsxPath string, sr ...Resource) error {
	e := &Loader{
		xlsxPath: xlsxPath,
		resource: sr,
	}

	return e.Import()
}

func (it *Loader) Import() (err error) {
	fmt.Printf("import start[%s]\n", it.xlsxPath)

	it.excel, err = reader.NewExcel(it.xlsxPath)
	if err != nil {
		return fmt.Errorf("NewExcel: %w", err)
	}
	defer func() {
		if closeErr := it.excel.Close(); closeErr != nil && err == nil {
			err = fmt.Errorf("excel.Close: %w", closeErr)
		}
	}()

	for _, resource := range it.resource {
		fmt.Printf("import resource[%s]\n", resource.GetName())
		err = it.importResource(resource)
		if err != nil {
			return fmt.Errorf("importResource[%s]: %w", resource.GetName(), err)
		}

		err = resource.AfterDo(resource)
		if err != nil {
			return fmt.Errorf("AfterDo[%s]: %w", resource.GetName(), err)
		}
	}

	return nil
}

func (it *Loader) importResource(resource Resource) error {
	err := resource.Prepare()
	if err != nil {
		return fmt.Errorf("prepare: %w", err)
	}

	_table := resource.GetTable()

	err = it.buildResourceTable(resource, _table)
	if err != nil {
		return fmt.Errorf("buildResourceTable: %w", err)
	}

	sheetData, err := resource.GetSheetData(it.excel)
	if err != nil {
		return fmt.Errorf("GetSheetData: %w", err)
	}

	err = it.importResourceToTable(resource, _table, sheetData)
	if err != nil {
		return fmt.Errorf("importResourceToTable: %w", err)
	}

	return nil
}

func (it *Loader) buildResourceTable(resource Resource, table create.Tabler) error {
	if resource.GetIsShowSql() {
		fmt.Println(table.SQL())
	}

	// TODO Optimize 合并 recreate 至 table
	isRecreateTable := resource.GetIsRecreateTable()
	_, err := resource.GetDB().GetTableColumns(context.TODO(), table.GetTableName())
	if err != nil {
		isRecreateTable = true
	} else {
		if isRecreateTable {
			err = resource.GetDB().Drop(context.TODO(), table.GetTableName())
			if err != nil {
				return fmt.Errorf("drop: %w", err)
			}
		}
	}

	if isRecreateTable {
		resource.ExecCommand(table.GetCommand())

		err = table.Build()
		if err != nil {
			return fmt.Errorf("build: %w", err)
		}
	}

	return nil
}

func (it *Loader) importResourceToTable(resource Resource, table create.Tabler, sheetData reader.SheetDataReader) error {
	var err error

	if len(resource.GetColumnDropDuplicates()) > 0 {
		err = sheetData.HandleRows(func(rows [][]string) ([][]string, error) {
			rows = slicex.SliceSliceDropDuplicates(rows, resource.GetColumnDropDuplicates())

			return rows, nil
		})
		if err != nil {
			return fmt.Errorf("HandleRows: %w", err)
		}
	}

	if !resource.GetIsAppendData() {
		err = resource.GetDB().Truncate(context.TODO(), table.GetTableName())
		if err != nil {
			return fmt.Errorf("truncate: %w", err)
		}
	}

	err = resource.GetDB().BulkInsertFromSliceSlice(context.TODO(), table.GetTableName(), resource.GetTitlesName(), sheetData.GetRowsWithAny(), resource.GetBatchInsertSize())
	if err != nil {
		return err
	}

	return nil
}
