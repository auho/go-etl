package transform

type Operator interface {
	Title() string
	GetFields() []string // source data 里的 key name
	Prepare() error
	Close() error
}

type SingleOperator interface {
	Operator
	Apply(map[string]any) (map[string]any, error)
}

type InsertOperator interface {
	Operator
	Keys() []string                // 处理后的 key name
	DefaultValues() map[string]any // 需要 implement clone important!
	Apply(map[string]any) ([]map[string]any, error)
	State() []string
}

type UpdateOperator interface {
	SingleOperator
}

type TransferOperator interface {
	SingleOperator
}
