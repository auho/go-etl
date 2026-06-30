package segword

import (
	"strings"

	"github.com/yanyiwu/gojieba"
)

type Seg struct {
	userHmm bool
	jieBa   *gojieba.Jieba
}

func NewSeg() *Seg {
	return &Seg{
		userHmm: true,
		jieBa:   gojieba.NewJieba(),
	}
}

func (s *Seg) tag(contents []string) results {
	var results results
	for _, content := range contents {
		items := s.jieBa.Tag(content)
		if len(items) <= 0 {
			continue
		}

		for _, item := range items {
			rets := strings.Split(item, "/")
			results = append(results, result{
				token: rets[0],
				flag:  rets[1],
			})

		}
	}

	return results
}

func (s *Seg) Close() error {
	s.jieBa.Free()

	return nil
}
