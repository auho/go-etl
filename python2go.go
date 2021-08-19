package go_etl

import (
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
)

type Python2Go struct {
	f    *os.File
	name string
	path string
	c    string
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

	b, err := ioutil.ReadAll(p.f)
	if err != nil {
		panic(err)
	}

	p.c = string(b)

	return p
}
func (p *Python2Go) Conversion() {
	p.clean()
	p.conversionQuote()
	p.conversionComment()
	p.conversionVar()
	p.conversionSlice()
	p.conversionDict()

	p.save()
}

func (p *Python2Go) save() {
	err := ioutil.WriteFile(p.name, []byte(p.c), 0644)
	if err != nil {
		panic(err)
	}
}

func (p *Python2Go) clean() {
	re := regexp.MustCompile("(?m)^from .+")
	p.c = re.ReplaceAllString(p.c, "")
}

func (p *Python2Go) conversionQuote() {
	p.c = strings.ReplaceAll(p.c, "'", `"`)
}

func (p *Python2Go) conversionComment() {
	re := regexp.MustCompile(`"""([^"]+)"""`)
	p.c = re.ReplaceAllString(p.c, `/*$1*/`)
}

func (p *Python2Go) conversionVar() {
	re := regexp.MustCompile(`(?m)([^\s]+)\s=\s"([^"]+)"`)
	p.c = re.ReplaceAllString(p.c, `const $1 = "$2"`)
}

func (p *Python2Go) conversionSlice() {
	re := regexp.MustCompile(`([^\s]+)\s=\s\[([^\]]+)\]`)
	p.c = re.ReplaceAllString(p.c, `var $1 = []string{$2}`)
}

func (p *Python2Go) conversionDict() {
	re := regexp.MustCompile(`([^\s]+)\s=\s{([^}]+)}`)
	p.c = re.ReplaceAllString(p.c, `var $1 = map[string]string{$2}`)

}

func (p *Python2Go) content() {
	fmt.Println(p.c)
}
