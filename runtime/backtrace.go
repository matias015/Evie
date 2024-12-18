package runtime

import "fmt"

type Backtrace struct {
	Trace []string
}

func NewBacktrace() *Backtrace {
	return &Backtrace{Trace: []string{}}
}

func (b *Backtrace) AddTrace(mod string, line int) {
	b.Trace = append(b.Trace, mod+":"+fmt.Sprint(line))
}
