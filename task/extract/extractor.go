package extract

// Extractor transforms content strings into structured rows.
// Keys and DefaultValues are prepared during Prepare so callers can
// retrieve them without re-computing.
type Extractor interface {
	Prepare() error
	Title() string
	Keys() []string
	// DefaultValues 返回默认值的独立副本。
	// 实现者必须克隆内部 map（如 maps.Clone），
	// 确保调用方修改返回值不影响内部状态。
	DefaultValues() map[string]any
	Extract(contents []string) Result
	Close() error
}
