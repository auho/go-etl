package assistant

import (
	"github.com/auho/go-etl/v2/insight/assistant/tablestructure"
	"github.com/auho/go-simple-db/v2"
)

type Rowsor interface {
	GetDB() *go_simple_db.SimpleDB
	GetName() string
	GetIdName() string
	TableName() string
	CommandExec(*tablestructure.Command)
}

type Ruler interface {
	GetDB() *go_simple_db.SimpleDB
	GetName() string
	GetNameLength() int
	GetIdName() string
	GetLabels() map[string]int
	GetKeywordLength() int
	LabelsName() []string
	TableName() string
	KeywordName() string
	KeywordLenName() string
	KeywordNumName() string
	CommandExec(*tablestructure.Command)
}

type Dataor interface {
	GetDB() *go_simple_db.SimpleDB
	GetName() string
	GetIdName() string
	TableName() string
	CommandExec(*tablestructure.Command)
}
