package means

type InsertMeans interface {
	GetTitle() string
	GetKeys() []string
	DefaultValues() map[string]any
	Insert([]string) []map[string]any
	Prepare() error
	Close() error
}

type UpdateMeans interface {
	GetTitle() string
	Update([]string) map[string]any
	Prepare() error
	Close() error
}

type Ruler interface {
	Name() string
	NameAlias() string
	TableName() string
	KeywordName() string
	KeywordNameAlias() string
	KeywordNumName() string
	KeywordNumNameAlias() string
	Labels() []string
	LabelsAlias() []string
	Fixed() map[string]any
	FixedAlias() map[string]any
	FixedKeys() []string
	FixedKeysAlias() []string
	ItemsAlias() ([]map[string]string, error)
	ItemsForRegexp() ([]map[string]string, error)
}
