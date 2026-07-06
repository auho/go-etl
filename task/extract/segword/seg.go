package segword

import (
	"strings"
	"sync"

	"github.com/yanyiwu/gojieba"
)

type Seg struct {
	userHmm bool
	jieBa   *gojieba.Jieba
	mu      sync.Mutex
}

func NewSeg() *Seg {
	return &Seg{
		userHmm: true,
		jieBa:   gojieba.NewJieba(),
	}
}

func (s *Seg) tag(contents []string) results {
	s.mu.Lock()
	defer s.mu.Unlock()

	var rets results
	for _, content := range contents {
		items := s.jieBa.Tag(content)
		if len(items) == 0 {
			continue
		}

		for _, item := range items {
			ret := strings.Split(item, "/")
			rets = append(rets, result{
				token: ret[0],
				flag:  ret[1],
			})

		}
	}

	return rets
}

func (s *Seg) Close() error {
	s.jieBa.Free()

	return nil
}
