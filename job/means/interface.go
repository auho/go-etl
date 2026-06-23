package means

type Inserter interface {
	GetTitle() string
	GetKeys() []string
	DefaultValues() map[string]any
	Insert([]string) []map[string]any
	Prepare() error
	Close() error
}

type Updater interface {
	GetTitle() string
	Update([]string) map[string]any
	Prepare() error
	Close() error
}

type Rule interface {
	Name() string
	NameAlias() string
	TableName() string
	KeywordName() string
	KeywordNameAlias() string
	KeywordNumName() string
	KeywordNumNameAlias() string
	KeywordAmountName() string
	KeywordAmountNameAlias() string
	Labels() []string
	LabelsAlias() []string
	LabelNumName() string
	LabelNumNameAlias() string
	Tags() []string
	TagsAlias() []string
	ItemsAlias() ([]map[string]string, error)
	ItemsForRegexp() ([]map[string]string, error)
}
