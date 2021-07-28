package segworder

type SegWordsMeans struct {
	SegWords
}

func NewSegWordsMeans() *SegWordsMeans {
	sw := &SegWordsMeans{}
	sw.prepare()

	return sw
}

func (sw *SegWordsMeans) GetKeys() []string {
	return []string{"word", "flag"}
}

func (sw *SegWordsMeans) Insert(contents []string) [][]interface{} {
	results := sw.Tag(contents)
	if results == nil {
		return nil
	}

	items := make([][]interface{}, len(results), len(results))
	for k, result := range results {
		items[k] = []interface{}{result[0], result[1]}
	}

	return items
}
