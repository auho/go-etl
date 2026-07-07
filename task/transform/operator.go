package transform

type Operator interface {
	Title() string
	Fields() []string // source data 里的 key name
	Prepare() error
	Close() error
}

type SingleOperator interface {
	Operator
	Apply(map[string]any) (map[string]any, error)
}

type InsertOperator interface {
	Operator
	Keys() []string // 处理后的 key name
	// DefaultValues 返回默认值的独立副本。
	// 实现者必须克隆内部 map（如 maps.Clone），
	// 确保调用方修改返回值不影响内部状态。
	DefaultValues() map[string]any
	Apply(map[string]any) ([]map[string]any, error)
	State() []string
}

type UpdateOperator interface {
	SingleOperator
}

type TransferOperator interface {
	SingleOperator
}
