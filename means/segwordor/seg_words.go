package segwordor

import (
	"strings"

	"github.com/yanyiwu/gojieba"
)

type SegWords struct {
	userHmm bool
	jieBa   *gojieba.Jieba
}

func (sw *SegWords) prepare() {
	sw.userHmm = true
	sw.jieBa = gojieba.NewJieba()
}

// tag
// return [[ word, flag ]]
func (sw *SegWords) tag(contents []string) [][]string {
	results := make([][]string, 0)
	for _, content := range contents {
		items := sw.jieBa.Tag(content)
		if len(items) <= 0 {
			continue
		}

		for _, item := range items {
			results = append(results, strings.Split(item, "/"))
		}
	}

	if len(results) <= 0 {
		return nil
	}

	return results
}

func (sw *SegWords) Close() {
	sw.jieBa.Free()
}
