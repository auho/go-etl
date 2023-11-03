package tag

type Ruler interface {
	Name() string
	TableName() string
	KeywordName() string
	KeywordNumName() string
	Labels() []string
	Fixed() map[string]string
	FixedKeys() []string
	Items() ([]map[string]string, error)
}
