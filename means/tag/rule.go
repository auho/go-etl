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
	Fixed() map[string]string
	FixedAlias() map[string]string
	FixedKeys() []string
	FixedKeysAlias() []string
	ItemsAlias() ([]map[string]string, error)
	ItemsForRegexp() ([]map[string]string, error)
}
