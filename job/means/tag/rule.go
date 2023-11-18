package tag

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
