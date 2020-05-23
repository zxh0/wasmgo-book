package aot

import (
	"fmt"
	"strings"
)

type printer struct {
	sb *strings.Builder
}

func newPrinter() printer {
	return printer{sb: &strings.Builder{}}
}

func (p *printer) print(s string) {
	//fmt.Print(s)
	p.sb.WriteString(s)
}

func (p *printer) println(s string) {
	p.print(s)
	p.print("\n")
}

func (p *printer) printf(format string, a ...interface{}) {
	p.print(fmt.Sprintf(format, a...))
}

func (p *printer) printIf(cond bool, s1, s2 string) {
	if cond {
		p.print(s1)
	} else {
		p.print(s2)
	}
}

func (p *printer) String() string {
	return p.sb.String()
}
