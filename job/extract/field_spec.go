package extract

type FieldSpec interface {
	Keys() []string
	DefaultValues() map[string]any
}
