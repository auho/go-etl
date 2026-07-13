package load

import (
	"context"
	"fmt"

	"github.com/auho/go-etl/v3/insight/assistant/excel/reader"
	"github.com/auho/go-etl/v3/insight/assistant/schema/create"
	"github.com/auho/go-etl/v3/tool/slicex"
)

// RunLoad loads one or more resources from an xlsx file into the database.
func RunLoad(xlsxPath string, r ...Resource) error {
	l := &Loader{
		xlsxPath:  xlsxPath,
		resources: r,
	}

	return l.Import()
}

type Loader struct {
	excel     *reader.Excel
	xlsxPath  string
	resources []Resource
}

func (l *Loader) Import() (err error) {
	fmt.Printf("import start[%s]\n", l.xlsxPath)

	l.excel, err = reader.NewExcel(l.xlsxPath)
	if err != nil {
		return fmt.Errorf("NewExcel: %w", err)
	}
	defer func() {
		if closeErr := l.excel.Close(); closeErr != nil && err == nil {
			err = fmt.Errorf("excel.Close: %w", closeErr)
		}
	}()

	for _, resource := range l.resources {
		fmt.Printf("import resource[%s]\n", resource.Name())
		err = l.importResource(resource)
		if err != nil {
			return fmt.Errorf("importResource[%s]: %w", resource.Name(), err)
		}

		err = resource.AfterDo(resource)
		if err != nil {
			return fmt.Errorf("AfterDo[%s]: %w", resource.Name(), err)
		}
	}

	return nil
}

func (l *Loader) importResource(resource Resource) error {
	err := resource.Prepare()
	if err != nil {
		return fmt.Errorf("prepare: %w", err)
	}

	_table := resource.Tabler()

	err = l.buildResourceTable(resource, _table)
	if err != nil {
		return fmt.Errorf("buildResourceTable: %w", err)
	}

	sheetData, err := resource.SheetData(l.excel)
	if err != nil {
		return fmt.Errorf("GetSheetData: %w", err)
	}

	err = l.importResourceToTable(resource, _table, sheetData)
	if err != nil {
		return fmt.Errorf("importResourceToTable: %w", err)
	}

	return nil
}

func (l *Loader) buildResourceTable(resource Resource, table create.Tabler) error {
	if resource.GetIsShowSql() {
		fmt.Println(table.SQL())
	}

	// TODO Optimize: merge recreate logic into table
	isRecreateTable := resource.GetIsRecreateTable()
	_, err := resource.DB().GetTableColumns(context.TODO(), table.GetTableName())
	if err != nil {
		isRecreateTable = true
	} else {
		if isRecreateTable {
			err = resource.DB().Drop(context.TODO(), table.GetTableName())
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

func (l *Loader) importResourceToTable(resource Resource, table create.Tabler, sheetData reader.SheetDataReader) error {
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
		err = resource.DB().Truncate(context.TODO(), table.GetTableName())
		if err != nil {
			return fmt.Errorf("truncate: %w", err)
		}
	}

	err = resource.DB().BulkInsertFromSliceSlice(context.TODO(), table.GetTableName(), resource.TitlesName(), sheetData.GetRowsWithAny(), resource.GetBatchInsertSize())
	if err != nil {
		return err
	}

	return nil
}
