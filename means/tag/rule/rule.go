package rule

var _ Ruler = (*DBRule)(nil)

type Ruler interface {
	TagsName() []string
	Items() []map[string]string
}
