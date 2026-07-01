package extract

type Rule interface {
	// Name returns the rule name.
	Name() string
	// NameAlias returns the alias for the rule name.
	NameAlias() string
	// TableName returns the table name associated with the rule.
	TableName() string
	// KeywordName returns the keyword field/column name.
	KeywordName() string
	// KeywordNameAlias returns the alias for the keyword field name.
	KeywordNameAlias() string
	// KeywordNumName returns the keyword number field/column name.
	KeywordNumName() string
	// KeywordNumNameAlias returns the alias for the keyword number field name.
	KeywordNumNameAlias() string
	// KeywordAmountName returns the keyword amount field/column name.
	KeywordAmountName() string
	// KeywordAmountNameAlias returns the alias for the keyword amount field name.
	KeywordAmountNameAlias() string
	// Labels returns the label names.
	Labels() []string
	// LabelsAlias returns the aliases for the label names.
	LabelsAlias() []string
	// LabelNumName returns the label number field/column name.
	LabelNumName() string
	// LabelNumNameAlias returns the alias for the label number field name.
	LabelNumNameAlias() string
	// Tags returns the tag names.
	Tags() []string
	// TagsAlias returns the aliases for the tag names.
	TagsAlias() []string
	// ItemsAlias returns the items resolved with their aliased column names.
	ItemsAlias() ([]map[string]string, error)
	// ItemsForRegexp returns the items formatted for regexp matching.
	ItemsForRegexp() ([]map[string]string, error)
}
