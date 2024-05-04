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

func (s *Seg) tag(contents []string) Results {
	var results Results
	for _, content := range contents {
		items := s.jieBa.Tag(content)
		if len(items) <= 0 {
			continue
		}

		for _, item := range items {
			rets := strings.Split(item, "/")
			results = append(results, Result{
				Token: rets[0],
				Flag:  rets[1],
			})

		}
	}

	return results
}

func (s *Seg) Close() error {
	s.jieBa.Free()

	return nil
}
