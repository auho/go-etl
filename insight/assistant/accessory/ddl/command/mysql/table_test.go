package mysql

import (
	"fmt"
	"testing"
)

var _tableName = "table_001"

func TestAll(t *testing.T) {
	t1 := &Table{}

	t1.setTable(_tableName, engineMyISAM, "", "")
	t1.AddPkBigInt("id", 20)
	t1.AddInt("int1", 11, 0, false)
	t1.AddInt("int2", 11, 0, false)
	t1.AddInt("int3", 11, 0, false)
	t1.AddBigInt("bigint2", 20, 0, false)
	t1.AddVarchar("vc1", 20, "")
	t1.AddVarchar("vc2", 200, "")
	t1.AddVarchar("vc3", 200, "")
	t1.AddText("t1")

	t1.AddKey("int1", 0)
	t1.AddUniqueKey("int3")

	t1.AddKey("vc1", 0)
	t1.AddKey("vc2", 10)
	t1.AddUniqueKey("vc3")

	fmt.Println(t1.SqlForCreate())
}
