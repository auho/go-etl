package go_etl

import (
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
)

type Python2Go struct {
	f    *os.File
	name string
	path string
}

func NewPython2Go(s string) *Python2Go {
	fs, err := os.Stat(s)
	if err != nil {
		panic(err)
	}

	f, err := os.OpenFile(s, os.O_RDONLY, 0755)
	if err != nil {
		panic(err)
	}

	p := &Python2Go{}
	p.name = fs.Name()
	p.path = f.Name()
	p.f = f

	return p
}

func (p *Python2Go) Conversion() {
	b, err := ioutil.ReadAll(p.f)
	if err != nil {
		panic(err)
	}

	s := string(b)

	re := regexp.MustCompile(`(?m)[^\s]+\s=\s'[^']+'`)
	items := re.FindAllString(s, -1)
	fmt.Println(items)

}
